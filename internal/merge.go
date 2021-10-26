// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package internal

import "github.com/drone/drone-go/drone"

// MergeStage merges the source stage with the destination.
func MergeStage(src, dst *drone.Stage) {
	dst.Version = src.Version
	dst.Created = src.Created
	dst.Updated = src.Updated
	for i, src := range src.Steps {
		dst := dst.Steps[i]
		MergeStep(src, dst)
	}
}

// MergeStep merges the source stage with the destination.
func MergeStep(src, dst *drone.Step) {
	dst.Version = src.Version
	dst.ID = src.ID
	dst.StageID = src.StageID
	dst.Started = src.Started
	dst.Stopped = src.Stopped
	dst.Version = src.Version
	dst.Schema = src.Schema
}
