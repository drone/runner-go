// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package runtime

import (
	"context"
	"io"
)

type (
	// Engine is the interface that must be implemented by a
	// pipeline execution engine.
	Engine interface {
		// Setup the pipeline environment.
		Setup(context.Context, Spec) error

		// Destroy the pipeline environment.
		Destroy(context.Context, Spec) error

		// Run runs the pipeline step.
		Run(context.Context, Spec, Step, io.Writer) (*State, error)
	}

	// Spec is an interface that must be implemented by all
	// pipeline specifications.
	Spec interface {
		// StepAt returns the step at the specified index.
		StepAt(int) Step

		// StepLen returns the number of steps.
		StepLen() int
	}

	// Step is an interface that must be implemented by all
	// pipeline steps.
	Step interface {
		// GetName returns the step name.
		GetName() string

		// GetDependencies returns the step dependencies
		// used to define an execution graph.
		GetDependencies() []string

		// GetEnviron returns the step environment variables.
		GetEnviron() map[string]string

		// SetEnviron updates the step environment variables.
		SetEnviron(map[string]string)

		// GetErrPolicy returns the step error policy.
		GetErrPolicy() ErrPolicy

		// GetRunPolicy returns the step run policy.
		GetRunPolicy() RunPolicy

		// GetSecretAt returns the secret at the specified
		// index.
		GetSecretAt(int) Secret

		// GetSecretLen returns the number of secrets.
		GetSecretLen() int

		// IsDetached returns true if the step is detached
		// and executed in the background.
		IsDetached() bool

		// Clone returns a copy of the Step.
		Clone() Step
	}

	// State reports the step state.
	State struct {
		// ExitCode returns the exit code of the exited step.
		ExitCode int

		// GetExited reports whether the step has exited.
		Exited bool

		// OOMKilled reports whether the step has been
		// killed by the process manager.
		OOMKilled bool
	}

	// Secret is an interface that must be implemented
	// by all pipeline secrets.
	Secret interface {
		// GetName returns the secret name.
		GetName() string

		// GetValue returns the secret value.
		GetValue() string

		// IsMasked returns true if the secret value should
		// be masked. If true the secret value is masked in
		// the logs.
		IsMasked() bool
	}
)
