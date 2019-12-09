// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package replacer

import (
	"bytes"
	"io"
	"testing"

	"github.com/drone/runner-go/pipeline/runtime/driver"
)

func TestReplace(t *testing.T) {
	secrets := []driver.Secret{
		&secret{Name: "DOCKER_USERNAME", Data: "octocat", Mask: false},
		&secret{Name: "DOCKER_PASSWORD", Data: "correct-horse-batter-staple", Mask: true},
		&secret{Name: "DOCKER_EMAIL", Data: "", Mask: true},
	}

	buf := new(bytes.Buffer)
	w := New(&nopCloser{buf}, secrets)
	w.Write([]byte("username octocat password correct-horse-batter-staple"))
	w.Close()

	if got, want := buf.String(), "username octocat password ******"; got != want {
		t.Errorf("Want masked string %s, got %s", want, got)
	}
}

// this test verifies that if there are no secrets to scan and
// mask, the io.WriteCloser is returned as-is.
func TestReplaceNone(t *testing.T) {
	secrets := []driver.Secret{
		&secret{Name: "DOCKER_USERNAME", Data: "octocat", Mask: false},
		&secret{Name: "DOCKER_PASSWORD", Data: "correct-horse-batter-staple", Mask: false},
	}

	buf := new(bytes.Buffer)
	w := &nopCloser{buf}
	r := New(w, secrets)
	if w != r {
		t.Errorf("Expect buffer returned with no replacer")
	}
}

type nopCloser struct {
	io.Writer
}

func (*nopCloser) Close() error {
	return nil
}

type secret struct {
	Name string
	Data string
	Mask bool
}

func (s *secret) GetName() string  { return s.Name }
func (s *secret) GetValue() string { return s.Data }
func (s *secret) IsMasked() bool   { return s.Mask }
