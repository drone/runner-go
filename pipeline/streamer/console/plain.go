// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package console

import (
	"fmt"
	"io"
	"strings"
)

// plain text line format with line number.
const plainf = "[%s:%d] %s\n"

type plain struct {
	base io.Writer
	name string
	seq  *sequence
}

func (w *plain) Write(b []byte) (int, error) {
	for _, part := range split(b) {
		fmt.Fprintf(w.base, plainf, w.name, w.seq.next(), part)
	}
	return len(b), nil
}

func (w *plain) Close() error {
	return nil
}

func split(b []byte) []string {
	s := string(b)
	s = strings.TrimSuffix(s, "\n")
	return strings.Split(s, "\n")
}
