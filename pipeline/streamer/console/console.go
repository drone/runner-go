// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package console provides a streamer that writes the pipeline
// output to stdout.
package console

import (
	"context"
	"io"
	"os"

	"github.com/drone/runner-go/pipeline"
)

var _ pipeline.Streamer = (*Console)(nil)

// Console implements a pipeline streamer that writes the
// pipeline logs to the console using os.Stdout.
type Console struct {
	seq *sequence
	col *sequence
	tty bool
}

// New returns a new console recorder.
func New(tty bool) *Console {
	return &Console{
		tty: tty,
		seq: new(sequence),
		col: new(sequence),
	}
}

// Stream returns an io.WriteCloser that prints formatted log
// lines to the console with step name, line number, and optional
// coloring.
func (s *Console) Stream(_ context.Context, _ *pipeline.State, name string) io.WriteCloser {
	if s.tty {
		return &pretty{
			base:  os.Stdout,
			color: colors[s.col.next()%len(colors)],
			name:  name,
			seq:   s.seq,
		}
	}
	return &plain{
		base: os.Stdout,
		name: name,
		seq:  s.seq,
	}
}
