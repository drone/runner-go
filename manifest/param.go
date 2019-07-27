// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package manifest

type (
	// Parameter represents an configuration parameter that
	// can be defined as a literal or as a reference
	// to a secret.
	Parameter struct {
		Value  interface{} `json:"value,omitempty"`
		Secret string      `json:"from_secret,omitempty" yaml:"from_secret"`
	}

	// parameter is a tempoary type used to unmarshal
	// parameters with references to secrets.
	parameter struct {
		Secret string `yaml:"from_secret"`
	}
)

// UnmarshalYAML implements yaml unmarshalling.
func (p *Parameter) UnmarshalYAML(unmarshal func(interface{}) error) error {
	d := new(parameter)
	err := unmarshal(d)
	if err == nil && d.Secret != "" {
		p.Secret = d.Secret
		return nil
	}
	var i interface{}
	err = unmarshal(&i)
	p.Value = i
	return err
}

// MarshalYAML implements yaml marshalling.
func (p *Parameter) MarshalYAML() (interface{}, error) {
	if p.Secret != "" {
		m := map[string]interface{}{}
		m["from_secret"] = p.Secret
		return m, nil
	}
	if p.Value != "" {
		return p.Value, nil
	}
	return nil, nil
}
