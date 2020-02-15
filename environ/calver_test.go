// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package environ

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCalver(t *testing.T) {
	a := calversions("v19.1.0-beta.20190318")
	b := map[string]string{
		"DRONE_CALVER":             "19.1.0-beta.20190318",
		"DRONE_CALVER_MAJOR":       "19",
		"DRONE_CALVER_MAJOR_MINOR": "19.1",
		"DRONE_CALVER_MINOR":       "1",
		"DRONE_CALVER_MICRO":       "0",
		"DRONE_CALVER_SHORT":       "19.1.0",
		"DRONE_CALVER_MODIFIER":    "beta.20190318",
	}
	if diff := cmp.Diff(a, b); diff != "" {
		t.Errorf("Unexpected calver variables")
		t.Log(diff)
	}
}

func TestCalverAlternate(t *testing.T) {
	a := calversions("2019.01.0002")
	b := map[string]string{
		"DRONE_CALVER":             "2019.01.0002",
		"DRONE_CALVER_MAJOR_MINOR": "2019.01",
		"DRONE_CALVER_MAJOR":       "2019",
		"DRONE_CALVER_MINOR":       "01",
		"DRONE_CALVER_MICRO":       "0002",
		"DRONE_CALVER_SHORT":       "2019.01.0002",
	}
	if diff := cmp.Diff(a, b); diff != "" {
		t.Errorf("Unexpected calver variables")
		t.Log(diff)
	}
}

func TestCalver_Invalid(t *testing.T) {
	tests := []string{
		"1.2.3",
		"1.2",
		"1",
		"0.12",
		"0.12.1",
	}
	for _, s := range tests {
		envs := calversions(s)
		if len(envs) != 0 {
			t.Errorf("Expect invalid calversion: %s", s)
		}
	}
}

func TestCalverParser(t *testing.T) {
	tests := []struct {
		s string
		v *calver
	}{
		{"09.01.02", &calver{"09", "01", "02", ""}},
		{"2009.01.02", &calver{"2009", "01", "02", ""}},
		{"2009.1.2", &calver{"2009", "1", "2", ""}},
		{"09.1.2", &calver{"09", "1", "2", ""}},
		{"9.1.2", &calver{"9", "1", "2", ""}},
		{"v19.1.0-beta.20190318", &calver{"19", "1", "0", "beta.20190318"}},

		// invalid values
		{"foo.bar.baz", nil},
		{"foo.bar", nil},
		{"foo.1", nil},
		{"foo", nil},
		{"1", nil},
		{"1.foo", nil},
	}

	for _, test := range tests {
		got, want := parseCalver(test.s), test.v
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Unexpected calver %s", test.s)
			t.Log(diff)
		}
	}
}
