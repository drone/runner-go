// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package secret

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"

	"github.com/drone/runner-go/logger"
	"github.com/drone/runner-go/manifest"

	"github.com/drone/drone-go/drone"
)

// Encrypted returns a new encrypted secret provider. The
// encrypted secret provider finds and decrypts secrets stored
// inline in the (yaml) configuration.
func Encrypted() Provider {
	return new(encrypted)
}

type encrypted struct{}

func (p *encrypted) Find(ctx context.Context, in *Request) (*drone.Secret, error) {
	logger := logger.FromContext(ctx).
		WithField("name", in.Name).
		WithField("kind", "secret")

	// lookup the named secret in the manifest. If the
	// secret does not exist, return a nil variable,
	// allowing the next secret controller in the chain
	// to be invoked.
	data, ok := getEncrypted(in.Conf, in.Name)
	if !ok {
		logger.Trace("secret: encrypted: no matching secret")
		return nil, nil
	}

	// if the build event is a pull request and the source
	// repository is a fork, the secret is not exposed to
	// the pipeline, for security reasons.
	if in.Repo.Private == false &&
		in.Build.Event == drone.EventPullRequest &&
		in.Build.Fork != "" {
		logger.Trace("secret: encrypted: restricted from forks")
		return nil, nil
	}

	decoded, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		logger.WithError(err).Debug("secret: encrypted: cannot decode")
		return nil, err
	}

	decrypted, err := decrypt(decoded, []byte(in.Repo.Secret))
	if err != nil {
		logger.WithError(err).Debug("secret: encrypted: cannot decrypt")
		return nil, err
	}

	logger.Trace("secret: encrypted: found matching secret")

	return &drone.Secret{
		Name: in.Name,
		Data: string(decrypted),
	}, nil
}

func getEncrypted(spec *manifest.Manifest, match string) (data string, ok bool) {
	for _, resource := range spec.Resources {
		secret, ok := resource.(*manifest.Secret)
		if !ok {
			continue
		}
		if secret.Name != match {
			continue
		}
		if secret.Data == "" {
			continue
		}
		return secret.Data, true
	}
	return
}

func decrypt(ciphertext []byte, key []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}
