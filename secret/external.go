// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package secret

import (
	"context"
	"time"

	"github.com/drone/runner-go/logger"
	"github.com/drone/runner-go/manifest"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/secret"
)

// External returns a new external secret provider. The
// external secret provider makes an external API call to find
// and return a named secret.
func External(endpoint, token string, insecure bool) Provider {
	provider := &external{}
	if endpoint != "" {
		provider.client = secret.Client(endpoint, token, insecure)
	}
	return provider
}

type external struct {
	client secret.Plugin
}

func (p *external) Find(ctx context.Context, in *Request) (*drone.Secret, error) {
	if p.client == nil {
		return nil, nil
	}

	logger := logger.FromContext(ctx).
		WithField("name", in.Name).
		WithField("kind", "secret")

	// lookup the named secret in the manifest. If the
	// secret does not exist, return a nil variable,
	// allowing the next secret controller in the chain
	// to be invoked.
	path, name, ok := getExternal(in.Conf, in.Name)
	if !ok {
		logger.Trace("secret: external: no matching secret")
		return nil, nil
	}

	// include a timeout to prevent an API call from
	// hanging the build process indefinitely. The
	// external service must return a request within
	// one minute.
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	req := &secret.Request{
		Name:  name,
		Path:  path,
		Repo:  *in.Repo,
		Build: *in.Build,
	}
	res, err := p.client.Find(ctx, req)
	if err != nil {
		logger.WithError(err).Debug("secret: external: cannot get secret")
		return nil, err
	}

	// if no error is returned and the secret is empty,
	// this indicates the client returned No Content,
	// and we should exit with no secret, but no error.
	if res.Data == "" {
		logger.Trace("secret: external: secret is empty")
		return nil, nil
	}

	logger.Trace("secret: external: found matching secret")

	return &drone.Secret{
		Name:        in.Name,
		Data:        res.Data,
		PullRequest: res.Pull,
	}, nil
}

func getExternal(spec *manifest.Manifest, match string) (path, name string, ok bool) {
	for _, resource := range spec.Resources {
		secret, ok := resource.(*manifest.Secret)
		if !ok {
			continue
		}
		if secret.Name != match {
			continue
		}
		if secret.Get.Name == "" && secret.Get.Path == "" {
			continue
		}
		return secret.Get.Path, secret.Get.Name, true
	}
	return
}
