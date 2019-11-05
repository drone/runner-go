// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package auths

import (
	"strings"
	"testing"

	"github.com/drone/drone-go/drone"
)

func TestEncode(t *testing.T) {
	endpoint := "docker.io"
	username := "octocat"
	password := "correct-horse-battery-staple"
	registry := &drone.Registry{
		Username: username,
		Password: password,
		Address:  endpoint,
	}
	got := Encode(registry, registry, registry)
	want := `{"auths":{"docker.io":{"auth":"b2N0b2NhdDpjb3JyZWN0LWhvcnNlLWJhdHRlcnktc3RhcGxl"}}}`
	if strings.TrimSpace(got) != strings.TrimSpace(want) {
		t.Errorf("Unexpected encoding: %q want %q", got, want)
	}
}
