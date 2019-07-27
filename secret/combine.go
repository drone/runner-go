// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package secret

import (
	"context"

	"github.com/drone/drone-go/drone"
)

// Combine returns a new combined secret provider, capable of
// sourcing secrets from multiple providers.
func Combine(sources ...Provider) Provider {
	return &combined{sources}
}

type combined struct {
	sources []Provider
}

func (p *combined) Find(ctx context.Context, in *Request) (*drone.Secret, error) {
	for _, source := range p.sources {
		secret, err := source.Find(ctx, in)
		if err != nil {
			return nil, err
		}
		if secret == nil {
			continue
		}
		// if the secret object is not nil, but is empty
		// we should assume the secret service returned a
		// 204 no content, and proceed to the next service
		// in the chain.
		if secret.Data == "" {
			continue
		}
		return secret, nil
	}
	return nil, nil
}
