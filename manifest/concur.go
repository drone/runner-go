// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package manifest

// Concurrency limits pipeline concurrency.
type Concurrency struct {
	Limit int `json:"limit,omitempty"`
}
