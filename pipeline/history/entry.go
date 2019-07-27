// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package history

import (
	"time"

	"github.com/drone/drone-go/drone"
)

// Entry represents a history entry.
type Entry struct {
	Stage   *drone.Stage `json:"stage"`
	Build   *drone.Build `json:"build"`
	Repo    *drone.Repo  `json:"repo"`
	Created time.Time    `json:"created"`
	Updated time.Time    `json:"updated"`
}

// ByTimestamp sorts a list of entries by timestamp
type ByTimestamp []*Entry

func (a ByTimestamp) Len() int      { return len(a) }
func (a ByTimestamp) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByTimestamp) Less(i, j int) bool {
	return a[i].Stage.ID > a[j].Stage.ID
}

// ByStatus sorts a list of entries by status
type ByStatus []*Entry

func (a ByStatus) Len() int      { return len(a) }
func (a ByStatus) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func (a ByStatus) Less(i, j int) bool {
	return order(a[i].Stage) < order(a[j].Stage)
}

func order(stage *drone.Stage) int64 {
	switch stage.Status {
	case drone.StatusPending:
		return 0
	case drone.StatusRunning:
		return 1
	default:
		return 2
	}
}
