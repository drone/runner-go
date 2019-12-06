// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package remote provides a reporter and streamer that sends the
// pipeline status and logs to the central server.
package remote

import (
	"context"
	"io"

	"github.com/drone/runner-go/client"
	"github.com/drone/runner-go/internal"
	"github.com/drone/runner-go/livelog"
	"github.com/drone/runner-go/pipeline"
)

var _ pipeline.Reporter = (*Remote)(nil)
var _ pipeline.Streamer = (*Remote)(nil)

// Remote implements a pipeline reporter that reports state
// changes and results to a remote server instance.
type Remote struct {
	client client.Client
}

// New returns a remote reporter.
func New(client client.Client) *Remote {
	return &Remote{
		client: client,
	}
}

// ReportStage reports the stage status.
func (s *Remote) ReportStage(ctx context.Context, state *pipeline.State) error {
	state.Lock()
	src := state.Stage
	cpy := internal.CloneStage(src)
	state.Unlock()
	err := s.client.Update(ctx, cpy)
	if err == nil {
		state.Lock()
		internal.MergeStage(cpy, src)
		state.Unlock()
	}
	return err
}

// ReportStep reports the step status.
func (s *Remote) ReportStep(ctx context.Context, state *pipeline.State, name string) error {
	src := state.Find(name)
	state.Lock()
	cpy := internal.CloneStep(src)
	state.Unlock()
	err := s.client.UpdateStep(ctx, cpy)
	if err == nil {
		state.Lock()
		internal.MergeStep(cpy, src)
		state.Unlock()
	}
	return err
}

// Stream returns an io.WriteCloser to stream the stdout
// and stderr of the pipeline step to the server.
func (s *Remote) Stream(ctx context.Context, state *pipeline.State, name string) io.WriteCloser {
	src := state.Find(name)
	return livelog.New(s.client, src.ID)
}
