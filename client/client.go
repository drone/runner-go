// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package client provides a client for using the runner API.
package client

import (
	"context"
	"errors"

	"github.com/drone/drone-go/drone"
)

// V1 is version 1 of the runner API
const V1 = "application/vnd.drone.runner.v1+json"

// ErrOptimisticLock is returned by if the struct being
// modified has a Version field and the value is not equal
// to the current value in the database
var ErrOptimisticLock = errors.New("Optimistic Lock Error")

type (
	// Filter is used to filter the builds that are pulled
	// from the queue.
	Filter struct {
		Kind    string            `json:"kind"`
		Type    string            `json:"type"`
		OS      string            `json:"os"`
		Arch    string            `json:"arch"`
		Variant string            `json:"variant"`
		Kernel  string            `json:"kernel"`
		Labels  map[string]string `json:"labels,omitempty"`
	}

	// File represents a file from the version control
	// repository. It is used by the runner to provide the
	// yaml configuration file to the runner.
	File struct {
		Data []byte
		Hash []byte
	}

	// Context provides the runner with the build context and
	// includes all environment data required to execute the
	// build.
	Context struct {
		Build   *drone.Build    `json:"build"`
		Config  *File           `json:"config"`
		Netrc   *drone.Netrc    `json:"netrc"`
		Repo    *drone.Repo     `json:"repository"`
		Secrets []*drone.Secret `json:"secrets"`
		System  *drone.System   `json:"system"`
	}
)

// A Client manages communication with the runner.
type Client interface {
	// Join notifies the server the runner is joining the cluster.
	Join(ctx context.Context, machine string) error

	// Leave notifies the server the runner is leaving the cluster.
	Leave(ctx context.Context, machine string) error

	// Ping sends a ping message to the server to test connectivity.
	Ping(ctx context.Context, machine string) error

	// Request requests the next available build stage for execution.
	Request(ctx context.Context, args *Filter) (*drone.Stage, error)

	// Accept accepts the build stage for execution.
	Accept(ctx context.Context, stage *drone.Stage) error

	// Detail gets the build stage details for execution.
	Detail(ctx context.Context, stage *drone.Stage) (*Context, error)

	// Update updates the build stage.
	Update(ctxt context.Context, step *drone.Stage) error

	// UpdateStep updates the build step.
	UpdateStep(ctx context.Context, stage *drone.Step) error

	// Watch watches for build cancellation requests.
	Watch(ctx context.Context, stage int64) (bool, error)

	// Batch batch writes logs to the build logs.
	Batch(ctx context.Context, step int64, lines []*drone.Line) error

	// Upload uploads the full logs to the server.
	Upload(ctx context.Context, step int64, lines []*drone.Line) error
}
