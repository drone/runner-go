// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Parity Public License
// that can be found in the LICENSE file.

// Package powershell provides functions for converting shell
// commands to powershell scripts.
package powershell

import (
	"bytes"
	"fmt"
	"strings"
)

// Suffix provides the shell script suffix.
const Suffix = ".ps1"

// Command returns the Powershell command and arguments.
func Command() (string, []string) {
	return "powershell", []string{
		"-noprofile",
		"-noninteractive",
		"-command",
	}
}

// Script converts a slice of individual shell commands to
// a powershell script.
func Script(commands []string, environ map[string]string) string {
	buf := new(bytes.Buffer)
	for k, v := range environ {
		fmt.Fprintln(buf)
		fmt.Fprintf(buf, exportScript, k, v)
	}
	fmt.Fprintln(buf)
	fmt.Fprintf(buf, optionScript)
	fmt.Fprintln(buf)
	for _, command := range commands {
		escaped := fmt.Sprintf("%q", "+ "+command)
		escaped = strings.Replace(escaped, "$", "`$", -1)
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
const optionScript = `$erroractionpreference = "stop"`

// exportScript is a helper script that is added to
// the build script to export environment variables.
const exportScript = "$Env:%s = %q"

// traceScript is a helper script that is added to
// the build script to trace a command.
const traceScript = `
echo %s
%s
`
