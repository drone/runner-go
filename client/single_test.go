// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package client

import (
	"context"
	"errors"
	"testing"

	"github.com/drone/drone-go/drone"
)

var noContext = context.Background()

func TestSingleFlight(t *testing.T) {
	mock := &mockRequestClient{
		out: &drone.Stage{},
		err: errors.New("some random error"),
	}
	client := NewSingleFlight("", "", false)
	client.Client = mock
	out, err := client.Request(noContext, nil)
	if got, want := out, mock.out; got != want {
		t.Errorf("Expect stage returned from request")
	}
	if got, want := err, mock.err; got != want {
		t.Errorf("Expect error returned from request")
	}
}

func TestSingleFlightPanic(t *testing.T) {
	mock := &mockRequestClientPanic{}
	client := NewSingleFlight("", "", false)
	client.Client = mock

	defer func() {
		if recover() != nil {
			t.Errorf("Expect Request to recover from panic")
		}
		client.mu.Lock()
		client.mu.Unlock()
	}()

	client.Request(noContext, nil)
}

func TestSingleFlightCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(noContext)
	cancel()
	client := NewSingleFlight("", "", false)
	client.Request(ctx, nil)
}

// mock client that returns a static stage and error
// from the request method.
type mockRequestClient struct {
	Client

	out *drone.Stage
	err error
}

func (m *mockRequestClient) Request(ctx context.Context, args *Filter) (*drone.Stage, error) {
	return m.out, m.err
}

// mock client that returns panics when the request
// method is invoked.
type mockRequestClientPanic struct {
	Client
}

func (m *mockRequestClientPanic) Request(ctx context.Context, args *Filter) (*drone.Stage, error) {
	panic("method not implemented")
}
