// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package console

import "testing"

func TestSequence(t *testing.T) {
	c := new(sequence)
	if got, want := c.curr(), 0; got != want {
		t.Errorf("Want curr sequence value %d, got %d", want, got)
	}
	if got, want := c.next(), 1; got != want {
		t.Errorf("Want next sequence value %d, got %d", want, got)
	}
	if got, want := c.curr(), 1; got != want {
		t.Errorf("Want curr sequence value %d, got %d", want, got)
	}
}
