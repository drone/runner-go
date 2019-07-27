// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package clone

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEnvironDefault(t *testing.T) {
	c := Config{}
	a := Environ(c)
	b := map[string]string{
		"GIT_AUTHOR_NAME":     "drone",
		"GIT_AUTHOR_EMAIL":    "noreply@drone",
		"GIT_COMMITTER_NAME":  "drone",
		"GIT_COMMITTER_EMAIL": "noreply@drone",
		"GIT_TERMINAL_PROMPT": "0",
	}
	if diff := cmp.Diff(a, b); diff != "" {
		t.Fail()
		t.Log(diff)
	}
}

func TestEnviron(t *testing.T) {
	c := Config{
		User: User{
			Name:  "The Octocat",
			Email: "octocat@github.com",
		},
		Trace:      true,
		SkipVerify: true,
	}
	a := Environ(c)
	b := map[string]string{
		"GIT_AUTHOR_NAME":     "The Octocat",
		"GIT_AUTHOR_EMAIL":    "octocat@github.com",
		"GIT_COMMITTER_NAME":  "The Octocat",
		"GIT_COMMITTER_EMAIL": "octocat@github.com",
		"GIT_TERMINAL_PROMPT": "0",
		"GIT_TRACE":           "true",
		"GIT_SSL_NO_VERIFY":   "true",
	}
	if diff := cmp.Diff(a, b); diff != "" {
		t.Fail()
		t.Log(diff)
	}
}
