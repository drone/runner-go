// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package provider

import (
	"context"
	"time"

	"github.com/drone/drone-go/plugin/environ"
	"github.com/drone/runner-go/logger"
)

// MultiExternal returns a new environment provider that
// is comprised of multiple external providers, and
// aggregates their results.
func MultiExternal(endpoints []string, token string, insecure bool) Provider {
	var sources []Provider
	for _, endpoint := range endpoints {
		sources = append(sources, External(
			endpoint, token, insecure,
		))
	}
	return Combine(sources...)
}

// External returns a new external environment variable
// provider. This provider makes an external API call to
// list and return environment variables.
func External(endpoint, token string, insecure bool) Provider {
	provider := &external{}
	if endpoint != "" {
		provider.client = environ.Client(endpoint, token, insecure)
	}
	return provider
}

type external struct {
	client environ.Plugin
}

func (p *external) List(ctx context.Context, in *Request) ([]*Variable, error) {
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

	req := &environ.Request{
		Repo:  *in.Repo,
		Build: *in.Build,
	}
	res, err := p.client.List(ctx, req)
	if err != nil {
		logger.WithError(err).Debug("environment: external: cannot get environment variable list")
		return nil, err
	}

	// if no error is returned and the list is empty,
	// this indicates the client returned No Content,
	// and we should exit with no credentials, but no error.
	if len(res) == 0 {
		logger.Trace("environment: external: environment variable list is empty")
		return nil, nil
	}

	logger.Trace("environment: external: environment variable list returned")

	var out []*Variable
	for _, v := range res {
		out = append(out, &Variable{
			Name: v.Name,
			Data: v.Data,
			Mask: v.Mask,
		})
	}
	return out, nil
}
