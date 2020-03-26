// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package provider

import (
	"context"
)

var noContext = context.Background()

type mockProvider struct {
	out []*Variable
	err error
}

func (p *mockProvider) List(context.Context, *Request) ([]*Variable, error) {
	return p.out, p.err
}
