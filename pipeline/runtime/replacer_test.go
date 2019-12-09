// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package runtime

import (
	"bytes"
	"io"
	"testing"
)

func TestReplace(t *testing.T) {
	secrets := []Secret{
		&mockSecret{Name: "DOCKER_USERNAME", Data: "octocat", Mask: false},
		&mockSecret{Name: "DOCKER_PASSWORD", Data: "correct-horse-batter-staple", Mask: true},
		&mockSecret{Name: "DOCKER_EMAIL", Data: "", Mask: true},
	}

	buf := new(bytes.Buffer)
	w := newReplacer(&nopCloser{buf}, secrets)
	w.Write([]byte("username octocat password correct-horse-batter-staple"))
	w.Close()

	if got, want := buf.String(), "username octocat password ******"; got != want {
		t.Errorf("Want masked string %s, got %s", want, got)
	}
}

// this test verifies that if there are no secrets to scan and
// mask, the io.WriteCloser is returned as-is.
func TestReplaceNone(t *testing.T) {
	secrets := []Secret{
		&mockSecret{Name: "DOCKER_USERNAME", Data: "octocat", Mask: false},
		&mockSecret{Name: "DOCKER_PASSWORD", Data: "correct-horse-batter-staple", Mask: false},
	}

	buf := new(bytes.Buffer)
	w := &nopCloser{buf}
	r := newReplacer(w, secrets)
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

type mockSecret struct {
	Name string
	Data string
	Mask bool
}

func (s *mockSecret) GetName() string  { return s.Name }
func (s *mockSecret) GetValue() string { return s.Data }
func (s *mockSecret) IsMasked() bool   { return s.Mask }
