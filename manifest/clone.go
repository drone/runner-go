// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package manifest

// Clone configures the git clone.
type Clone struct {
	Disable    bool `json:"disable,omitempty"`
	Depth      int  `json:"depth,omitempty"`
	SkipVerify bool `json:"skip_verify,omitempty" yaml:"skip_verify"`
	Trace      bool `json:"trace,omitempty"`
}
