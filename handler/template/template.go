// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package template

import (
	"regexp"
	"strings"
	"time"
)

//go:generate togo tmpl -func funcMap -format html

// regular expression to extract the pull request number
// from the git ref (e.g. refs/pulls/{d}/head)
var re = regexp.MustCompile("\\d+")

// mirros the func map in template.go
var funcMap = map[string]interface{}{
	"timestamp": func(v int64) string {
		return time.Unix(v, 0).UTC().Format("2006-01-02T15:04:05Z")
	},
	"pr": func(s string) string {
		return re.FindString(s)
	},
	"sha": func(s string) string {
		if len(s) > 8 {
			s = s[:8]
		}
		return s
	},
	"tag": func(s string) string {
		return strings.TrimPrefix(s, "refs/tags/")
	},
	"done": func(s string) bool {
		return s != "pending" && s != "running"
	},
}
