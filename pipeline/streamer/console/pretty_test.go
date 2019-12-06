// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package console

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPretty(t *testing.T) {
	buf := new(bytes.Buffer)

	sess := New(true)
	w := sess.Stream(nil, nil, "clone").(*pretty)
	w.base = buf
	w.Write([]byte("hello\nworld"))
	w.Close()

	got, want := buf.String(), "\x1b[33m[clone:1]\x1b[0m hello\n\x1b[33m[clone:2]\x1b[0m world\n"
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Invalid plain text log output")
		t.Log(diff)
	}
}
