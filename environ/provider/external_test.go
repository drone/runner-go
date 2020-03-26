// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package provider

import (
	"context"
	"errors"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/environ"
	"github.com/google/go-cmp/cmp"
)

func TestExternal(t *testing.T) {
	req := &Request{
		Build: &drone.Build{Event: drone.EventPush},
		Repo:  &drone.Repo{Private: false},
	}
	res := []*environ.Variable{
		{
			Name: "a",
			Data: "b",
			Mask: true,
		},
	}

	want := []*Variable{{Name: "a", Data: "b", Mask: true}}
	provider := External("http://localhost", "secret", false)
	provider.(*external).client = &mockPlugin{out: res}
	got, err := provider.List(noContext, req)
	if err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf(diff)
	}
}

// This test verifies that if the remote API call to the
// external plugin returns an error, the provider returns the
// error to the caller.
func TestExternal_ClientError(t *testing.T) {
	req := &Request{
		Build: &drone.Build{Event: drone.EventPush},
		Repo:  &drone.Repo{Private: false},
	}
	want := errors.New("Not Found")
	provider := External("http://localhost", "secret", false)
	provider.(*external).client = &mockPlugin{err: want}
	_, got := provider.List(noContext, req)
	if got != want {
		t.Errorf("Want error %s, got %s", want, got)
	}
}

// This test verifies that if no endpoint is configured the
// provider exits immediately and returns a nil slice and nil
// error.
func TestExternal_NoEndpoint(t *testing.T) {
	provider := External("", "", false)
	res, err := provider.List(noContext, nil)
	if err != nil {
		t.Errorf("Expect nil error, provider disabled")
	}
	if res != nil {
		t.Errorf("Expect nil secret, provider disabled")
	}
}

// This test verifies that nil credentials and a nil error
// are returned if the registry endpoint returns no content.
func TestExternal_NotFound(t *testing.T) {
	req := &Request{
		Repo:  &drone.Repo{},
		Build: &drone.Build{},
	}
	provider := External("http://localhost", "secret", false)
	provider.(*external).client = &mockPlugin{}
	res, err := provider.List(noContext, req)
	if err != nil {
		t.Errorf("Expect nil error, registry list empty")
	}
	if res != nil {
		t.Errorf("Expect nil registry credentials")
	}
}

// This test verifies that multiple external providers
// are combined into a single provider that concatenates
// the results.
func TestMultiExternal(t *testing.T) {
	provider := MultiExternal([]string{"https://foo", "https://bar"}, "correct-horse-batter-staple", true).(*combined)
	if len(provider.sources) != 2 {
		t.Errorf("Expect two provider sources")
	}
}

type mockPlugin struct {
	out []*environ.Variable
	err error
}

func (m *mockPlugin) List(context.Context, *environ.Request) ([]*environ.Variable, error) {
	return m.out, m.err
}
