// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package manifest

import (
	"testing"

	"github.com/buildkite/yaml"
	"github.com/google/go-cmp/cmp"
)

var mockSecretYaml = []byte(`
---
kind: secret
name: username

data: b2N0b2NhdA==
`)

var mockSecret = &Secret{
	Kind: "secret",
	Name: "username",
	Data: "b2N0b2NhdA==",
}

func TestSecretUnmarshal(t *testing.T) {
	a := new(Secret)
	b := mockSecret
	yaml.Unmarshal(mockSecretYaml, a)
	if diff := cmp.Diff(a, b); diff != "" {
		t.Error("Failed to parse secret")
		t.Log(diff)
	}
}

func TestSecretValidate(t *testing.T) {
	secret := new(Secret)

	secret.Data = "some-data"
	if err := secret.Validate(); err != nil {
		t.Error(err)
		return
	}

	secret.Get.Path = "secret/data/docker"
	if err := secret.Validate(); err != nil {
		t.Error(err)
		return
	}

	secret.Data = ""
	secret.Get.Path = ""
	if err := secret.Validate(); err == nil {
		t.Errorf("Expect invalid secret error")
	}
}
