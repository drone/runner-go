// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package livelog

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/client"

	"github.com/google/go-cmp/cmp"
)

var optNoTS = cmpopts.IgnoreFields(drone.Line{}, "Timestamp")

// TestLineWriterClose tests if closing the Writer triggers streaming of all pending lines and upload of the full history.
func TestLineWriterClose(t *testing.T) {
	c := newMockClient()
	w := New(c, 1)
	w.SetInterval(time.Hour) // make sure it does not stream data
	w.lineList.lineCnt = 4   // new lines are starting from the Number=4

	w.Write([]byte("foo\nbar\n"))

	a := w.lineList.peekPending()
	b := []*drone.Line{
		{Number: 4, Message: "foo\n"},
		{Number: 5, Message: "bar\n"},
		{Number: 6, Message: ""},
	}
	if diff := cmp.Diff(a, b, optNoTS); diff != "" {
		t.Fail()
		t.Log(diff)
	}

	if len(c.uploaded) != 0 || len(c.lines) != 0 {
		t.Error("Expected nothing has been streamed or uploaded")
	}

	w.Close()

	if diff := cmp.Diff(c.lines, b, optNoTS); diff != "" {
		t.Error("Expected all output has been streamed")
		t.Log(diff)
	}

	if diff := cmp.Diff(c.uploaded, b, optNoTS); diff != "" {
		t.Error("Expected all output has been uploaded")
		t.Log(diff)
	}

	if len(w.lineList.peekPending()) > 0 {
		t.Errorf("Expect empty buffer")
	}
}

// TestLineWriterStreaming tests if streaming is done correctly through a client.
func TestLineWriterStreaming(t *testing.T) {
	c := newMockClient()
	w := New(c, 1)
	w.SetInterval(time.Nanosecond)

	w.Write([]byte("foo"))
	c.waitUpload()

	var a, b []*drone.Line

	a = w.lineList.peekPending()
	if len(a) != 0 {
		t.Errorf("Expected that all lines are uploaded, but there are still %d pending lines", len(a))
	}

	a = c.lines
	b = []*drone.Line{{Number: 0, Message: "foo"}}
	if diff := cmp.Diff(a, b, optNoTS); diff != "" {
		t.Fail()
		t.Log(diff)
	}

	w.Write([]byte("bar"))
	c.waitUpload()

	a = w.lineList.peekPending()
	if len(a) != 0 {
		t.Errorf("Expected that all lines are uploaded, but there are still %d pending lines", len(a))
	}

	a = c.lines
	b = []*drone.Line{{Number: 0, Message: "foo"}, {Number: 1, Message: "bar"}}
	if diff := cmp.Diff(a, b, optNoTS); diff != "" {
		t.Fail()
		t.Log(diff)
	}

	w.Close()

	a = c.uploaded
	if diff := cmp.Diff(a, b, optNoTS); diff != "" {
		t.Fail()
		t.Log(diff)
	}
}

// TestLineWriterLimit tests if the history contains only last uploaded content after the limit has been breached.
func TestLineWriterLimit(t *testing.T) {
	c := newMockClient()

	w := New(c, 0)
	if got, want := w.GetLimit(), defaultLimit; got != want {
		t.Errorf("Expect default buffer limit %d, got %d", want, got)
	}

	w.SetLimit(6)

	if got, want := w.GetLimit(), 6; got != want {
		t.Errorf("Expect custom buffer limit %d, got %d", want, got)
	}

	w.Write([]byte("foo"))
	w.Write([]byte("bar"))
	w.Write([]byte("baz")) // this write overflows the buffer, so "foo" is removed from the history

	if got, want := w.GetSize(), 6; got != want {
		t.Errorf("Expect buffer size %d, got %d", want, got)
	}

	a := w.lineList.peekHistory()
	b := []*drone.Line{{Number: 1, Message: "bar"}, {Number: 2, Message: "baz"}}
	if diff := cmp.Diff(a, b, optNoTS); diff != "" {
		t.Fail()
		t.Log(diff)
	}

	w.Write([]byte("boss")) // "boss" and "baz" are 7 bytes, so "bar" and "baz" are removed

	a = w.lineList.peekHistory()
	b = []*drone.Line{{Number: 3, Message: "boss"}}
	if diff := cmp.Diff(a, b, optNoTS); diff != "" {
		t.Fail()
		t.Log(diff)
	}

	w.Write([]byte("xy")) // this "xy" fits in the buffer so nothing should be removed now

	a = w.lineList.peekHistory()
	b = []*drone.Line{{Number: 3, Message: "boss"}, {Number: 4, Message: "xy"}}
	if diff := cmp.Diff(a, b, optNoTS); diff != "" {
		t.Fail()
		t.Log(diff)
	}

	w.Close()
}

// TestLineWriterLimitStopStreaming tests if streaming has been stopped after the buffer overflow.
func TestLineWriterLimitStopStreaming(t *testing.T) {
	c := newMockClient()
	w := New(c, 0)
	w.SetLimit(8)
	w.SetInterval(time.Nanosecond)

	w.Write([]byte("foo"))
	if uploaded := c.waitUpload(); !uploaded || len(c.lines) != 1 {
		t.Errorf("Expected %d lines streamed, got %d", 1, len(c.lines))
	}

	w.Write([]byte("bar"))
	if uploaded := c.waitUpload(); !uploaded || len(c.lines) != 2 {
		t.Errorf("Expected %d lines streamed, got %d", 2, len(c.lines))
	}

	w.Write([]byte("baz")) // overflow! streaming should be aborted
	if uploaded := c.waitUpload(); uploaded || len(c.lines) != 2 {
		t.Errorf("Expected streaming has been stopped. Streamed %d lines, expected %d", len(c.lines), 2)
	}

	w.Close()

	if len(c.lines) != 2 {
		t.Errorf("Closing should not trigged output streaming. Streamed %d lines, expected %d", len(c.lines), 2)
	}
}

// TestLineWriterOverLimit tests weird situation when data is written in chunks that exceed the limit.
func TestLineWriterOverLimit(t *testing.T) {
	c := newMockClient()

	w := New(c, 0)
	w.SetLimit(4)

	w.Write([]byte("foobar")) // over the limit, nothing should be written

	if got, want := w.GetSize(), 0; got != want {
		t.Errorf("Expect buffer size %d, got %d", want, got)
	}

	w.Close()

	if len(c.uploaded) != 0 {
		t.Error("there should be no uploaded lines")
	}
}

func BenchmarkWriter_Write(b *testing.B) {
	b.ReportAllocs()
	c := &dummyClient{}
	w := New(c, 0)
	p := []byte("Lorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt\nut labore et dolore magna aliqua.\n")
	for i := 0; i < b.N; i++ {
		w.Write(p)
	}
	w.Close()
}

type mockClient struct {
	client.Client
	uploadDone chan struct{}
	lines      []*drone.Line
	uploaded   []*drone.Line
}

func newMockClient() *mockClient {
	return &mockClient{uploadDone: make(chan struct{}, 1)}
}

// waitUpload waits a while for streaming to complete. Writer's interval should be set to very low value before this call.
func (m *mockClient) waitUpload() bool {
	select {
	case <-m.uploadDone:
		return true
	case <-time.After(10 * time.Millisecond):
		return false
	}
}

func (m *mockClient) Batch(ctx context.Context, id int64, lines []*drone.Line) error {
	m.lines = append(m.lines, lines...)
	select {
	case m.uploadDone <- struct{}{}:
	default:
	}
	return nil
}

func (m *mockClient) Upload(ctx context.Context, id int64, lines []*drone.Line) error {
	m.uploaded = lines
	return nil
}

type dummyClient struct{ client.Client }

func (m *dummyClient) Batch(context.Context, int64, []*drone.Line) error  { return nil }
func (m *dummyClient) Upload(context.Context, int64, []*drone.Line) error { return nil }
