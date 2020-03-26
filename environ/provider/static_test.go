// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package provider

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStatic(t *testing.T) {
	in := map[string]string{"a": "b"}

	got, err := Static(in).List(noContext, nil)
	if err != nil {
		t.Error(err)
		return
	}

	want := []*Variable{
		{
			Name: "a",
			Data: "b",
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf(diff)
	}
}
