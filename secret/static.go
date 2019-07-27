// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package secret

import (
	"context"
	"strings"

	"github.com/drone/runner-go/logger"

	"github.com/drone/drone-go/drone"
)

// Static returns a new static secret provider. The static
// secret provider finds and returns a named secret from the
// static list.
func Static(secrets []*drone.Secret) Provider {
	return &static{secrets}
}

// StaticVars returns a new static secret provider. The static
// secret provider finds and returns a named secret from the
// static key value pairs.
func StaticVars(vars map[string]string) Provider {
	var secrets []*drone.Secret
	for k, v := range vars {
		secrets = append(secrets, &drone.Secret{
			Name: k,
			Data: v,
		})
	}
	return Static(secrets)
}

type static struct {
	secrets []*drone.Secret
}

func (p *static) Find(ctx context.Context, in *Request) (*drone.Secret, error) {
	logger := logger.FromContext(ctx).
		WithField("name", in.Name).
		WithField("kind", "secret")

	for _, secret := range p.secrets {
		if !strings.EqualFold(secret.Name, in.Name) {
			continue
		}
		// The secret can be restricted to non-pull request
		// events. If the secret is restricted, return
		// empty results.
		if secret.PullRequest == false &&
			in.Build.Event == drone.EventPullRequest {
			logger.Trace("secret: database: restricted from pull requests")
			continue
		}

		logger.Trace("secret: database: found matching secret")
		return secret, nil
	}

	logger.Trace("secret: database: no matching secret")
	return nil, nil
}
