// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package manifest

import (
	"encoding/json"
	"io/ioutil"

	"github.com/google/go-cmp/cmp"
)

func diff(file string) (string, error) {
	a, err := ParseFile(file)
	if err != nil {
		return "", err
	}
	d, err := ioutil.ReadFile(file + ".golden")
	if err != nil {
		return "", err
	}
	b := new(Manifest)
	err = json.Unmarshal(d, b)
	if err != nil {
		return "", err
	}
	return cmp.Diff(a, b), nil
}
