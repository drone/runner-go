// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// +build !windows

package shell

import "github.com/drone/runner-go/shell/bash"

// Suffix provides the shell script suffix.
const Suffix = ""

// Command returns the powershell command and arguments.
func Command() (string, []string) {
	return bash.Command()
}

// Script converts a slice of individual shell commands to
// a powershell script.
func Script(commands []string) string {
	return bash.Script(commands)
}
