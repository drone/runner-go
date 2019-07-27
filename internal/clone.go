// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package internal

import "github.com/drone/drone-go/drone"

// CloneRepo returns a copy of the Repository.
func CloneRepo(src *drone.Repo) *drone.Repo {
	dst := new(drone.Repo)
	*dst = *src
	return dst
}

// CloneBuild returns a copy of the Build.
func CloneBuild(src *drone.Build) *drone.Build {
	dst := new(drone.Build)
	*dst = *src
	dst.Stages = append(src.Stages[:0:0], src.Stages...)
	dst.Params = map[string]string{}
	for k, v := range src.Params {
		dst.Params[k] = v
	}
	for i, v := range src.Stages {
		dst.Stages[i] = CloneStage(v)
	}
	return dst
}

// CloneStage returns a copy of the Stage.
func CloneStage(src *drone.Stage) *drone.Stage {
	dst := new(drone.Stage)
	*dst = *src
	dst.DependsOn = append(src.DependsOn[:0:0], src.DependsOn...)
	dst.Steps = append(src.Steps[:0:0], src.Steps...)
	dst.Labels = map[string]string{}
	for k, v := range src.Labels {
		dst.Labels[k] = v
	}
	for i, v := range src.Steps {
		dst.Steps[i] = CloneStep(v)
	}
	return dst
}

// CloneStep returns a copy of the Step.
func CloneStep(src *drone.Step) *drone.Step {
	dst := new(drone.Step)
	*dst = *src
	return dst
}
