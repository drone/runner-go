// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package livelog

import (
	"context"
	"testing"
	"time"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/client"

	"github.com/google/go-cmp/cmp"
)

func TestLineWriterSingle(t *testing.T) {
	client := new(mockClient)
	w := New(client, 1)
	w.SetInterval(time.Duration(0))
	w.num = 4
	w.Write([]byte("foo\nbar\n"))

	a := w.pending
	b := []*drone.Line{
		{Number: 4, Message: "foo\n"},
		{Number: 5, Message: "bar\n"},
		{Number: 6, Message: ""},
	}
	if diff := cmp.Diff(a, b); diff != "" {
		t.Fail()
		t.Log(diff)
	}

	w.Close()
	a = client.uploaded
	if diff := cmp.Diff(a, b); diff != "" {
		t.Fail()
		t.Log(diff)
	}

	if len(w.pending) > 0 {
		t.Errorf("Expect empty buffer")
	}
}

func TestLineWriterLimit(t *testing.T) {
	client := new(mockClient)
	w := New(client, 0)
	if got, want := w.limit, defaultLimit; got != want {
		t.Errorf("Expect default buffer limit %d, got %d", want, got)
	}
	w.SetLimit(6)
	if got, want := w.limit, 6; got != want {
		t.Errorf("Expect custom buffer limit %d, got %d", want, got)
	}

	w.Write([]byte("foo"))
	w.Write([]byte("bar"))
	w.Write([]byte("baz"))

	if got, want := w.size, 6; got != want {
		t.Errorf("Expect buffer size %d, got %d", want, got)
	}

	a := w.history
	b := []*drone.Line{
		{Number: 1, Message: "bar"},
		{Number: 2, Message: "baz"},
	}
	if diff := cmp.Diff(a, b); diff != "" {
		t.Fail()
		t.Log(diff)
	}
}

type mockClient struct {
	client.Client
	lines    []*drone.Line
	uploaded []*drone.Line
}

func (m *mockClient) Batch(ctx context.Context, id int64, lines []*drone.Line) error {
	m.lines = append(m.lines, lines...)
	return nil
}

func (m *mockClient) Upload(ctx context.Context, id int64, lines []*drone.Line) error {
	m.uploaded = lines
	return nil
}
