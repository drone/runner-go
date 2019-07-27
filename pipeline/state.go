// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package pipeline

import (
	"sync"
	"time"

	"github.com/drone/drone-go/drone"
)

// State stores the pipeline state.
type State struct {
	sync.Mutex

	Build  *drone.Build
	Repo   *drone.Repo
	Stage  *drone.Stage
	System *drone.System
}

// Cancel cancels the pipeline.
func (s *State) Cancel() {
	s.Lock()
	s.skipall()
	s.killall()
	s.update()
	s.Unlock()
}

// Cancelled returns true if the pipeline is cancelled.
func (s *State) Cancelled() bool {
	s.Lock()
	v := s.killed()
	s.Unlock()
	return v
}

// Fail fails the named pipeline step with error.
func (s *State) Fail(name string, err error) {
	s.Lock()
	v := s.find(name)
	s.fail(v, err)
	s.update()
	s.Unlock()
}

// FailAll fails the entire pipeline.
func (s *State) FailAll(err error) {
	s.Lock()
	s.failAll(err)
	s.skipall()
	s.update()
	s.Unlock()
}

// Failed returns true if the pipeline failed.
func (s *State) Failed() bool {
	s.Lock()
	v := s.failed()
	s.Unlock()
	return v
}

// Skip skips the named pipeline step.
func (s *State) Skip(name string) {
	s.Lock()
	v := s.find(name)
	s.skip(v)
	s.update()
	s.Unlock()
}

// SkipAll skips all pipeilne steps.
func (s *State) SkipAll() {
	s.Lock()
	s.skipall()
	s.update()
	s.Unlock()
}

// Skipped returns true if all pipeline steps are skipped.
func (s *State) Skipped() bool {
	s.Lock()
	v := s.skipped()
	s.Unlock()
	return v
}

// Start sets the named pipeline step to started.
func (s *State) Start(name string) {
	s.Lock()
	v := s.find(name)
	s.start(v)
	s.Unlock()
}

// Finish sets the pipeline step to finished.
func (s *State) Finish(name string, code int) {
	s.Lock()
	v := s.find(name)
	s.finish(v, code)
	s.update()
	s.Unlock()
}

// FinishAll finishes all pipeline steps.
func (s *State) FinishAll() {
	s.Lock()
	s.finishAll()
	s.update()
	s.Unlock()
}

// Find returns the named pipeline step.
func (s *State) Find(name string) *drone.Step {
	s.Lock()
	v := s.find(name)
	s.Unlock()
	return v
}

//
// Helper functions. INTERNAL USE ONLY
//

// helper function skips all pipeline steps.
func (s *State) skipall() {
	for _, v := range s.Stage.Steps {
		s.skip(v)
	}
}

// helper function that updates the state of an individual step
// to indicate the step to skipped.
func (s *State) skip(v *drone.Step) {
	if v.Status == drone.StatusPending {
		v.Started = time.Now().Unix()
		v.Stopped = time.Now().Unix()
		v.Status = drone.StatusSkipped
		v.ExitCode = 0
		v.Error = ""
	}
}

// helper function returns true if the overall pipeline is
// finished and remaining steps skipped.
func (s *State) skipped() bool {
	if s.finished() == false {
		return false
	}
	for _, v := range s.Stage.Steps {
		if v.Status == drone.StatusSkipped {
			return true
		}
	}
	return false
}

// helper function kills all pipeline steps.
func (s *State) killall() {
	s.Stage.Error = ""
	s.Stage.ExitCode = 0
	s.Stage.Status = drone.StatusKilled
	s.Stage.Stopped = time.Now().Unix()
	if s.Stage.Started == 0 {
		s.Stage.Started = s.Stage.Stopped
	}
	for _, v := range s.Stage.Steps {
		s.kill(v)
	}
}

// helper function that updates the state of an individual step
// to indicate the step to killed.
func (s *State) kill(v *drone.Step) {
	if v.Status == drone.StatusRunning {
		v.Status = drone.StatusKilled
		v.Stopped = time.Now().Unix()
		v.ExitCode = 137
		v.Error = ""
		if v.Started == 0 {
			v.Started = v.Stopped
		}
	}
}

// helper function returns true if the overall pipeline status
// is killed.
func (s *State) killed() bool {
	return s.Stage.Status == drone.StatusKilled
}

// helper function that updates the state of an individual step
// to indicate the step is started.
func (s *State) start(v *drone.Step) {
	if v.Status == drone.StatusPending {
		v.Status = drone.StatusRunning
		v.Started = time.Now().Unix()
		v.Stopped = 0
		v.ExitCode = 0
		v.Error = ""
	}
}

// helper function updates the state of an individual step
// based on the exit code.
func (s *State) finish(v *drone.Step, code int) {
	switch v.Status {
	case drone.StatusRunning, drone.StatusPending:
	default:
		return
	}
	v.ExitCode = code
	v.Stopped = time.Now().Unix()
	if v.Started == 0 {
		v.Started = v.Stopped
	}
	switch code {
	case 0, 78:
		v.Status = drone.StatusPassing
	default:
		v.Status = drone.StatusFailing
	}
}

// helper function returns true if the overall pipeline status
// is failing.
func (s *State) finished() bool {
	for _, v := range s.Stage.Steps {
		switch v.Status {
		case drone.StatusRunning, drone.StatusPending:
			return false
		}
	}
	return true
}

// helper function finishes all pipeline steps.
func (s *State) finishAll() {
	for _, v := range s.Stage.Steps {
		switch v.Status {
		case drone.StatusPending:
			s.skip(v)
		case drone.StatusRunning:
			s.finish(v, 0)
		}
	}
	switch s.Stage.Status {
	case drone.StatusRunning, drone.StatusPending:
		s.Stage.Stopped = time.Now().Unix()
		s.Stage.Status = drone.StatusPassing
		if s.Stage.Started == 0 {
			s.Stage.Started = s.Stage.Stopped
		}
		if s.failed() {
			s.Stage.Status = drone.StatusFailing
		}
	default:
		if s.Stage.Started == 0 {
			s.Stage.Started = time.Now().Unix()
		}
		if s.Stage.Stopped == 0 {
			s.Stage.Stopped = time.Now().Unix()
		}
	}
}

// helper function fails an individual step.
func (s *State) fail(v *drone.Step, err error) {
	v.Status = drone.StatusError
	v.Error = err.Error()
	v.ExitCode = 255
	v.Stopped = time.Now().Unix()
	if v.Started == 0 {
		v.Started = v.Stopped
	}
}

// helper function fails the overall pipeline.
func (s *State) failAll(err error) {
	switch s.Stage.Status {
	case drone.StatusPending,
		drone.StatusRunning:
		s.Stage.Status = drone.StatusError
		s.Stage.Error = err.Error()
		s.Stage.Stopped = time.Now().Unix()
		if s.Stage.Started == 0 {
			s.Stage.Started = s.Stage.Stopped
		}
	}
}

// helper function returns true if the overall pipeline status
// is failing.
func (s *State) failed() bool {
	switch s.Stage.Status {
	case drone.StatusFailing,
		drone.StatusError,
		drone.StatusKilled:
		return true
	}
	for _, v := range s.Stage.Steps {
		if v.ErrIgnore {
			continue
		}
		switch v.Status {
		case drone.StatusFailing,
			drone.StatusError,
			drone.StatusKilled:
			return true
		}
	}
	return false
}

// helper function updates the build and stage based on the
// aggregate
func (s *State) update() {
	for _, v := range s.Stage.Steps {
		switch v.Status {
		case drone.StatusKilled:
			s.Stage.ExitCode = 137
			s.Stage.Status = drone.StatusKilled
			s.Build.Status = drone.StatusKilled
			return
		case drone.StatusError:
			s.Stage.Status = drone.StatusError
			s.Build.Status = drone.StatusError
			return
		case drone.StatusFailing:
			if v.ErrIgnore == false {
				s.Stage.Status = drone.StatusFailing
				s.Build.Status = drone.StatusFailing
				return
			}
		}
	}
}

// helper function finds the step by name.
func (s *State) find(name string) *drone.Step {
	for _, step := range s.Stage.Steps {
		if step.Name == name {
			return step
		}
	}
	panic("step not found: " + name)
}
