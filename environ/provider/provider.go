// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package provider provides environment variables to
// a pipeline.
package provider

import (
	"context"

	"github.com/drone/drone-go/drone"
)

// Request provides arguments for requesting a environment
// variables from an environment Provider.
type Request struct {
	Repo  *drone.Repo
	Build *drone.Build
}

// Provider is the interface that must be implemented by an
// environment provider.
type Provider interface {
	// List returns a list of environment variables.
	List(context.Context, *Request) (map[string]string, error)
}
