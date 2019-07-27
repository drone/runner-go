// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package shell provides functions for converting shell commands
// to shell scripts.
package bash

import (
	"bytes"
	"fmt"
	"strings"
)

// Suffix provides the shell script suffix. For posix systems
// this value is an empty string.
const Suffix = ""

// Command returns the shell command and arguments.
func Command() (string, []string) {
	return "/bin/sh", []string{"-e"}
}

// Script converts a slice of individual shell commands to
// a posix-compliant shell script.
func Script(commands []string) string {
	buf := new(bytes.Buffer)
	fmt.Fprintln(buf)
	fmt.Fprintf(buf, optionScript)
	fmt.Fprintln(buf)
	for _, command := range commands {
		escaped := fmt.Sprintf("%q", command)
		escaped = strings.Replace(escaped, "$", `\$`, -1)
		buf.WriteString(fmt.Sprintf(
			traceScript,
			escaped,
			command,
		))
	}
	return buf.String()
}

// optionScript is a helper script this is added to the build
// to set shell options, in this case, to exit on error.
const optionScript = "set -e"

// traceScript is a helper script that is added to
// the build script to trace a command.
const traceScript = `
echo + %s
%s
`
