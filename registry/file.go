// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package registry

import (
	"context"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/logger"
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

func (p *file) List(ctx context.Context, _ *Request) ([]*drone.Registry, error) {
	if p.path == "" {
		return nil, nil
	}

	logger := logger.FromContext(ctx)
	logger.WithField("path", p.path).
		Trace("registry: file: parsing credentials file")

	// load the registry credentials from the file.
	res, err := auths.ParseFile(p.path)
	if err != nil {
		logger.WithError(err).
			Debug("registry: file: cannot parse credentials file")
		return nil, err
	}

	// if no error is returned and the list is empty,
	// this indicates the client returned No Content,
	// and we should exit with no credentials, but no error.
	if len(res) == 0 {
		logger.Trace("registry: file: credential list is empty")
		return nil, nil
	}

	for _, v := range res {
		logger.
			WithField("address", v.Address).
			WithField("username", v.Username).
			Trace("registry: file: received credentials")
	}

	return res, err
}
