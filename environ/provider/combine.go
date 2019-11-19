// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package provider

import (
	"context"

	"github.com/drone/runner-go/environ"
)

// Combine returns a new combined environment provider,
// capable of sourcing environment variables from multiple
// providers.
func Combine(sources ...Provider) Provider {
	return &combined{sources}
}

type combined struct {
	sources []Provider
}

func (p *combined) List(ctx context.Context, in *Request) (map[string]string, error) {
	out := map[string]string{}
	for _, source := range p.sources {
		got, err := source.List(ctx, in)
		if err != nil {
			return nil, err
		}
		if got != nil {
			out = environ.Combine(got, out)
		}
	}
	return out, nil
}
