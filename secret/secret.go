// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package secret provides secrets to a pipeline.
package secret

import (
	"context"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/manifest"
)

// Request provides arguments for requesting a secret from
// a secret Provider.
type Request struct {
	Name  string
	Repo  *drone.Repo
	Build *drone.Build
	Conf  *manifest.Manifest
}

// Provider is the interface that must be implemented by a
// secret provider.
type Provider interface {
	// Find finds and returns a requested secret.
	Find(context.Context, *Request) (*drone.Secret, error)
}
