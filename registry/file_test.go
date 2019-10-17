// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package registry

import (
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/google/go-cmp/cmp"
)

func TestFile(t *testing.T) {
	p := File("auths/testdata/config.json")
	got, err := p.List(noContext, nil)
	if err != nil {
		t.Error(err)
		return
	}
	want := []*drone.Registry{
		{
			Address:  "index.docker.io",
			Username: "octocat",
			Password: "correct-horse-battery-staple",
		},
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf(diff)
	}
}

func TestFileEmptyPath(t *testing.T) {
	p := File("")
	out, err := p.List(noContext, nil)
	if err != nil {
		t.Error(err)
	}
	if len(out) != 0 {
		t.Errorf("Expect empty registry credentials")
	}
}
