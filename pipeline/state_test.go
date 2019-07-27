// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package pipeline

import (
	"errors"
	"testing"

	"github.com/drone/drone-go/drone"
)

func TestStateKill(t *testing.T) {
	step := &drone.Step{Name: "clone", Status: drone.StatusPending}
	state := &State{
		Stage: &drone.Stage{
			Steps: []*drone.Step{step},
		},
	}

	state.kill(step)
	if got, want := step.Status, drone.StatusPending; got != want {
		t.Errorf("Expect a non-running step cannot be killed")
	}

	step.Status = drone.StatusRunning
	state.kill(step)
	if got, want := step.Status, drone.StatusKilled; got != want {
		t.Errorf("Want status %s, got %s", want, got)
	}
	if got, want := step.Error, ""; got != want {
		t.Errorf("Want error %q, got %q", want, got)
	}
	if got, want := step.ExitCode, 137; got != want {
		t.Errorf("Want exit code %d, got %d", want, got)
	}
	if got, want := step.Stopped, step.Started; got != want {
		t.Errorf("Want stopped %d, got %d", want, got)
	}
	if step.Started == 0 {
		t.Errorf("Expect step started is non-zero value")
	}
}

func TestStateKilled(t *testing.T) {
	state := &State{}
	state.Stage = &drone.Stage{Status: drone.StatusError}
	if state.killed() == true {
		t.Errorf("Expect killed false, got true")
	}
	state.Stage.Status = drone.StatusKilled
	if state.killed() == false {
		t.Errorf("Expect killed true, got false")
	}
}

func TestStateFinished(t *testing.T) {
	step := &drone.Step{}
	state := &State{
		Stage: &drone.Stage{
			Steps: []*drone.Step{step},
		},
	}
	step.Status = drone.StatusRunning
	if state.finished() == true {
		t.Errorf("Expect finished false")
	}
	step.Status = drone.StatusPending
	if state.finished() == true {
		t.Errorf("Expect finished false")
	}
	step.Status = drone.StatusKilled
	if state.finished() == false {
		t.Errorf("Expect finished true")
	}
}

func TestStateSkipped(t *testing.T) {
	state := &State{
		Stage: &drone.Stage{
			Steps: []*drone.Step{
				{Status: drone.StatusPassing},
				{Status: drone.StatusRunning},
			},
		},
	}
	if state.skipped() == true {
		t.Errorf("Expect skipped false")
	}

	state.Stage.Steps[1].Status = drone.StatusPassing
	if state.skipped() == true {
		t.Errorf("Expect skipped false")
	}

	state.Stage.Steps[1].Status = drone.StatusSkipped
	if state.skipped() == false {
		t.Errorf("Expect skipped true")
	}
}

func TestStateStarted(t *testing.T) {
	step := &drone.Step{Name: "clone", Status: drone.StatusPending}
	state := &State{
		Stage: &drone.Stage{
			Steps: []*drone.Step{step},
		},
	}

	state.start(step)
	if got, want := step.Status, drone.StatusRunning; got != want {
		t.Errorf("Want status %s, got %s", want, got)
	}
	if got, want := step.Error, ""; got != want {
		t.Errorf("Want error %q, got %q", want, got)
	}
	if got, want := step.ExitCode, 0; got != want {
		t.Errorf("Want exit code %d, got %d", want, got)
	}
	if got, want := step.Stopped, int64(0); got != want {
		t.Errorf("Want stopped %d, got %d", want, got)
	}
	if step.Started == 0 {
		t.Errorf("Expect step started is non-zero value")
	}
}

func TestStateFinish(t *testing.T) {
	t.Skip()
}

func TestStateFail(t *testing.T) {
	step := &drone.Step{Name: "clone"}
	state := &State{
		Stage: &drone.Stage{
			Steps: []*drone.Step{step},
		},
	}
	state.fail(step, errors.New("this is an error"))
	if got, want := step.Status, drone.StatusError; got != want {
		t.Errorf("Want status %s, got %s", want, got)
	}
	if got, want := step.Error, "this is an error"; got != want {
		t.Errorf("Want error %q, got %q", want, got)
	}
	if got, want := step.ExitCode, 255; got != want {
		t.Errorf("Want exit code %d, got %d", want, got)
	}
	if got, want := step.Stopped, step.Started; got != want {
		t.Errorf("Want started %d, got %d", want, got)
	}
	if step.Stopped == 0 {
		t.Errorf("Expect step stopped is non-zero value")
	}
}

func TestStateFind(t *testing.T) {
	step := &drone.Step{Name: "clone"}
	state := &State{
		Stage: &drone.Stage{
			Steps: []*drone.Step{step},
		},
	}
	if got, want := state.find("clone"), step; got != want {
		t.Errorf("Expect find returns the named step")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expect recover from panic")
		}
	}()

	state.find("test")
}
