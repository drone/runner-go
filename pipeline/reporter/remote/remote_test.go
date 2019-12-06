// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package remote

import (
	"context"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/client"
	"github.com/drone/runner-go/livelog"
	"github.com/drone/runner-go/pipeline"

	"github.com/google/go-cmp/cmp"
)

var nocontext = context.Background()

func TestReportStep(t *testing.T) {
	step := &drone.Step{Name: "clone"}
	state := &pipeline.State{
		Stage: &drone.Stage{
			Steps: []*drone.Step{step},
		},
	}

	c := new(mockClient)
	r := New(c)
	err := r.ReportStep(nocontext, state, step.Name)
	if err != nil {
		t.Error(err)
	}
	if state.Stage.Steps[0] != step {
		t.Errorf("Expect step updated, not replaced")
	}
	after := &drone.Step{
		Name:    "clone",
		ID:      1,
		StageID: 2,
		Started: 1561256080,
		Stopped: 1561256090,
		Version: 42,
	}
	if diff := cmp.Diff(after, step); diff != "" {
		t.Errorf("Expect response merged with step")
		t.Log(diff)
	}
}

func TestReportStage(t *testing.T) {
	stage := &drone.Stage{
		Created: 0,
		Updated: 0,
		Version: 0,
	}
	state := &pipeline.State{
		Stage: stage,
	}

	c := new(mockClient)
	r := New(c)
	err := r.ReportStage(nocontext, state)
	if err != nil {
		t.Error(err)
	}
	if state.Stage != stage {
		t.Errorf("Expect stage updated, not replaced")
	}
	after := &drone.Stage{
		Created: 1561256080,
		Updated: 1561256090,
		Version: 42,
	}
	if diff := cmp.Diff(after, state.Stage); diff != "" {
		t.Errorf("Expect response merged with stage")
		t.Log(diff)
	}
}

func TestStream(t *testing.T) {
	state := &pipeline.State{
		Stage: &drone.Stage{
			Steps: []*drone.Step{
				{
					ID:   1,
					Name: "clone",
				},
			},
		},
	}

	c := new(mockClient)
	r := New(c)
	w := r.Stream(nocontext, state, "clone")

	if _, ok := w.(*livelog.Writer); !ok {
		t.Errorf("Expect livelog writer")
	}
}

type mockClient struct {
	*client.HTTPClient
}

func (m *mockClient) Update(_ context.Context, stage *drone.Stage) error {
	stage.Version = 42
	stage.Created = 1561256080
	stage.Updated = 1561256090
	return nil
}

func (m *mockClient) UpdateStep(_ context.Context, step *drone.Step) error {
	step.ID = 1
	step.StageID = 2
	step.Started = 1561256080
	step.Stopped = 1561256090
	step.Version = 42
	return nil
}
