// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package clone

import "testing"

func TestFetchFlags(t *testing.T) {
	var args Args
	if got, want := fetchFlags(args), ""; got != want {
		t.Errorf("Want %q, got %q", want, got)
	}
	args.Tags = true
	if got, want := fetchFlags(args), "--tags"; got != want {
		t.Errorf("Want %q, got %q", want, got)
	}
	args.Tags = false
	args.Depth = 50
	if got, want := fetchFlags(args), "--depth=50"; got != want {
		t.Errorf("Want %q, got %q", want, got)
	}
}

func TestMergeFlags(t *testing.T) {
	var args Args
	if got, want := mergeFlags(args), ""; got != want {
		t.Errorf("Want %q, got %q", want, got)
	}
	args.NoFF = true
	if got, want := mergeFlags(args), "--no-ff"; got != want {
		t.Errorf("Want %q, got %q", want, got)
	}
}

func TestIsTag(t *testing.T) {
	tests := []struct {
		s string
		v bool
	}{
		{
			s: "refs/heads/master",
			v: false,
		},
		{
			s: "refs/pull/1/head",
			v: false,
		},
		{
			s: "refs/tags/v1.0.0",
			v: true,
		},
	}

	for _, test := range tests {
		if got, want := isTag(test.s), test.v; got != want {
			t.Errorf("Want tag %v for %s", want, test.s)
		}
	}
}

func TestIsPullRequst(t *testing.T) {
	tests := []struct {
		s string
		v bool
	}{
		{
			s: "refs/heads/master",
			v: false,
		},
		{
			s: "refs/pull/1/head",
			v: true,
		},
		{
			s: "refs/pull/2/merge",
			v: true,
		},
	}

	for _, test := range tests {
		if got, want := isPullRequest(test.s), test.v; got != want {
			t.Errorf("Want pull request %v for %s", want, test.s)
		}
	}
}
