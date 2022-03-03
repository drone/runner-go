// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package livelog provides a Writer that collects pipeline
// output and streams to the central server.
package livelog

import (
	"context"
	"strings"
	"time"

	"github.com/drone/runner-go/client"
)

// defaultLimit is the default maximum log size in bytes.
const defaultLimit = 5242880 // 5MB

// Writer is an io.WriteCloser that sends logs to the server.
type Writer struct {
	client client.Client

	id int64

	interval time.Duration
	lineList *list

	stopStreamFn func()
	doneStream   <-chan struct{}
	ready        chan struct{}
}

// New returns a new Writer.
func New(client client.Client, id int64) *Writer {
	streamCtx, stopStream := context.WithCancel(context.Background())

	b := &Writer{
		client:       client,
		id:           id,
		interval:     time.Second,
		lineList:     makeList(defaultLimit),
		stopStreamFn: stopStream,
		doneStream:   streamCtx.Done(),
		ready:        make(chan struct{}, 1),
	}

	// a call to stopStream() stops this goroutine.
	// this happens when the Close method is called or after overflow of output data (>limit).
	go b.start()

	return b
}

// SetLimit sets the Writer limit.
func (b *Writer) SetLimit(limit int) {
	b.lineList.SetLimit(limit)
}

// GetLimit returns the Writer limit.
func (b *Writer) GetLimit() int {
	return b.lineList.GetLimit()
}

// GetSize returns amount of output data the Writer currently holds.
func (b *Writer) GetSize() int {
	return b.lineList.GetSize()
}

// SetInterval sets the Writer flusher interval.
func (b *Writer) SetInterval(interval time.Duration) {
	b.interval = interval
}

// Write uploads the live log stream to the server.
func (b *Writer) Write(p []byte) (n int, err error) {
	if isOverLimit := b.lineList.Push(p); isOverLimit {
		b.stopStreamFn()
	}

	select {
	case b.ready <- struct{}{}:
	default:
	}

	return len(p), nil
}

// Close closes the writer and uploads the full contents to the server.
func (b *Writer) Close() error {
	select {
	case <-b.doneStream:
	default:
		b.stopStreamFn()
		_ = b.flush() // send all pending lines
	}

	return b.upload() // upload full log history
}

// upload uploads the full log history to the server.
func (b *Writer) upload() error {
	return b.client.Upload(context.Background(), b.id, b.lineList.History())
}

// flush batch uploads all buffered logs to the server.
func (b *Writer) flush() error {
	lines := b.lineList.Pending()
	if len(lines) == 0 {
		return nil
	}

	return b.client.Batch(context.Background(), b.id, lines)
}

func (b *Writer) start() {
	for {
		select {
		case <-b.doneStream:
			return
		case <-b.ready:
			select {
			case <-b.doneStream:
				return
			case <-time.After(b.interval):
				// we intentionally ignore errors. log streams
				// are ephemeral and are considered low priority
				// because they are not required for drone to
				// operator, and the impact of failure is minimal
				_ = b.flush()
			}
		}
	}
}

func split(p []byte) []string {
	s := string(p)
	v := []string{s}
	// kubernetes buffers the output and may combine
	// multiple lines into a single block of output.
	// Split into multiple lines.
	//
	// note that docker output always inclines a line
	// feed marker. This needs to be accounted for when
	// splitting the output into multiple lines.
	if strings.Contains(strings.TrimSuffix(s, "\n"), "\n") {
		v = strings.SplitAfter(s, "\n")
	}
	return v
}
