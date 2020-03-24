// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package auths

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/drone/drone-go/drone"
)

type (
	// config represents the Docker client configuration,
	// typically located at ~/.docker/config.json
	config struct {
		Auths map[string]auth `json:"auths"`
	}

	// auth stores the registry authentication string.
	auth struct {
		Auth     string `json:"auth"`
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}
)

// Parse parses the registry credential from the reader.
func Parse(r io.Reader) ([]*drone.Registry, error) {
	c := new(config)
	err := json.NewDecoder(r).Decode(c)
	if err != nil {
		return nil, err
	}
	var auths []*drone.Registry
	for k, v := range c.Auths {
		username, password := v.Username, v.Password
		if v.Auth != "" {
			username, password = decode(v.Auth)
		}
		auths = append(auths, &drone.Registry{
			Address:  hostname(k),
			Username: username,
			Password: password,
		})
	}
	return auths, nil
}

// ParseFile parses the registry credential file.
func ParseFile(filepath string) ([]*drone.Registry, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

// ParseString parses the registry credential file.
func ParseString(s string) ([]*drone.Registry, error) {
	return Parse(strings.NewReader(s))
}

// ParseBytes parses the registry credential file.
func ParseBytes(b []byte) ([]*drone.Registry, error) {
	return Parse(bytes.NewReader(b))
}

// Header returns the json marshaled, base64 encoded
// credential string that can be passed to the docker
// registry authentication header.
func Header(username, password string) string {
	v := struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}{
		Username: username,
		Password: password,
	}
	buf, _ := json.Marshal(&v)
	return base64.URLEncoding.EncodeToString(buf)
}

// encode returns the encoded credentials.
func encode(username, password string) string {
	return base64.StdEncoding.EncodeToString(
		[]byte(username + ":" + password),
	)
}

// decode returns the decoded credentials.
func decode(s string) (username, password string) {
	d, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return
	}
	parts := strings.SplitN(string(d), ":", 2)
	if len(parts) > 0 {
		username = parts[0]
	}
	if len(parts) > 1 {
		password = parts[1]
	}
	return
}

// hostname returns the trimmed hostname from the
// registry url.
func hostname(s string) string {
	uri, _ := url.Parse(s)
	if uri.Host != "" {
		s = uri.Host
	}
	return s
}
