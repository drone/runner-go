// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Parity Public License
// that can be found in the LICENSE file.

// Package handler provides HTTP handlers that expose pipeline
// state and status.
package handler

import (
	"net/http"
	"sort"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/pipeline/history"
)

// HandleHealth returns a http.HandlerFunc that returns a 200
// if the service is healthly.
func HandleHealth(t *history.History) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO(bradrydzewski) iterate through the list of
		// pending or running stages and write an error message
		// if the timeout is exceeded.
		nocache(w)
		w.WriteHeader(200)
	}
}

// HandleIndex returns a http.HandlerFunc that displays a list
// of currently and previously executed builds.
func HandleIndex(t *history.History) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := t.Entries()

		s1 := history.ByTimestamp(d)
		s2 := history.ByStatus(d)
		sort.Sort(s1)
		sort.Sort(s2)

		if r.Header.Get("Accept") == "application/json" {
			nocache(w)
			renderJSON(w, d)
		} else {
			nocache(w)
			render(w, "index.tmpl", &data{Items: d})
		}
	}
}

// data is a template data structure that provides helper
// functions for calculating the system state.
type data struct {
	Items []*history.Entry
}

// helper function returns true if the history is empty.
func (d *data) Empty() bool {
	return len(d.Items) == 0
}

// helper function returns true if no running builds exists.
func (d *data) Idle() bool {
	for _, item := range d.Items {
		switch item.Stage.Status {
		case drone.StatusPending, drone.StatusRunning:
			return false
		}
	}
	return true
}
