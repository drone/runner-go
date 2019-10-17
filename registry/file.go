// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package registry

import (
	"context"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/registry/auths"
)

// File returns a new registry credential provider that
// parses and returns credentials from the Docker user
// configuration file.
func File(path string) Provider {
	return &file{path}
}

type file struct {
	path string
}

func (p *file) List(context.Context, *Request) ([]*drone.Registry, error) {
	if p.path == "" {
		return nil, nil
	}
	return auths.ParseFile(p.path)
}
