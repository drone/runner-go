// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package provider

import "context"

// Static returns a new static environment variable provider.
// The static provider finds and returns the static list
// of static environment variables.
func Static(params map[string]string) Provider {
	return &static{params}
}

type static struct {
	params map[string]string
}

func (p *static) List(context.Context, *Request) ([]*Variable, error) {
	return ToSlice(p.params), nil
}
