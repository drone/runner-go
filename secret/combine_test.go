// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Parity Public License
// that can be found in the LICENSE file.

package secret

import (
	"errors"
	"testing"

	"github.com/drone/drone-go/drone"
)

func TestCombine(t *testing.T) {
	secrets := []*drone.Secret{
		{Name: "docker_username", Data: "octocat"},
		{Name: "docker_password", Data: "correct-horse-battery-staple"},
	}
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPush},
	}
	service := Combine(Static(secrets[:1]), Static(secrets[1:]))
	secret, err := service.Find(noContext, args)
	if err != nil {
		t.Error(err)
		return
	}
	if secret != secrets[1] {
		t.Errorf("expect docker_password")
	}
}

func TestCombine_Error(t *testing.T) {
	args := &Request{
		Name:  "slack_token",
		Build: &drone.Build{Event: drone.EventPush},
	}
	want := errors.New("cannot find secret")
	mock := &mockProvider{err: want}
	service := Combine(mock)
	_, got := service.Find(noContext, args)
	if got != want {
		t.Errorf("expect error")
	}
}

func TestCombine_NotFound(t *testing.T) {
	secrets := []*drone.Secret{
		{Name: "docker_username", Data: "octocat"},
		{Name: "docker_password", Data: "correct-horse-battery-staple"},
	}
	args := &Request{
		Name:  "slack_token",
		Build: &drone.Build{Event: drone.EventPush},
	}
	service := Combine(Static(secrets))
	secret, err := service.Find(noContext, args)
	if err != nil {
		t.Error(err)
	}
	if secret != nil {
		t.Errorf("expect nil secret")
	}
}

func TestCombine_Empty(t *testing.T) {
	secrets := []*drone.Secret{
		{Name: "docker_username", Data: ""},
		{Name: "docker_password", Data: ""},
	}
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPush},
	}
	service := Combine(Static(secrets))
	secret, err := service.Find(noContext, args)
	if err != nil {
		t.Error(err)
	}
	if secret != nil {
		t.Errorf("expect nil secret")
	}
}
