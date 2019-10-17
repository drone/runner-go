// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package registry

import (
	"context"

	"github.com/drone/drone-go/drone"
)

// Static returns a new static registry credential provider.
// The static secret provider finds and returns the static list
// of registry credentials.
func Static(registries []*drone.Registry) Provider {
	return &static{registries}
}

type static struct {
	registries []*drone.Registry
}

func (p *static) List(context.Context, *Request) ([]*drone.Registry, error) {
	return p.registries, nil
}
