// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package registry

import (
	"context"

	"github.com/drone/drone-go/drone"
)

// Combine returns a new combined registry provider, capable of
// sourcing registry credentials from multiple providers.
func Combine(sources ...Provider) Provider {
	return &combined{sources}
}

type combined struct {
	sources []Provider
}

func (p *combined) List(ctx context.Context, in *Request) ([]*drone.Registry, error) {
	var out []*drone.Registry
	for _, source := range p.sources {
		list, err := source.List(ctx, in)
		if err != nil {
			return nil, err
		}
		out = append(out, list...)
	}
	return out, nil
}
