// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package manifest

// Workspace configures the project path on disk.
type Workspace struct {
	Path string `json:"path,omitempty"`
}
