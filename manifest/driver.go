// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package manifest

// registered drivers.
var drivers []Driver

// Register registers the parsing driver.
func Register(driver Driver) {
	drivers = append(drivers, driver)
}

// Driver defines a parser driver that can be used to parse
// resource-specific Yaml documents.
type Driver func(r *RawResource) (Resource, bool, error)
