// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package registry

import (
	"testing"

	"github.com/drone/drone-go/drone"
)

func TestStatic(t *testing.T) {
	a := &drone.Registry{}
	b := &drone.Registry{}
	p := Static([]*drone.Registry{a, b})
	out, err := p.List(noContext, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if len(out) != 2 {
		t.Errorf("Expect combined registry output")
		return
	}
	if out[0] != a {
		t.Errorf("Unexpected registry at index 0")
	}
	if out[1] != b {
		t.Errorf("Unexpected registry at index 1")
	}
}
