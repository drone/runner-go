// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package manifest

import (
	"testing"

	"github.com/buildkite/yaml"
)

func TestEnv(t *testing.T) {
	tests := []struct {
		yaml  string
		value string
		from  string
	}{
		{
			yaml:  "bar",
			value: "bar",
		},
		{
			yaml: "from_secret: username",
			from: "username",
		},
	}
	for _, test := range tests {
		in := []byte(test.yaml)
		out := new(Variable)
		err := yaml.Unmarshal(in, out)
		if err != nil {
			t.Error(err)
			return
		}
		if got, want := out.Value, test.value; got != want {
			t.Errorf("Want variable value %q, got %q", want, got)
		}
		if got, want := out.Secret, test.from; got != want {
			t.Errorf("Want variable from_secret %q, got %q", want, got)
		}
	}
}
