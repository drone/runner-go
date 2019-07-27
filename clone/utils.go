// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package clone

import (
	"fmt"
	"strings"
)

// fetchFlags is a helper function that returns a string of
// optional git-fetch command line flags.
func fetchFlags(args Args) string {
	var flags []string
	if depth := args.Depth; depth > 0 {
		flag := fmt.Sprintf("--depth=%d", depth)
		flags = append(flags, flag)
	}
	if args.Tags {
		flags = append(flags, "--tags")
	}
	return strings.Join(flags, " ")
}

// mergeFlags is a helper function that returns a string of
// optional git-merge command line flags.
func mergeFlags(args Args) string {
	var flags []string
	if args.NoFF {
		flags = append(flags, "--no-ff")
	}
	return strings.Join(flags, " ")
}

// isTag returns true if the reference path points to
// a tag object.
func isTag(ref string) bool {
	return strings.HasPrefix(ref, "refs/tags/")
}

// isPullRequest returns true if the reference path points to
// a pull request object.
func isPullRequest(ref string) bool {
	return strings.HasPrefix(ref, "refs/pull/") ||
		strings.HasPrefix(ref, "refs/pull-requests/") ||
		strings.HasPrefix(ref, "refs/merge-requests/")
}
