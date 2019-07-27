// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package manifest

import (
	"errors"

	"github.com/buildkite/yaml"
)

var _ Resource = (*Signature)(nil)

type (
	// Signature is a resource that provides an hmac
	// signature of combined resources. This signature
	// can be used to validate authenticity and prevent
	// tampering.
	Signature struct {
		Version string `json:"version,omitempty"`
		Kind    string `json:"kind"`
		Type    string `json:"type"`
		Name    string `json:"name"`
		Hmac    string `json:"hmac"`
	}
)

func init() {
	Register(signatureFunc)
}

func signatureFunc(r *RawResource) (Resource, bool, error) {
	if r.Kind != KindSignature {
		return nil, false, nil
	}
	out := new(Signature)
	err := yaml.Unmarshal(r.Data, out)
	return out, true, err
}

// GetVersion returns the resource version.
func (s *Signature) GetVersion() string { return s.Version }

// GetKind returns the resource kind.
func (s *Signature) GetKind() string { return s.Kind }

// GetType returns the resource type.
func (s *Signature) GetType() string { return s.Type }

// GetName returns the resource name.
func (s *Signature) GetName() string { return s.Name }

// Validate returns an error if the signature is invalid.
func (s Signature) Validate() error {
	if s.Hmac == "" {
		return errors.New("yaml: invalid signature. missing hash")
	}
	return nil
}
