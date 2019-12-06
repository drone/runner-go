// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package console

import (
	"fmt"
	"io"
)

// pretty line format with line number and coloring.
const prettyf = "\033[%s[%s:%d]\033[0m %s\n"

// available terminal colors
var colors = []string{
	"32m", // green
	"33m", // yellow
	"34m", // blue
	"35m", // magenta
	"36m", // cyan
}

type pretty struct {
	base  io.Writer
	color string
	name  string
	seq   *sequence
}

func (w *pretty) Write(b []byte) (int, error) {
	for _, part := range split(b) {
		fmt.Fprintf(w.base, prettyf, w.color, w.name, w.seq.next(), part)
	}
	return len(b), nil
}

func (w *pretty) Close() error {
	return nil
}
