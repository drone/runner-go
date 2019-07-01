// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Parity Public License
// that can be found in the LICENSE file.

package secret

import (
	"context"
	"testing"

	"github.com/drone/drone-go/drone"
)

var noContext = context.Background()

func TestStatic(t *testing.T) {
	secrets := []*drone.Secret{
		{Name: "docker_username"},
		{Name: "docker_password"},
	}
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPush},
	}
	service := Static(secrets)
	secret, err := service.Find(noContext, args)
	if err != nil {
		t.Error(err)
		return
	}
	if secret != secrets[1] {
		t.Errorf("expect docker_password")
	}
}

func TestStaticVars(t *testing.T) {
	secrets := map[string]string{
		"docker_username": "octocat",
		"docker_password": "correct-horse-battery-staple",
	}
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPush},
	}
	service := StaticVars(secrets)
	secret, err := service.Find(noContext, args)
	if err != nil {
		t.Error(err)
		return
	}
	if secret.Data != secrets["docker_password"] {
		t.Errorf("expect docker_password")
	}
}

func TestStaticNotFound(t *testing.T) {
	secrets := []*drone.Secret{
		{Name: "docker_username"},
		{Name: "docker_password"},
	}
	args := &Request{
		Name:  "slack_token",
		Build: &drone.Build{Event: drone.EventPush},
	}
	service := Static(secrets)
	secret, err := service.Find(noContext, args)
	if err != nil {
		t.Error(err)
		return
	}
	if secret != nil {
		t.Errorf("Expect secret not found")
	}
}

func TestStaticPullRequestDisabled(t *testing.T) {
	secrets := []*drone.Secret{
		{Name: "docker_username"},
		{Name: "docker_password", PullRequest: false},
	}
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPullRequest},
	}
	service := Static(secrets)
	secret, err := service.Find(noContext, args)
	if err != nil {
		t.Error(err)
		return
	}
	if secret != nil {
		t.Errorf("Expect secret not found")
	}
}

func TestStaticPullRequestEnabled(t *testing.T) {
	secrets := []*drone.Secret{
		{Name: "docker_username"},
		{Name: "docker_password", PullRequest: true},
	}
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPullRequest},
	}
	service := Static(secrets)
	secret, err := service.Find(noContext, args)
	if err != nil {
		t.Error(err)
		return
	}
	if err != nil {
		t.Error(err)
		return
	}
	if secret != secrets[1] {
		t.Errorf("expect docker_username")
	}
}
