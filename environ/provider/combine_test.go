// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package provider

import (
	"errors"
	"testing"
)

func TestCombine(t *testing.T) {
	a := map[string]string{"a": "b"}
	b := map[string]string{"c": "d"}
	aa := Static(a)
	bb := Static(b)
	p := Combine(aa, bb)
	out, err := p.List(noContext, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if len(out) != 2 {
		t.Errorf("Expect combined variable output")
		return
	}
	if out["a"] != "b" {
		t.Errorf("Missing variable")
	}
	if out["c"] != "d" {
		t.Errorf("Missing variable")
	}
}

func TestCombineError(t *testing.T) {
	e := errors.New("not found")
	m := mockProvider{err: e}
	p := Combine(&m)
	_, err := p.List(noContext, nil)
	if err != e {
		t.Errorf("Expect error")
	}
}
