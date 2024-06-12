// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package auths

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/drone/drone-go/drone"
	"golang.org/x/oauth2"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/sts/v1"
)

const (
	gcpAudienceFormat = "//iam.googleapis.com/projects/%s/locations/global/workloadIdentityPools/%s/providers/%s"
	gcpScopeURL       = "https://www.googleapis.com/auth/cloud-platform"
)

type (
	// config represents the Docker client configuration,
	// typically located at ~/.docker/config.json
	config struct {
		Auths map[string]auth `json:"auths"`
	}

	// auth stores the registry authentication string.
	auth struct {
		Auth      string `json:"auth"`
		Username  string `json:"username,omitempty"`
		Password  string `json:"password,omitempty"`
		OIDCToken string `json:"oidc_token,omitempty"`
	}
)

// Parse parses the registry credential from the reader.
func Parse(r io.Reader) ([]*drone.Registry, error) {
	c := new(config)
	err := json.NewDecoder(r).Decode(c)
	if err != nil {
		return nil, err
	}
	var auths []*drone.Registry
	for k, v := range c.Auths {
		username, password := v.Username, v.Password
		if v.Auth != "" {
			username, password = decode(v.Auth)
		}
		auths = append(auths, &drone.Registry{
			Address:  hostname(k),
			Username: username,
			Password: password,
		})
		if v.OIDCToken != "" {
			auths = append(auths, &drone.Registry{
				Address: k,
				Token:   v.OIDCToken,
			})
		}
	}
	return auths, nil
}

// ParseFile parses the registry credential file.
func ParseFile(filepath string) ([]*drone.Registry, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

// ParseString parses the registry credential file.
func ParseString(s string) ([]*drone.Registry, error) {
	return Parse(strings.NewReader(s))
}

// ParseBytes parses the registry credential file.
func ParseBytes(b []byte) ([]*drone.Registry, error) {
	return Parse(bytes.NewReader(b))
}

// Header returns the json marshaled, base64 encoded
// credential string that can be passed to the docker
// registry authentication header.
func Header(username, password string, token string) string {
	v := struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
		Token    string `json:"identitytoken,omitempty"`
	}{
		Username: username,
		Password: password,
		Token:    token,
	}
	buf, _ := json.Marshal(&v)
	return base64.URLEncoding.EncodeToString(buf)
}

// encode returns the encoded credentials.
func encode(username, password string) string {
	return base64.StdEncoding.EncodeToString(
		[]byte(username + ":" + password),
	)
}

// decode returns the decoded credentials.
func decode(s string) (username, password string) {
	d, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return
	}
	parts := strings.SplitN(string(d), ":", 2)
	if len(parts) > 0 {
		username = parts[0]
	}
	if len(parts) > 1 {
		password = parts[1]
	}
	return
}

// hostname returns the trimmed hostname from the
// registry url.
func hostname(s string) string {
	uri, _ := url.Parse(s)
	if uri != nil && uri.Host != "" {
		s = uri.Host
	}
	return s
}

func GetGcpFederalToken(idToken, projectNumber, poolId, providerId string) (string, error) {
	ctx := context.Background()
	stsService, err := sts.NewService(ctx, option.WithoutAuthentication())
	if err != nil {
		return "", err
	}

	audience := fmt.Sprintf(gcpAudienceFormat, projectNumber, poolId, providerId)

	tokenRequest := &sts.GoogleIdentityStsV1ExchangeTokenRequest{
		GrantType:          "urn:ietf:params:oauth:grant-type:token-exchange",
		SubjectToken:       idToken,
		Audience:           audience,
		Scope:              gcpScopeURL,
		RequestedTokenType: "urn:ietf:params:oauth:token-type:access_token",
		SubjectTokenType:   "urn:ietf:params:oauth:token-type:id_token",
	}

	tokenResponse, err := stsService.V1.Token(tokenRequest).Do()
	if err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func GetGoogleCloudAccessToken(federatedToken string, serviceAccountEmail string) (string, error) {
	ctx := context.Background()
	token := &oauth2.Token{AccessToken: federatedToken}
	service, err := iamcredentials.NewService(ctx, option.WithTokenSource(oauth2.StaticTokenSource(token)))
	if err != nil {
		return "", err
	}

	name := "projects/-/serviceAccounts/" + serviceAccountEmail
	// rb (request body) specifies parameters for generating an access token.
	rb := &iamcredentials.GenerateAccessTokenRequest{
		Scope: []string{gcpScopeURL},
	}
	// Generate an access token for the service account using the specified parameters
	resp, err := service.Projects.ServiceAccounts.GenerateAccessToken(name, rb).Do()
	if err != nil {
		return "", err
	}

	return resp.AccessToken, nil
}
