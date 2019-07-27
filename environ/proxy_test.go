// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package environ

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProxy(t *testing.T) {
	defer func() {
		getenv = os.Getenv
	}()

	getenv = func(s string) string {
		switch s {
		case "no_proxy":
			return "http://dummy.no.proxy"
		case "http_proxy":
			return "http://dummy.http.proxy"
		case "https_proxy":
			return "http://dummy.https.proxy"
		default:
			return ""
		}
	}

	a := map[string]string{
		"no_proxy":    "http://dummy.no.proxy",
		"NO_PROXY":    "http://dummy.no.proxy",
		"http_proxy":  "http://dummy.http.proxy",
		"HTTP_PROXY":  "http://dummy.http.proxy",
		"https_proxy": "http://dummy.https.proxy",
		"HTTPS_PROXY": "http://dummy.https.proxy",
	}
	b := Proxy()
	if diff := cmp.Diff(a, b); diff != "" {
		t.Fail()
		t.Log(diff)
	}
}

func Test_envAnyCase(t *testing.T) {
	defer func() {
		getenv = os.Getenv
	}()

	getenv = func(s string) string {
		switch s {
		case "foo":
			return "bar"
		default:
			return ""
		}
	}

	if envAnyCase("FOO") != "bar" {
		t.Errorf("Expect environment variable sourced from lowercase variant")
	}

	getenv = func(s string) string {
		switch s {
		case "FOO":
			return "bar"
		default:
			return ""
		}
	}

	if envAnyCase("foo") != "bar" {
		t.Errorf("Expect environment variable sourced from uppercase variant")
	}

	if envAnyCase("bar") != "" {
		t.Errorf("Expect zero value string when environment variable does not exit")
	}
}
