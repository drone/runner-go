// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package manifest

import "testing"

type resourceImpl struct {
	Name, Kind, Type, Version string
}

func (r *resourceImpl) GetVersion() string { return r.Version }
func (r *resourceImpl) GetKind() string    { return r.Kind }
func (r *resourceImpl) GetType() string    { return r.Type }
func (r *resourceImpl) GetName() string    { return r.Name }

func TestLookup(t *testing.T) {
	want := &resourceImpl{Name: "default"}
	m := &Manifest{
		Resources: []Resource{want},
	}
	got, err := Lookup("default", m)
	if err != nil {
		t.Error(err)
	}
	if got != want {
		t.Errorf("Expect resource not found error")
	}
}

func TestLookupNotFound(t *testing.T) {
	m := &Manifest{
		Resources: []Resource{
			&Secret{
				Kind: "secret",
				Name: "password",
			},
		},
	}
	_, err := Lookup("default", m)
	if err == nil {
		t.Errorf("Expect resource not found error")
	}
}

func TestNameMatch(t *testing.T) {
	tests := []struct {
		a, b  string
		match bool
	}{
		{"a", "b", false},
		{"a", "a", true},
		{"", "default", true},
	}
	for _, test := range tests {
		got, want := isNameMatch(test.a, test.b), test.match
		if got != want {
			t.Errorf("Expect %q and %q match is %v", test.a, test.b, want)
		}
	}
}
