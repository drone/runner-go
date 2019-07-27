// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package manifest

import (
	"errors"

	"github.com/buildkite/yaml"
)

var _ Resource = (*Secret)(nil)

type (
	// Secret is a resource that provides encrypted data
	// and pointers to external data (i.e. from vault).
	Secret struct {
		Version string    `json:"version,omitempty"`
		Kind    string    `json:"kind,omitempty"`
		Type    string    `json:"type,omitempty"`
		Name    string    `json:"name,omitempty"`
		Data    string    `json:"data,omitempty"`
		Get     SecretGet `json:"get,omitempty"`
	}

	// SecretGet defines a request to get a secret from
	// an external sevice at the specified path, or with the
	// specified name.
	SecretGet struct {
		Path string `json:"path,omitempty"`
		Name string `json:"name,omitempty"`
		Key  string `json:"key,omitempty"`
	}
)

func init() {
	Register(secretFunc)
}

func secretFunc(r *RawResource) (Resource, bool, error) {
	if r.Kind != KindSecret {
		return nil, false, nil
	}
	out := new(Secret)
	err := yaml.Unmarshal(r.Data, out)
	return out, true, err
}

// GetVersion returns the resource version.
func (s *Secret) GetVersion() string { return s.Version }

// GetKind returns the resource kind.
func (s *Secret) GetKind() string { return s.Kind }

// GetType returns the resource type.
func (s *Secret) GetType() string { return s.Type }

// GetName returns the resource name.
func (s *Secret) GetName() string { return s.Name }

// Validate returns an error if the secret is invalid.
func (s *Secret) Validate() error {
	if len(s.Data) == 0 &&
		len(s.Get.Path) == 0 &&
		len(s.Get.Name) == 0 {
		return errors.New("yaml: invalid secret resource")
	}
	return nil
}
