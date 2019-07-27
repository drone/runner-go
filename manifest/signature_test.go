// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package manifest

import (
	"testing"

	"github.com/buildkite/yaml"
	"github.com/google/go-cmp/cmp"
)

var mockSignatureYaml = []byte(`
---
kind: signature
hmac: N2NmYjA3ODQwNTY1ODFlY2E5MGJmOWI1NDk0NDFhMTEK
`)

var mockSignature = &Signature{
	Kind: "signature",
	Hmac: "N2NmYjA3ODQwNTY1ODFlY2E5MGJmOWI1NDk0NDFhMTEK",
}

func TestSignatureUnmarshal(t *testing.T) {
	a := new(Signature)
	b := mockSignature
	yaml.Unmarshal(mockSignatureYaml, a)
	if diff := cmp.Diff(a, b); diff != "" {
		t.Error("Failed to parse signature")
		t.Log(diff)
	}
}

func TestSignatureValidate(t *testing.T) {
	sig := Signature{Hmac: "1234"}
	if err := sig.Validate(); err != nil {
		t.Error(err)
		return
	}

	sig.Hmac = ""
	if err := sig.Validate(); err == nil {
		t.Errorf("Expect invalid signature error")
	}
}
