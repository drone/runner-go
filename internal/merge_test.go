// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package internal

import (
	"testing"

	"github.com/drone/drone-go/drone"

	"github.com/google/go-cmp/cmp"
)

func TestMergeStep(t *testing.T) {
	src := &drone.Step{
		ID:      1,
		StageID: 2,
		Started: 1561256095,
		Stopped: 1561256092,
		Version: 1,
	}
	dst := &drone.Step{
		ID: 1,
	}

	MergeStep(src, dst)
	if src == dst {
		t.Errorf("Except copy of step, got reference")
	}

	after := &drone.Step{
		ID:      1,
		StageID: 2,
		Started: 1561256095,
		Stopped: 1561256092,
		Version: 1,
	}
	if diff := cmp.Diff(after, src); diff != "" {
		t.Errorf("Expect src not modified")
		t.Log(diff)
	}
	if diff := cmp.Diff(after, dst); diff != "" {
		t.Errorf("Expect src values copied to dst")
		t.Log(diff)
	}
}

func TestMergeStage(t *testing.T) {
	dst := &drone.Stage{
		ID: 1,
		Steps: []*drone.Step{
			{
				ID: 1,
			},
		},
	}
	src := &drone.Stage{
		ID:      1,
		Created: 1561256095,
		Updated: 1561256092,
		Version: 1,
		Steps: []*drone.Step{
			{
				ID:      1,
				StageID: 2,
				Started: 1561256095,
				Stopped: 1561256092,
				Version: 1,
			},
		},
	}

	MergeStage(src, dst)
	if src == dst {
		t.Errorf("Except copy of stage, got reference")
	}

	after := &drone.Stage{
		ID:      1,
		Created: 1561256095,
		Updated: 1561256092,
		Version: 1,
		Steps: []*drone.Step{
			{
				ID:      1,
				StageID: 2,
				Started: 1561256095,
				Stopped: 1561256092,
				Version: 1,
			},
		},
	}
	if diff := cmp.Diff(after, dst); diff != "" {
		t.Errorf("Expect src values copied to dst")
		t.Log(diff)
	}
}
