// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package manifest

type (
	// Variable represents an environment variable that
	// can be defined as a string literal or as a reference
	// to a secret.
	Variable struct {
		Value  string `json:"value,omitempty"`
		Secret string `json:"from_secret,omitempty" yaml:"from_secret"`
	}

	// variable is a temporary type used to unmarshal
	// variables with references to secrets.
	variable struct {
		Value  string
		Secret string `yaml:"from_secret"`
	}
)

// UnmarshalYAML implements yaml unmarshalling.
func (v *Variable) UnmarshalYAML(unmarshal func(interface{}) error) error {
	d := new(variable)
	err := unmarshal(&d.Value)
	if err != nil {
		err = unmarshal(d)
	}
	v.Value = d.Value
	v.Secret = d.Secret
	return err
}

// MarshalYAML implements yaml marshalling.
func (v *Variable) MarshalYAML() (interface{}, error) {
	if v.Secret != "" {
		m := map[string]interface{}{}
		m["from_secret"] = v.Secret
		return m, nil
	}
	if v.Value != "" {
		return v.Value, nil
	}
	return nil, nil
}
