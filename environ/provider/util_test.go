// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package provider

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestToMap(t *testing.T) {
	in := []*Variable{
		{
			Name: "foo",
			Data: "bar",
		},
	}
	want := map[string]string{
		"foo": "bar",
	}
	got := ToMap(in)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Log(diff)
		t.Errorf("Unexpected map value")
	}
}

func TestFromMap(t *testing.T) {
	in := map[string]string{
		"foo": "bar",
	}
	want := []*Variable{
		{
			Name: "foo",
			Data: "bar",
		},
	}
	got := ToSlice(in)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Log(diff)
		t.Errorf("Unexpected variable list")
	}
}

func TestFilterMasked(t *testing.T) {
	in := []*Variable{
		{
			Name: "foo",
			Data: "bar",
			Mask: false,
		},
		{
			Name: "baz",
			Data: "qux",
			Mask: true,
		},
	}
	want := in[1:]
	got := FilterMasked(in)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Log(diff)
		t.Errorf("Unexpected variable list")
	}
}

func TestFilterUnmasked(t *testing.T) {
	in := []*Variable{
		{
			Name: "foo",
			Data: "bar",
			Mask: true,
		},
		{
			Name: "baz",
			Data: "qux",
			Mask: false,
		},
	}
	want := in[1:]
	got := FilterUnmasked(in)
	if diff := cmp.Diff(want, got); diff != "" {
		t.Log(diff)
		t.Errorf("Unexpected variable list")
	}
}
