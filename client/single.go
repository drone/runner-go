// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package client

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/drone/drone-go/drone"
)

var _ Client = (*SingleFlight)(nil)

// SingleFlight wraps a Client and limits to a single in-flight
// request to pull items from the queue.
type SingleFlight struct {
	Client
	mu sync.Mutex
}

// NewSingleFlight returns a Client that is limited to a single in-flight
// request to pull items from the queue.
func NewSingleFlight(endpoint, secret string, skipverify bool) *SingleFlight {
	return &SingleFlight{Client: New(endpoint, secret, skipverify)}
}

// Request requests the next available build stage for execution.
func (t *SingleFlight) Request(ctx context.Context, args *Filter) (*drone.Stage, error) {
	// if the context is canceled there is no need to make
	// the request and we can exit early.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	// if is critical to unlock the mutex when the function
	// exits. although a panic is unlikely it is critical that
	// we recover from the panic to avoid deadlock.
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
		t.mu.Unlock()
	}()
	// lock the mutex to ensure only a single in-flight
	// request to request a resource from the server queue.
	t.mu.Lock()
	return t.Client.Request(ctx, args)
}
