// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package registry

import (
	"context"
	"time"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/registry"
	"github.com/drone/runner-go/logger"
)

// External returns a new external registry credentials
// provider. The external credentials provider makes an
// external API call to list and return credentials.
func External(endpoint, token string, insecure bool) Provider {
	provider := &external{}
	if endpoint != "" {
		provider.client = registry.Client(endpoint, token, insecure)
	}
	return provider
}

type external struct {
	client registry.Plugin
}

func (p *external) List(ctx context.Context, in *Request) ([]*drone.Registry, error) {
	if p.client == nil {
		return nil, nil
	}

	logger := logger.FromContext(ctx)

	// include a timeout to prevent an API call from
	// hanging the build process indefinitely. The
	// external service must return a request within
	// one minute.
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	req := &registry.Request{
		Repo:  *in.Repo,
		Build: *in.Build,
	}
	res, err := p.client.List(ctx, req)
	if err != nil {
		logger.WithError(err).Debug("registry: external: cannot get credentials")
		return nil, err
	}

	// if no error is returned and the list is empty,
	// this indicates the client returned No Content,
	// and we should exit with no credentials, but no error.
	if len(res) == 0 {
		logger.Trace("registry: external: credential list is empty")
		return nil, nil
	}

	for _, v := range res {
		logger.
			WithField("address", v.Address).
			WithField("username", v.Username).
			Trace("registry: external: received credentials")
	}

	return res, nil
}
