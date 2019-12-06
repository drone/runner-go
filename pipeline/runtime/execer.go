// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package runtime

import (
	"context"
	"sync"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/environ"
	"github.com/drone/runner-go/logger"
	"github.com/drone/runner-go/pipeline"

	"github.com/hashicorp/go-multierror"
	"github.com/natessilva/dag"
	"golang.org/x/sync/semaphore"
)

// Execer executes the pipeline.
type Execer struct {
	mu       sync.Mutex
	engine   Engine
	reporter pipeline.Reporter
	streamer pipeline.Streamer
	sem      *semaphore.Weighted
}

// NewExecer returns a new execer.
func NewExecer(
	reporter pipeline.Reporter,
	streamer pipeline.Streamer,
	engine Engine,
	threads int64,
) *Execer {
	exec := &Execer{
		reporter: reporter,
		streamer: streamer,
		engine:   engine,
	}
	if threads > 0 {
		// optional semaphore that limits the number of steps
		// that can execute concurrently.
		exec.sem = semaphore.NewWeighted(threads)
	}
	return exec
}

// Exec executes the intermediate representation of the pipeline
// and returns an error if execution fails.
func (e *Execer) Exec(ctx context.Context, spec Spec, state *pipeline.State) error {
	defer e.engine.Destroy(noContext, spec)

	if err := e.engine.Setup(noContext, spec); err != nil {
		state.FailAll(err)
		return e.reporter.ReportStage(noContext, state)
	}

	// create a directed graph, where each vertex in the graph
	// is a pipeline step.
	var d dag.Runner
	for i := 0; i < spec.StepLen(); i++ {
		step := spec.StepAt(i)
		d.AddVertex(step.GetName(), func() error {
			return e.exec(ctx, state, spec, step)
		})
	}

	// create the vertex edges from the values configured in the
	// depends_on attribute.
	for i := 0; i < spec.StepLen(); i++ {
		step := spec.StepAt(i)
		for _, dep := range step.GetDependencies() {
			d.AddEdge(dep, step.GetName())
		}
	}

	var result error
	if err := d.Run(); err != nil {
		multierror.Append(result, err)
	}

	// once pipeline execution completes, notify the state
	// manager that all steps are finished.
	state.FinishAll()
	if err := e.reporter.ReportStage(noContext, state); err != nil {
		multierror.Append(result, err)
	}
	return result
}

func (e *Execer) exec(ctx context.Context, state *pipeline.State, spec Spec, step Step) error {
	var result error

	select {
	case <-ctx.Done():
		state.Cancel()
		return nil
	default:
	}

	log := logger.FromContext(ctx)
	log = log.WithField("step.name", step.GetName())
	ctx = logger.WithContext(ctx, log)

	if e.sem != nil {
		// the semaphore limits the number of steps that can run
		// concurrently. acquire the semaphore and release when
		// the pipeline completes.
		if err := e.sem.Acquire(ctx, 1); err != nil {
			return nil
		}

		defer func() {
			// recover from a panic to ensure the semaphore is
			// released to prevent deadlock. we do not expect a
			// panic, however, we are being overly cautious.
			if r := recover(); r != nil {
				// TODO(bradrydzewski) log the panic.
			}
			// release the semaphore
			e.sem.Release(1)
		}()
	}

	switch {
	case state.Skipped():
		return nil
	case state.Cancelled():
		return nil
	case step.GetRunPolicy() == RunNever:
		return nil
	case step.GetRunPolicy() == RunAlways:
		break
	case step.GetRunPolicy() == RunOnFailure && state.Failed() == false:
		state.Skip(step.GetName())
		return e.reporter.ReportStep(noContext, state, step.GetName())
	case step.GetRunPolicy() == RunOnSuccess && state.Failed():
		state.Skip(step.GetName())
		return e.reporter.ReportStep(noContext, state, step.GetName())
	}

	state.Start(step.GetName())
	err := e.reporter.ReportStep(noContext, state, step.GetName())
	if err != nil {
		return err
	}

	copy := step.Clone()

	// the pipeline environment variables need to be updated to
	// reflect the current state of the build and stage.
	state.Lock()
	copy.SetEnviron(
		environ.Combine(
			copy.GetEnviron(),
			environ.Build(state.Build),
			environ.Stage(state.Stage),
			environ.Step(findStep(state, step.GetName())),
		),
	)
	state.Unlock()

	// writer used to stream build logs.
	wc := e.streamer.Stream(noContext, state, step.GetName())
	wc = newReplacer(wc, secretSlice(step))

	// if the step is configured as a daemon, it is detached
	// from the main process and executed separately.
	if step.IsDetached() {
		go func() {
			e.engine.Run(ctx, spec, copy, wc)
			wc.Close()
		}()
		return nil
	}

	exited, err := e.engine.Run(ctx, spec, copy, wc)

	// close the stream. If the session is a remote session, the
	// full log buffer is uploaded to the remote server.
	if err := wc.Close(); err != nil {
		multierror.Append(result, err)
	}

	if exited != nil {
		state.Finish(step.GetName(), exited.ExitCode)
		err := e.reporter.ReportStep(noContext, state, step.GetName())
		if err != nil {
			multierror.Append(result, err)
		}
		// if the exit code is 78 the system will skip all
		// subsequent pending steps in the pipeline.
		if exited.ExitCode == 78 {
			state.SkipAll()
		}
		return result
	}

	switch err {
	case context.Canceled, context.DeadlineExceeded:
		state.Cancel()
		return nil
	}

	// if the step failed with an internal error (as opposed to a
	// runtime error) the step is failed.
	state.Fail(step.GetName(), err)
	err = e.reporter.ReportStep(noContext, state, step.GetName())
	if err != nil {
		multierror.Append(result, err)
	}
	return result
}

// helper function returns the named step from the state.
func findStep(state *pipeline.State, name string) *drone.Step {
	for _, step := range state.Stage.Steps {
		if step.Name == name {
			return step
		}
	}
	panic("step not found: " + name)
}

// helper function returns an array of secrets from the
// pipeline step.
func secretSlice(step Step) []Secret {
	var secrets []Secret
	for i := 0; i < step.GetSecretLen(); i++ {
		secrets = append(secrets, step.GetSecretAt(i))
	}
	return secrets
}
