// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package pipeline

import (
	"context"
	"io"
)

// A Streamer streams the pipeline logs.
type Streamer interface {
	// Stream returns an io.WriteCloser to stream the stdout
	// and stderr of the pipeline step.
	Stream(context.Context, *State, string) io.WriteCloser
}

// NopStreamer returns a noop streamer.
func NopStreamer() Streamer {
	return new(nopStreamer)
}

type nopStreamer struct{}

func (*nopStreamer) Stream(context.Context, *State, string) io.WriteCloser {
	return new(nopWriteCloser)
}

type nopWriteCloser struct{}

func (*nopWriteCloser) Close() error                { return nil }
func (*nopWriteCloser) Write(p []byte) (int, error) { return len(p), nil }
