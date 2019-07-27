// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package secret

import (
	"context"
	"errors"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/secret"
	"github.com/drone/runner-go/manifest"
)

func TestExternal(t *testing.T) {
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPush},
		Repo:  &drone.Repo{Private: false},
		Conf: &manifest.Manifest{
			Resources: []manifest.Resource{
				&manifest.Secret{
					Name: "docker_password",
					Get: manifest.SecretGet{
						Path: "docker",
						Name: "password",
					},
				},
			},
		},
	}
	extern := &drone.Secret{Name: "docker_password", Data: "correct-horse-battery-staple"}
	provider := External("http://localhost", "secret", false)
	provider.(*external).client = &mockPlugin{sec: extern}
	result, err := provider.Find(noContext, args)
	if err != nil {
		t.Errorf("Expect nil error, secret found")
	}
	if extern == nil {
		t.Errorf("Expect secret returned, got nil")
		return
	}
	if got, want := result.Name, extern.Name; got != want {
		t.Errorf("Want secret name %s, got %s", want, got)
	}
	if got, want := result.Data, extern.Data; got != want {
		t.Errorf("Want secret value %s, got %s", want, got)
	}
}

// This test verifies that if no endpoint is configured the
// provider exits immediately and returns a nil secret and nil
// error.
func TestExternal_NoEndpoint(t *testing.T) {
	provider := External("", "", false)
	sec, err := provider.Find(noContext, nil)
	if err != nil {
		t.Errorf("Expect nil error, provider disabled")
	}
	if sec != nil {
		t.Errorf("Expect nil secret, provider disabled")
	}
}

// This test verifies that a nil secret and nil error are
// returned if no corresponding secret resource entry is found
// in the manifest. No error is returned because Not Found is
// not considered an error.
func TestExternal_NotFound(t *testing.T) {
	args := &Request{
		Name:  "docker_username",
		Build: &drone.Build{Event: drone.EventPush},
		Conf: &manifest.Manifest{
			Resources: []manifest.Resource{
				&manifest.Signature{
					Name: "signature",
					Hmac: "<signature>",
				},
			},
		},
	}
	provider := External("http://localhost", "secret", false)
	sec, err := provider.Find(noContext, args)
	if err != nil {
		t.Errorf("Expect nil error, secret not found")
	}
	if sec != nil {
		t.Errorf("Expect nil secret, secret not found")
	}
}

// This test verifies that a nil secret and nil error are
// returned if no matching secret resources is found in the
// manifest.
func TestExternal_NoMatch(t *testing.T) {
	args := &Request{
		Name:  "docker_username",
		Build: &drone.Build{Event: drone.EventPush},
		Conf: &manifest.Manifest{
			Resources: []manifest.Resource{
				&manifest.Secret{
					Name: "docker_password",
				},
			},
		},
	}
	provider := External("http://localhost", "secret", false)
	sec, err := provider.Find(noContext, args)
	if err != nil {
		t.Errorf("Expect nil error, no matching secret")
	}
	if sec != nil {
		t.Errorf("Expect nil secret, no matching secret")
	}
}

// This test verifies that if a secret is requested, and secret
// resource exists in the manifest but with no path or name, a
// nil secret and nil error are returned.
func TestExternal_NoPath(t *testing.T) {
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPush},
		Conf: &manifest.Manifest{
			Resources: []manifest.Resource{
				&manifest.Secret{
					Name: "docker_password",
					Get: manifest.SecretGet{
						Path: "",
						Name: "",
					},
				},
			},
		},
	}
	provider := External("http://localhost", "secret", false)
	sec, err := provider.Find(noContext, args)
	if err != nil {
		t.Errorf("Expect nil error, no path")
	}
	if sec != nil {
		t.Errorf("Expect nil secret, no path")
	}
}

// This test verifies that if the remote API call to the
// external plugin returns an error, the provider returns the
// error to the caller.
func TestExternal_ClientError(t *testing.T) {
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPush},
		Repo:  &drone.Repo{Private: false},
		Conf: &manifest.Manifest{
			Resources: []manifest.Resource{
				&manifest.Secret{
					Name: "docker_password",
					Get: manifest.SecretGet{
						Path: "docker",
						Name: "password",
					},
				},
			},
		},
	}
	want := errors.New("not found")
	provider := External("http://localhost", "secret", false)
	provider.(*external).client = &mockPlugin{err: want}
	_, got := provider.Find(noContext, args)
	if got != want {
		t.Errorf("Expect error returned from client")
	}
}

// This test verifies that if a secret with an emtpy value is
// returned from the external plugin, a nil secret and nil error
// are returned by the provider.
func TestExternal_EmptySecret(t *testing.T) {
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPush},
		Repo:  &drone.Repo{Private: false},
		Conf: &manifest.Manifest{
			Resources: []manifest.Resource{
				&manifest.Secret{
					Name: "docker_password",
					Get: manifest.SecretGet{
						Path: "docker",
						Name: "password",
					},
				},
			},
		},
	}
	res := &drone.Secret{Name: "docker_password", Data: ""}
	provider := External("http://localhost", "secret", false)
	provider.(*external).client = &mockPlugin{sec: res}
	sec, err := provider.Find(noContext, args)
	if err != nil {
		t.Errorf("Expect nil error, secret not found")
	}
	if sec != nil {
		t.Errorf("Expect nil secret, secret not found")
	}
}

type mockPlugin struct {
	sec *drone.Secret
	err error
}

func (m *mockPlugin) Find(context.Context, *secret.Request) (*drone.Secret, error) {
	return m.sec, m.err
}
