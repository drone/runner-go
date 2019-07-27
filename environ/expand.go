// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package environ

import "os"

// function used to expand environment variables.
var getenv = os.Getenv

// Expand is a helper function to expand the PATH parameter in
// the pipeline environment.
func Expand(env map[string]string) map[string]string {
	c := map[string]string{}
	for k, v := range env {
		c[k] = v
	}
	if path := c["PATH"]; path != "" {
		c["PATH"] = os.Expand(path, getenv)
	}
	return c
}
