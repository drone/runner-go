// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package clone provides utilities for cloning commits.
package clone

import "fmt"

//
// IMPORTANT: DO NOT MODIFY THIS FILE
//
// this file must not be changed unless the changes have been
// discussed and approved by the project maintainers in the
// GitHub issue tracker.
//

// Args provide arguments to clone a repository.
type Args struct {
	Branch string
	Commit string
	Ref    string
	Remote string
	Depth  int
	Tags   bool
	NoFF   bool
}

// Commands returns posix-compliant commands to clone a
// repository and checkout a commit.
func Commands(args Args) []string {
	switch {
	case isTag(args.Ref):
		return tag(args)
	case isPullRequest(args.Ref):
		return pull(args)
	default:
		return branch(args)
	}
}

// branch returns posix-compliant commands to clone a repository
// and checkout the named branch.
func branch(args Args) []string {
	return []string{
		"git init",
		fmt.Sprintf("git remote add origin %s", args.Remote),
		fmt.Sprintf("git fetch %s origin +refs/heads/%s:", fetchFlags(args), args.Branch),
		fmt.Sprintf("git checkout %s -b %s", args.Commit, args.Branch),
	}
}

// tag returns posix-compliant commands to clone a repository
// and checkout the tag by reference path.
func tag(args Args) []string {
	return []string{
		"git init",
		fmt.Sprintf("git remote add origin %s", args.Remote),
		fmt.Sprintf("git fetch %s origin +%s:", fetchFlags(args), args.Ref),
		"git checkout -qf FETCH_HEAD",
	}
}

// pull returns posix-compliant commands to clone a repository
// and checkout the pull request by reference path.
func pull(args Args) []string {
	return []string{
		"git init",
		fmt.Sprintf("git remote add origin %s", args.Remote),
		fmt.Sprintf("git fetch %s origin +refs/heads/%s:", fetchFlags(args), args.Branch),
		fmt.Sprintf("git checkout %s", args.Branch),
		fmt.Sprintf("git fetch origin %s:", args.Ref),
		fmt.Sprintf("git merge %s %s", mergeFlags(args), args.Commit),
	}
}
