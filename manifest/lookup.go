// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package manifest

import "errors"

// Lookup returns the named resource from the Manifest.
func Lookup(name string, manifest *Manifest) (Resource, error) {
	for _, resource := range manifest.Resources {
		if isNameMatch(resource.GetName(), name) {
			return resource, nil
		}
	}
	return nil, errors.New("resource not found")
}

// helper function returns true if the name matches.
func isNameMatch(a, b string) bool {
	return a == b ||
		(a == "" && b == "default") ||
		(b == "" && a == "default")
}
