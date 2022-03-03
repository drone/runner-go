// Copyright 2022 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package livelog provides a Writer that collects pipeline
// output and streams to the central server.

package livelog

import (
	"sync"
	"time"

	"github.com/drone/drone-go/drone"
)

type node struct {
	drone.Line
	next *node
}

type list struct {
	sync.Mutex

	lineCnt int
	lineNow time.Time

	size  int
	limit int

	last *node

	history      *node
	historyCount int

	pending      *node
	pendingCount int
}

func makeList(limit int) *list {
	return &list{
		lineCnt: 0,
		lineNow: time.Now(),
		limit:   limit,
	}
}

func (l *list) SetLimit(limit int) {
	l.Lock()
	l.limit = limit
	l.Unlock()
}

func (l *list) GetLimit() int {
	l.Lock()
	limit := l.limit
	l.Unlock()
	return limit
}

func (l *list) GetSize() int {
	l.Lock()
	size := l.size
	l.Unlock()
	return size
}

func (l *list) Push(p []byte) (overflow bool) {
	l.Lock()
	for _, part := range split(p) {
		line := drone.Line{
			Number:    l.lineCnt,
			Message:   part,
			Timestamp: int64(time.Since(l.lineNow).Seconds()),
		}

		l.lineNow = time.Now()
		l.lineCnt++

		overflow = overflow || l.push(line)
	}
	l.Unlock()

	return overflow
}

func (l *list) push(line drone.Line) bool {
	n := &node{
		Line: line,
		next: nil,
	}

	// put the element to list

	l.size += len(line.Message)

	if l.last != nil {
		l.last.next = n
	}
	l.last = n

	if l.history == nil {
		l.history = n
	}
	l.historyCount++

	if l.pending == nil {
		l.pending = n
	}
	l.pendingCount++

	// overflow check

	var overflow bool

	for l.size > l.limit && l.history != nil {
		drop := l.history
		next := drop.next

		if l.pending == drop {
			l.pending = next
			l.pendingCount--
		}

		l.history = next
		l.historyCount--

		if l.history == nil {
			l.last = nil
		}

		l.size -= len(drop.Line.Message)

		overflow = true
	}

	return overflow
}

func (l *list) peekPending() (lines []*drone.Line) {
	l.Lock()
	lines = toSlice(l.pendingCount, l.pending)
	l.Unlock()
	return
}

// Pending returns lines added since the previous call to this method.
func (l *list) Pending() (lines []*drone.Line) {
	l.Lock()
	lines = toSlice(l.pendingCount, l.pending)
	l.pending = nil
	l.pendingCount = 0
	l.Unlock()
	return
}

func (l *list) peekHistory() (lines []*drone.Line) {
	l.Lock()
	lines = toSlice(l.historyCount, l.history)
	l.Unlock()
	return lines
}

// History returns full history stored in the buffer and clears the buffer.
func (l *list) History() (lines []*drone.Line) {
	l.Lock()
	lines = toSlice(l.historyCount, l.history)
	l.history = nil
	l.historyCount = 0
	l.pending = nil
	l.pendingCount = 0
	l.last = nil
	l.size = 0
	l.Unlock()
	return lines
}

func toSlice(count int, head *node) []*drone.Line {
	if count == 0 {
		return nil
	}

	lines := make([]*drone.Line, count)
	for i, n := 0, head; n != nil; i, n = i+1, n.next {
		lines[i] = &n.Line
	}

	return lines
}
