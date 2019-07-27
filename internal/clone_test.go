// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package internal

import (
	"testing"

	"github.com/drone/drone-go/drone"

	"github.com/google/go-cmp/cmp"
)

func TestCloneRepo(t *testing.T) {
	src := &drone.Repo{
		ID:          1,
		UID:         "2",
		UserID:      3,
		Namespace:   "octocat",
		Name:        "hello-world",
		Slug:        "octocat/hello-world",
		SCM:         "git",
		HTTPURL:     "https://github.com/octocat/hello-world.git",
		SSHURL:      "git@github.com:octocat/hello-world.git",
		Link:        "https://github.com/octocat/hello-world",
		Branch:      "master",
		Private:     true,
		Visibility:  "public",
		Active:      true,
		Config:      ".drone.yml",
		Trusted:     true,
		Protected:   true,
		IgnoreForks: true,
		IgnorePulls: true,
		Timeout:     60,
		Counter:     50,
		Synced:      1561256365,
		Created:     1561256505,
		Updated:     1561256511,
		Version:     1,
	}
	dst := CloneRepo(src)
	if src == dst {
		t.Errorf("Except copy of repo, got reference")
	}
	if diff := cmp.Diff(src, dst); diff != "" {
		t.Errorf("Expect copy of values")
		t.Log(diff)
	}
}

func TestCloneBuild(t *testing.T) {
	src := &drone.Build{
		ID:           1,
		RepoID:       2,
		Number:       3,
		Parent:       4,
		Status:       drone.StatusFailing,
		Error:        "",
		Event:        drone.EventPush,
		Action:       "created",
		Link:         "https://github.com/octocat/Hello-World/commit/7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
		Timestamp:    1561256041,
		Title:        "",
		Message:      "updated README",
		Before:       "553c2077f0edc3d5dc5d17262f6aa498e69d6f8e",
		After:        "762941318ee16e59dabbacb1b4049eec22f0d303",
		Ref:          "refs/heads/master",
		Fork:         "spaceghost/hello-world",
		Source:       "develop",
		Target:       "master",
		Author:       "octocat",
		AuthorName:   "The Octocat",
		AuthorEmail:  "octocat@github.com",
		AuthorAvatar: "https://avatars2.githubusercontent.com/u/251370",
		Sender:       "spaceghost",
		Params:       map[string]string{"memory": "high"},
		Cron:         "nightly",
		Deploy:       "production",
		Started:      1561256065,
		Finished:     1561256082,
		Created:      1561256050,
		Updated:      1561256052,
		Version:      1,
		Stages: []*drone.Stage{
			{
				ID:        1,
				BuildID:   2,
				Number:    3,
				Name:      "build",
				Kind:      "pipeline",
				Type:      "docker",
				Status:    drone.StatusPassing,
				Error:     "",
				ErrIgnore: true,
				ExitCode:  0,
				Machine:   "server1",
				OS:        "linux",
				Arch:      "amd64",
				Variant:   "",
				Kernel:    "",
				Limit:     0,
				Started:   1561256065,
				Stopped:   1561256505,
				Created:   1561256356,
				Updated:   1561256082,
				Version:   1,
				OnSuccess: true,
				OnFailure: true,
				DependsOn: []string{"clone"},
				Labels:    map[string]string{"foo": "bar"},
			},
		},
	}
	dst := CloneBuild(src)
	if src == dst {
		t.Errorf("Except copy of build, got reference")
	}
	if diff := cmp.Diff(src, dst); diff != "" {
		t.Errorf("Expect copy of values")
		t.Log(diff)
	}

	if src.Stages[0] == dst.Stages[0] {
		t.Errorf("Except copy of stages, got reference")
	}
	if diff := cmp.Diff(src.Stages[0], dst.Stages[0]); diff != "" {
		t.Errorf("Expect copy of stage values")
		t.Log(diff)
	}
}

func TestCloneStage(t *testing.T) {
	src := &drone.Stage{
		ID:        1,
		BuildID:   2,
		Number:    3,
		Name:      "build",
		Kind:      "pipeline",
		Type:      "docker",
		Status:    drone.StatusPassing,
		Error:     "",
		ErrIgnore: true,
		ExitCode:  0,
		Machine:   "server1",
		OS:        "linux",
		Arch:      "amd64",
		Variant:   "",
		Kernel:    "",
		Limit:     0,
		Started:   1561256065,
		Stopped:   1561256505,
		Created:   1561256356,
		Updated:   1561256082,
		Version:   1,
		OnSuccess: true,
		OnFailure: true,
		DependsOn: []string{"clone"},
		Labels:    map[string]string{"foo": "bar"},
		Steps: []*drone.Step{
			{
				ID:        1,
				StageID:   2,
				Number:    3,
				Name:      "foo",
				Status:    drone.StatusFailing,
				Error:     "",
				ErrIgnore: false,
				ExitCode:  255,
				Started:   1561256065,
				Stopped:   1561256082,
				Version:   1,
			},
		},
	}
	dst := CloneStage(src)
	if src == dst {
		t.Errorf("Except copy of stage, got reference")
	}
	if src.Steps[0] == dst.Steps[0] {
		t.Errorf("Except copy of step, got reference")
	}
	if diff := cmp.Diff(src, dst); diff != "" {
		t.Errorf("Expect copy of step values")
		t.Log(diff)
	}
	if diff := cmp.Diff(src.Steps[0], dst.Steps[0]); diff != "" {
		t.Errorf("Expect copy of stage values")
		t.Log(diff)
	}
}

func TestCloneStep(t *testing.T) {
	src := &drone.Step{
		ID:        1,
		StageID:   2,
		Number:    3,
		Name:      "foo",
		Status:    drone.StatusFailing,
		Error:     "",
		ErrIgnore: false,
		ExitCode:  255,
		Started:   1561256065,
		Stopped:   1561256082,
		Version:   1,
	}
	dst := CloneStep(src)
	if src == dst {
		t.Errorf("Except copy of step, got reference")
	}
	if diff := cmp.Diff(src, dst); diff != "" {
		t.Errorf("Expect copy of values")
		t.Log(diff)
	}
	dst.ID = 101
	dst.StageID = 102
	dst.Number = 103
	dst.Name = "bar"
	dst.ErrIgnore = true
	dst.ExitCode = 0
	dst.Status = drone.StatusPassing
	dst.Started = 1561256356
	dst.Stopped = 1561256365
	dst.Version = 2
	if diff := cmp.Diff(src, dst); diff == "" {
		t.Errorf("Expect copy of values, got reference")
	}
}
