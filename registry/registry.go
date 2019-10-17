// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package registry provides registry credentials used
// to pull private images from a registry.
package registry

import (
	"context"

	"github.com/drone/drone-go/drone"
)

// Request provides arguments for requesting a secret from
// a secret Provider.
type Request struct {
	Repo  *drone.Repo
	Build *drone.Build
}

// Provider is the interface that must be implemented by a
// registry provider.
type Provider interface {
	// Find finds and returns a list of registry credentials.
	List(context.Context, *Request) ([]*drone.Registry, error)
}
