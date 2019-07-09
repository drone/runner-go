// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Parity Public License
// that can be found in the LICENSE file.

// +build !windows

// Package shell provides functions for converting shell commands
// to shell scripts.
package shell

import (
	"bytes"
	"fmt"
	"strings"
)

// Suffix provides the shell script suffix.
const Suffix = ".sh"

// Command returns the shell command and arguments.
func Command() (string, []string) {
	return "/bin/sh", []string{"-e"}
}

// Script converts a slice of individual shell commands to
// a posix-compliant shell script.
func Script(commands []string) string {
	var buf bytes.Buffer
	for _, command := range commands {
		escaped := fmt.Sprintf("%q", command)
		escaped = strings.Replace(escaped, "$", `\$`, -1)
		buf.WriteString(fmt.Sprintf(
			traceScript,
			escaped,
			command,
		))
	}
	return fmt.Sprintf(
		buildScript,
		buf.String(),
	)
}

// buildScript is a helper script this is added to the build
// to prepare the environment and execute the build commands.
const buildScript = `
set -e
%s
`

// traceScript is a helper script that is added to
// the build script to trace a command.
const traceScript = `
echo + %s
%s
`
