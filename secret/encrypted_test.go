// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package secret

import (
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/manifest"
)

func TestEncrypted(t *testing.T) {
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPullRequest},
		Repo:  &drone.Repo{Private: false, Secret: passcode},
		Conf: &manifest.Manifest{
			Resources: []manifest.Resource{
				&manifest.Secret{
					Name: "docker_password",
					Data: ciphertext,
				},
			},
		},
	}
	service := Encrypted()
	secret, err := service.Find(noContext, args)
	if err != nil {
		t.Errorf("Expect secret found and decrypted, got error. %s", err)
	}
	if secret == nil {
		t.Errorf("Expect secret found and decrypted")
	}
	if got, want := secret.Data, plaintext; got != want {
		t.Errorf("Expect plaintext %s, got %s", want, got)
	}
}

// This test verifies that a nil secret and nil error are
// returned if no corresponding secret resource entry is found
// in the manifest. No error is returned because Not Found is
// not considered an error.
func TestEncrypted_NotFound(t *testing.T) {
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPush},
		Conf: &manifest.Manifest{
			Resources: []manifest.Resource{
				&manifest.Signature{
					Name: "signature",
					Hmac: "<some signature>",
				},
				&manifest.Secret{
					Name: "docker_username",
					Data: ciphertext,
				},
			},
		},
	}
	service := Encrypted()
	secret, err := service.Find(noContext, args)
	if err != nil {
		t.Errorf("Expect nil error when secret not exists")
	}
	if secret != nil {
		t.Errorf("Expect nil secret when secret not exists")
	}
}

// This test verifies that if a secret is decrypted and the
// value is an empty string, a nil secret and nil error are
// returned by the provider.
func TestEncrypted_EmptyData(t *testing.T) {
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPush},
		Conf: &manifest.Manifest{
			Resources: []manifest.Resource{
				&manifest.Secret{
					Name: "docker_password",
					Data: "", // should skip
				},
			},
		},
	}
	service := Encrypted()
	secret, err := service.Find(noContext, args)
	if err != nil {
		t.Errorf("Expect nil error when secret not exists")
	}
	if secret != nil {
		t.Errorf("Expect nil secret when secret not exists")
	}
}

// This test verifies that encrypted secrets are not exposed
// the all the following criteria is met, a) the repository is
// private b) the event is a pull request and c) the pull
// requests comes from a fork.
func TestEncrypted_PullRequest_Public_Fork(t *testing.T) {
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPullRequest, Fork: "octocat/hello-world"},
		Repo:  &drone.Repo{Private: false},
		Conf: &manifest.Manifest{
			Resources: []manifest.Resource{
				&manifest.Secret{
					Name: "docker_password",
					Data: ciphertext,
				},
			},
		},
	}
	service := Encrypted()
	secret, err := service.Find(noContext, args)
	if err != nil {
		t.Errorf("Expect nil error when secret not exists")
	}
	if secret != nil {
		t.Errorf("Expect nil secret when secret not exists")
	}
}

// This test verifies that an error decoding an ecrypted secret
// is returned to the caller.
func TestEncrypted_DecodeError(t *testing.T) {
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPullRequest},
		Repo:  &drone.Repo{Private: false, Secret: passcode},
		Conf: &manifest.Manifest{
			Resources: []manifest.Resource{
				&manifest.Secret{
					Name: "docker_password",
					Data: "<invalid text>",
				},
			},
		},
	}
	service := Encrypted()
	_, err := service.Find(noContext, args)
	if err == nil {
		t.Errorf("Expect base64 decode error")
	}
}

// This test verifies that if a secret key is invalid, resulting
// in a decryption error, the error is returned to the caller.
func TestEncrypted_Decrypt_InvalidKey(t *testing.T) {
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPullRequest},
		Repo:  &drone.Repo{Private: false, Secret: "id0cvCE0F1tgWeR4WvAXwChyRdDi8gHL"},
		Conf: &manifest.Manifest{
			Resources: []manifest.Resource{
				&manifest.Secret{
					Name: "docker_password",
					Data: ciphertext,
				},
			},
		},
	}
	service := Encrypted()
	_, err := service.Find(noContext, args)
	if err == nil {
		t.Errorf("Expect decryption error")
	}
}

// This test verifies that if a secret key is malformed,
// resulting in a decryption error, the error is returned to the
// caller.
func TestEncrypted_Decrypt_MalformedKey(t *testing.T) {
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPullRequest},
		Repo:  &drone.Repo{Private: false, Secret: "<invalid key>"},
		Conf: &manifest.Manifest{
			Resources: []manifest.Resource{
				&manifest.Secret{
					Name: "docker_password",
					Data: ciphertext,
				},
			},
		},
	}
	service := Encrypted()
	_, err := service.Find(noContext, args)
	if err == nil {
		t.Errorf("Expect decryption error")
	}
}

// This test verifies that if an encrypted secret is malformed,
// resulting in a decryption error, the error is returned to the
// caller.
func TestEncrypted_Decrypt_MalformedBlock(t *testing.T) {
	args := &Request{
		Name:  "docker_password",
		Build: &drone.Build{Event: drone.EventPullRequest},
		Repo:  &drone.Repo{Private: false, Secret: passcode},
		Conf: &manifest.Manifest{
			Resources: []manifest.Resource{
				&manifest.Secret{
					Name: "docker_password",
					Data: "<malformed ciphertext>",
				},
			},
		},
	}
	service := Encrypted()
	_, err := service.Find(noContext, args)
	if err == nil {
		t.Errorf("Expect decryption error")
	}
}

var (
	plaintext  = `correct-horse-battery-staple`
	ciphertext = `6OEK5rMVitI+S7//bLJSuwbIGq7/rtlj/V30DhhUa1sYWQd+CsCb/YUUhJ1F6pRBTVlBmiFHxq4=`
	passcode   = `O9AhOJ23FHtiiXejAl381QG4fg9Zv3LK`
)
