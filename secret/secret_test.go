// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package secret

import (
	"context"

	"github.com/drone/drone-go/drone"
)

type mockProvider struct {
	sec *drone.Secret
	err error
}

func (p *mockProvider) Find(context.Context, *Request) (*drone.Secret, error) {
	return p.sec, p.err
}
