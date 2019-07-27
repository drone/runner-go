// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package clone

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCommandsTag(t *testing.T) {
	args := Args{
		Depth:  50,
		Remote: "https://github.com/octocat/hello-world.git",
		Ref:    "refs/tags/v1.0.0",
	}
	got := Commands(args)
	want := []string{
		"git init",
		"git remote add origin https://github.com/octocat/hello-world.git",
		"git fetch --depth=50 origin +refs/tags/v1.0.0:",
		"git checkout -qf FETCH_HEAD",
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Fail()
		t.Log(diff)
	}
}

func TestCommandsBranch(t *testing.T) {
	args := Args{
		Branch: "develop",
		Commit: "3650a5d21bbf086fa8d2f16b0067ddeecfa604df",
		Depth:  50,
		NoFF:   true,
		Remote: "https://github.com/octocat/hello-world.git",
		Ref:    "refs/heads/develop",
		Tags:   true,
	}
	got := Commands(args)
	want := []string{
		"git init",
		"git remote add origin https://github.com/octocat/hello-world.git",
		"git fetch --depth=50 --tags origin +refs/heads/develop:",
		"git checkout 3650a5d21bbf086fa8d2f16b0067ddeecfa604df -b develop",
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Log(want)
		t.Fail()
		t.Log(diff)
	}
}

func TestCommandsPullRequest(t *testing.T) {
	args := Args{
		Branch: "master",
		Commit: "3650a5d21bbf086fa8d2f16b0067ddeecfa604df",
		Depth:  50,
		NoFF:   true,
		Remote: "https://github.com/octocat/hello-world.git",
		Ref:    "refs/pull/42/head",
		Tags:   true,
	}
	got := Commands(args)
	want := []string{
		"git init",
		"git remote add origin https://github.com/octocat/hello-world.git",
		"git fetch --depth=50 --tags origin +refs/heads/master:",
		"git checkout master",
		"git fetch origin refs/pull/42/head:",
		"git merge --no-ff 3650a5d21bbf086fa8d2f16b0067ddeecfa604df",
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Log(want)
		t.Fail()
		t.Log(diff)
	}
}
