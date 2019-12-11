// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package history implements a tracer that provides access to
// pipeline execution history.
package history

import (
	"context"
	"sync"
	"time"

	"github.com/drone/runner-go/internal"
	"github.com/drone/runner-go/pipeline"
)

var _ pipeline.Reporter = (*History)(nil)

// default number of historical entries.
const defaultLimit = 25

// History tracks pending, running and complete pipeline stages
// processed by the system.
type History struct {
	sync.Mutex
	base  pipeline.Reporter
	limit int
	items []*Entry
}

// New returns a new History recorder that wraps the base
// reporter.
func New(base pipeline.Reporter) *History {
	return &History{base: base}
}

// ReportStage adds or updates the pipeline history.
func (h *History) ReportStage(ctx context.Context, state *pipeline.State) error {
	h.Lock()
	h.update(state)
	h.prune()
	h.Unlock()
	return h.base.ReportStage(ctx, state)
}

// ReportStep adds or updates the pipeline history.
func (h *History) ReportStep(ctx context.Context, state *pipeline.State, name string) error {
	h.Lock()
	h.update(state)
	h.prune()
	h.Unlock()
	return h.base.ReportStep(ctx, state, name)
}

// Entries returns a list of entries.
func (h *History) Entries() []*Entry {
	h.Lock()
	var entries []*Entry
	for _, src := range h.items {
		dst := new(Entry)
		*dst = *src
		entries = append(entries, dst)
	}
	h.Unlock()
	return entries
}

// Entry returns the entry by id.
func (h *History) Entry(id int64) *Entry {
	h.Lock()
	defer h.Unlock()
	for _, src := range h.items {
		if src.Stage.ID == id {
			dst := new(Entry)
			*dst = *src
			return dst
		}
	}
	return nil
}

// Limit returns the history limit.
func (h *History) Limit() int {
	if h.limit == 0 {
		return defaultLimit
	}
	return h.limit
}

func (h *History) update(state *pipeline.State) {
	for _, v := range h.items {
		if v.Stage.ID == state.Stage.ID {
			v.Stage = internal.CloneStage(state.Stage)
			v.Build = internal.CloneBuild(state.Build)
			v.Repo = internal.CloneRepo(state.Repo)
			v.Updated = time.Now().UTC()
			return
		}
	}
	h.items = append(h.items, &Entry{
		Stage:   internal.CloneStage(state.Stage),
		Build:   internal.CloneBuild(state.Build),
		Repo:    internal.CloneRepo(state.Repo),
		Created: time.Now(),
		Updated: time.Now(),
	})
}

func (h *History) prune() {
	if len(h.items) > h.Limit() {
		h.items = h.items[:h.Limit()-1]
	}
}
