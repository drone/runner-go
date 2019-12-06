// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package console

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPlain(t *testing.T) {
	buf := new(bytes.Buffer)

	sess := New(false)
	w := sess.Stream(nil, nil, "clone").(*plain)
	w.base = buf
	w.Write([]byte("hello\nworld"))
	w.Close()

	got, want := buf.String(), "[clone:1] hello\n[clone:2] world\n"
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Invalid plain text log output")
		t.Log(diff)
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		before string
		after  []string
	}{
		{
			before: "hello world",
			after:  []string{"hello world"},
		},
		{
			before: "hello world\n",
			after:  []string{"hello world"},
		},
		{
			before: "hello\nworld\n",
			after:  []string{"hello", "world"},
		},
		{
			before: "hello\n\nworld\n",
			after:  []string{"hello", "", "world"},
		},
		{
			before: "\nhello\n\nworld\n",
			after:  []string{"", "hello", "", "world"},
		},
		{
			before: "\n",
			after:  []string{""},
		},
		{
			before: "\n\n",
			after:  []string{"", ""},
		},
	}
	for _, test := range tests {
		b := []byte(test.before)
		got, want := split(b), test.after
		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Invalid split")
			t.Log(diff)
		}
	}
}
