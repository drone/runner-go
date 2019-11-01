// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package powershell

import (
	"reflect"
	"testing"
)

func TestCommands(t *testing.T) {
	cmd, args := Command()
	{
		got, want := cmd, "powershell"
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want command %v, got %v", want, got)
		}
	}
	{
		got, want := args, []string{
			"-noprofile",
			"-noninteractive",
			"-command",
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Want args %v, got %v", want, got)
		}
	}
}

func TestScript(t *testing.T) {
	got, want := Script([]string{"go build", "go test"}), exampleScript
	if got != want {
		t.Errorf("Want %q, got %q", want, got)
	}
}

var exampleScript = `
$erroractionpreference = "stop"

echo "+ go build"
go build
if ($LastExitCode -gt 0) { exit $LastExitCode }

echo "+ go test"
go test
if ($LastExitCode -gt 0) { exit $LastExitCode }
`
