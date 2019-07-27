// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package environ

import (
	"strings"
)

// Proxy returns the http_proxy variables.
func Proxy() map[string]string {
	environ := map[string]string{}
	if value := envAnyCase("no_proxy"); value != "" {
		environ["no_proxy"] = value
		environ["NO_PROXY"] = value
	}
	if value := envAnyCase("http_proxy"); value != "" {
		environ["http_proxy"] = value
		environ["HTTP_PROXY"] = value
	}
	if value := envAnyCase("https_proxy"); value != "" {
		environ["https_proxy"] = value
		environ["HTTPS_PROXY"] = value
	}
	return environ
}

// helper function returns the environment variable value
// using a case-insenstive environment name.
func envAnyCase(name string) (value string) {
	name = strings.ToUpper(name)
	if value := getenv(name); value != "" {
		return value
	}
	name = strings.ToLower(name)
	if value := getenv(name); value != "" {
		return value
	}
	return
}
