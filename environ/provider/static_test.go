// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package provider

import (
	"reflect"
	"testing"
)

func TestStatic(t *testing.T) {
	a := map[string]string{"a": "b"}
	p := Static(a)
	b, err := p.List(noContext, nil)
	if err != nil {
		t.Error(err)
		return
	}
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Unexpected environment variable output")
	}
}
