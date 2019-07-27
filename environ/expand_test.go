// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package environ

import (
	"os"
	"testing"
)

func TestExpand(t *testing.T) {
	defer func() {
		getenv = os.Getenv
	}()

	getenv = func(string) string {
		return "/bin:/usr/local/bin"
	}

	before := map[string]string{
		"USER": "root",
		"HOME": "/home/$USER", // does not expect
		"PATH": "/go/bin:$PATH",
	}

	after := Expand(before)
	if got, want := after["PATH"], "/go/bin:/bin:/usr/local/bin"; got != want {
		t.Errorf("Got PATH %q, want %q", got, want)
	}
	if got, want := after["USER"], "root"; got != want {
		t.Errorf("Got USER %q, want %q", got, want)
	}
	// only the PATH variable should expand. No other variables
	// should be expanded.
	if got, want := after["HOME"], "/home/$USER"; got != want {
		t.Errorf("Got HOME %q, want %q", got, want)
	}
}
