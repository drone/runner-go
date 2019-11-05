// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package auths

import (
	"bytes"
	"encoding/json"

	"github.com/drone/drone-go/drone"
)

// Encode encodes the registry credentials to using the
// docker config json format and returns the resulting
// data in string format.
func Encode(registry ...*drone.Registry) string {
	c := new(config)
	c.Auths = map[string]auth{}
	for _, r := range registry {
		c.Auths[r.Address] = auth{
			Auth: encode(r.Username, r.Password),
		}
	}
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.Encode(c)
	return buf.String()
}
