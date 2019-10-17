// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package registry

import (
	"context"

	"github.com/drone/drone-go/drone"
)

var noContext = context.Background()

type mockProvider struct {
	out []*drone.Registry
	err error
}

func (p *mockProvider) List(context.Context, *Request) ([]*drone.Registry, error) {
	return p.out, p.err
}
