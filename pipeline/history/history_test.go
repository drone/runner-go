// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package history

import (
	"context"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/pipeline"

	"github.com/google/go-cmp/cmp"
)

func TestReportStage(t *testing.T) {
	r := &drone.Repo{ID: 1}
	b := &drone.Build{ID: 2, Params: map[string]string{}}
	s := &drone.Stage{ID: 3, Labels: map[string]string{}}

	v := New(&nopReporter{})
	v.ReportStage(nil, &pipeline.State{
		Repo:  r,
		Build: b,
		Stage: s,
	})

	if v.items[0].Repo == r {
		t.Errorf("Expect copy of repository")
	}
	if v.items[0].Build == b {
		t.Errorf("Expect copy of build")
	}
	if v.items[0].Stage == s {
		t.Errorf("Expect copy of stage")
	}

	if diff := cmp.Diff(v.items[0].Repo, r); diff != "" {
		t.Errorf("Expect repository data copied")
		t.Log(diff)
	}
	if diff := cmp.Diff(v.items[0].Build, b); diff != "" {
		t.Errorf("Expect build data copied")
		t.Log(diff)
	}
	if diff := cmp.Diff(v.items[0].Stage, s); diff != "" {
		t.Errorf("Expect stage data copied")
		t.Log(diff)
	}
	if v.items[0].Updated.IsZero() {
		t.Errorf("Expect created timestamp non-zero")
	}
	if v.items[0].Created.IsZero() {
		t.Errorf("Expect updated timestamp non-zero")
	}
}

func TestReportStep(t *testing.T) {
	r := &drone.Repo{ID: 1}
	b := &drone.Build{ID: 2, Params: map[string]string{}}
	s := &drone.Stage{ID: 3, Labels: map[string]string{}}

	v := New(&nopReporter{})
	v.ReportStep(nil, &pipeline.State{
		Repo:  r,
		Build: b,
		Stage: s,
	}, "foo")

	if v.items[0].Repo == r {
		t.Errorf("Expect copy of repository")
	}
	if v.items[0].Build == b {
		t.Errorf("Expect copy of build")
	}
	if v.items[0].Stage == s {
		t.Errorf("Expect copy of stage")
	}

	if diff := cmp.Diff(v.items[0].Repo, r); diff != "" {
		t.Errorf("Expect repository data copied")
		t.Log(diff)
	}
	if diff := cmp.Diff(v.items[0].Build, b); diff != "" {
		t.Errorf("Expect build data copied")
		t.Log(diff)
	}
	if diff := cmp.Diff(v.items[0].Stage, s); diff != "" {
		t.Errorf("Expect stage data copied")
		t.Log(diff)
	}
	if v.items[0].Updated.IsZero() {
		t.Errorf("Expect created timestamp non-zero")
	}
	if v.items[0].Created.IsZero() {
		t.Errorf("Expect updated timestamp non-zero")
	}
}

func TestEntries(t *testing.T) {
	s := new(drone.Stage)
	v := History{}
	v.items = append(v.items, &Entry{Stage: s})

	list := v.Entries()
	if got, want := len(list), len(v.items); got != want {
		t.Errorf("Want %d entries, got %d", want, got)
	}

	if v.items[0] == list[0] {
		t.Errorf("Expect copy of Entry, got reference")
	}
	if v.items[0].Stage != list[0].Stage {
		t.Errorf("Expect reference to Stage, got copy")
	}
}

func TestEntry(t *testing.T) {
	s1 := &drone.Stage{ID: 1}
	s2 := &drone.Stage{ID: 2}
	v := History{}
	v.items = append(v.items, &Entry{Stage: s1}, &Entry{Stage: s2})

	if got := v.Entry(99); got != nil {
		t.Errorf("Want nil when stage not found")
	}
	if got := v.Entry(s1.ID); got == nil {
		t.Errorf("Want entry by stage ID, got nil")
		return
	}
}

func TestLimit(t *testing.T) {
	v := History{}
	if got, want := v.Limit(), defaultLimit; got != want {
		t.Errorf("Want default limit %d, got %d", want, got)
	}
	v.limit = 5
	if got, want := v.Limit(), 5; got != want {
		t.Errorf("Want custom limit %d, got %d", want, got)
	}
}

func TestInsert(t *testing.T) {
	stage := &drone.Stage{ID: 42, Labels: map[string]string{}}
	state := &pipeline.State{
		Stage: stage,
		Repo:  &drone.Repo{},
		Build: &drone.Build{},
	}

	v := History{}
	v.update(state)

	if v.items[0].Stage == stage {
		t.Errorf("Expect stage replaced")
	}
	if v.items[0].Updated.IsZero() {
		t.Errorf("Expect entry timestamp updated")
	}
	if diff := cmp.Diff(v.items[0].Stage, stage); diff != "" {
		t.Errorf("Expect stage data copied")
		t.Log(diff)
	}
}

func TestUpdate(t *testing.T) {
	stage1 := &drone.Stage{ID: 42, Labels: map[string]string{}}
	stage2 := &drone.Stage{ID: 42, Labels: map[string]string{}}
	state := &pipeline.State{
		Stage: stage1,
		Repo:  &drone.Repo{},
		Build: &drone.Build{},
	}

	v := History{}
	v.items = append(v.items, &Entry{Stage: stage1})
	v.update(state)

	if v.items[0].Stage == stage1 {
		t.Errorf("Expect stage replaced")
	}
	if v.items[0].Stage == stage2 {
		t.Errorf("Expect stage replaced")
	}
	if v.items[0].Updated.IsZero() {
		t.Errorf("Expect entry timestamp updated")
	}
	if diff := cmp.Diff(v.items[0].Stage, stage1); diff != "" {
		t.Errorf("Expect stage data copied")
		t.Log(diff)
	}
}

func TestPrune(t *testing.T) {
	v := History{}
	v.limit = 3
	v.items = append(v.items, nil)
	v.items = append(v.items, nil)
	v.items = append(v.items, nil)
	v.items = append(v.items, nil)
	v.items = append(v.items, nil)
	v.prune()
	if got, want := len(v.items), 2; got != want {
		t.Errorf("Want pruned entry len %d, got %d", want, got)
	}
}

type nopReporter struct{}

func (*nopReporter) ReportStage(context.Context, *pipeline.State) error        { return nil }
func (*nopReporter) ReportStep(context.Context, *pipeline.State, string) error { return nil }
