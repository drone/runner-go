// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package history

import (
	"sort"
	"testing"

	"github.com/drone/drone-go/drone"

	"github.com/google/go-cmp/cmp"
)

func TestSort(t *testing.T) {
	before := []*Entry{
		{Stage: &drone.Stage{ID: 1, Status: drone.StatusPassing}},
		{Stage: &drone.Stage{ID: 2, Status: drone.StatusPassing}},
		{Stage: &drone.Stage{ID: 3, Status: drone.StatusPending}},
		{Stage: &drone.Stage{ID: 4, Status: drone.StatusRunning}},
		{Stage: &drone.Stage{ID: 5, Status: drone.StatusPassing}},
	}

	after := []*Entry{
		{Stage: &drone.Stage{ID: 3, Status: drone.StatusPending}},
		{Stage: &drone.Stage{ID: 4, Status: drone.StatusRunning}},
		{Stage: &drone.Stage{ID: 5, Status: drone.StatusPassing}},
		{Stage: &drone.Stage{ID: 2, Status: drone.StatusPassing}},
		{Stage: &drone.Stage{ID: 1, Status: drone.StatusPassing}},
	}

	s1 := ByTimestamp(before)
	s2 := ByStatus(before)
	sort.Sort(s1)
	sort.Sort(s2)

	if diff := cmp.Diff(before, after); diff != "" {
		t.Errorf("Expect entries sorted by status")
		t.Log(diff)
	}
}
