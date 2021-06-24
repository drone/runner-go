// Copyright 2021 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package container

import "testing"

func TestIsRestrictedVolume(t *testing.T) {
	restrictedPaths := []string{
		"/",
		"../../../../../../../../../../../../var/run",
		"/var/run",
		"//var/run",
		"/var/run/",
		"/var/run/.",
		"/var//run/",
		"/var/run//",
		"/var/run/test/..",
		"/./var/run",
		"/var/./run",
	}

	allowedPaths := []string{
		"/drone",
		"/drone/var/run",
		"/development",
		"/var/lib",
		"/etc/ssh",
	}

	for _, path := range restrictedPaths {
		if result := IsRestrictedVolume(path); result != true {
			t.Errorf("Test failed for restricted path %q", path)
		}
	}

	for _, path := range allowedPaths {
		if result := IsRestrictedVolume(path); result != false {
			t.Errorf("Test failed for allowed path %q", path)
		}
	}
}
