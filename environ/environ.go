// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package environ provides utilities for generating environment
// variables for a build pipeline.
package environ

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/drone/drone-go/drone"
)

// regular expression to extract the pull request number
// from the git ref (e.g. refs/pulls/{d}/head)
var re = regexp.MustCompile("\\d+")

// System returns a set of environment variables containing
// system metadata.
func System(system *drone.System) map[string]string {
	return map[string]string{
		"CI":                    "true",
		"DRONE":                 "true",
		"DRONE_SYSTEM_PROTO":    system.Proto,
		"DRONE_SYSTEM_HOST":     system.Host,
		"DRONE_SYSTEM_HOSTNAME": system.Host,
		"DRONE_SYSTEM_VERSION":  fmt.Sprint(system.Version),
	}
}

// Repo returns a set of environment variables containing
// repository metadata.
func Repo(repo *drone.Repo) map[string]string {
	return map[string]string{
		"DRONE_REPO":            repo.Slug,
		"DRONE_REPO_SCM":        repo.SCM,
		"DRONE_REPO_OWNER":      repo.Namespace,
		"DRONE_REPO_NAMESPACE":  repo.Namespace,
		"DRONE_REPO_NAME":       repo.Name,
		"DRONE_REPO_LINK":       repo.Link,
		"DRONE_REPO_BRANCH":     repo.Branch,
		"DRONE_REMOTE_URL":      repo.HTTPURL,
		"DRONE_GIT_HTTP_URL":    repo.HTTPURL,
		"DRONE_GIT_SSH_URL":     repo.SSHURL,
		"DRONE_REPO_VISIBILITY": repo.Visibility,
		"DRONE_REPO_PRIVATE":    fmt.Sprint(repo.Private),

		// these are legacy configuration parameters for backward
		// compatibility with drone 0.8. These are deprecated and
		// should not be relied upon going forward.
		"CI_REPO":         repo.Slug,
		"CI_REPO_NAME":    repo.Slug,
		"CI_REPO_LINK":    repo.Link,
		"CI_REPO_REMOTE":  repo.HTTPURL,
		"CI_REMOTE_URL":   repo.HTTPURL,
		"CI_REPO_PRIVATE": fmt.Sprint(repo.Private),
	}
}

// Stage returns a set of environment variables containing
// stage metadata.
func Stage(stage *drone.Stage) map[string]string {
	env := map[string]string{
		"DRONE_STAGE_KIND":       stage.Kind,
		"DRONE_STAGE_TYPE":       stage.Type,
		"DRONE_STAGE_NAME":       stage.Name,
		"DRONE_STAGE_NUMBER":     fmt.Sprint(stage.Number),
		"DRONE_STAGE_MACHINE":    stage.Machine,
		"DRONE_STAGE_OS":         stage.OS,
		"DRONE_STAGE_ARCH":       stage.Arch,
		"DRONE_STAGE_VARIANT":    stage.Variant,
		"DRONE_STAGE_VERSION":    fmt.Sprint(stage.Version),
		"DRONE_STAGE_STATUS":     "success",
		"DRONE_STAGE_STARTED":    fmt.Sprint(stage.Started),
		"DRONE_STAGE_FINISHED":   fmt.Sprint(stage.Stopped),
		"DRONE_STAGE_DEPENDS_ON": strings.Join(stage.DependsOn, ","),
		"DRONE_CARD_PATH":        "/dev/stdout",
	}
	if isStageFailing(stage) {
		env["DRONE_STAGE_STATUS"] = "failure"
		env["DRONE_FAILED_STEPS"] = strings.Join(failedSteps(stage), ",")
	}
	if stage.Started == 0 {
		env["DRONE_STAGE_STARTED"] = fmt.Sprint(time.Now().Unix())
	}
	if stage.Stopped == 0 {
		env["DRONE_STAGE_FINISHED"] = fmt.Sprint(time.Now().Unix())
	}
	return env
}

// Step returns a set of environment variables containing the
// step metadata.
func Step(step *drone.Step) map[string]string {
	return map[string]string{
		"DRONE_STEP_NAME":   step.Name,
		"DRONE_STEP_NUMBER": fmt.Sprint(step.Number),
	}
}

// StepArgs returns a set of environment variables containing
// the step name and number.
func StepArgs(name string, number int64) map[string]string {
	return map[string]string{
		"DRONE_STEP_NAME":   name,
		"DRONE_STEP_NUMBER": fmt.Sprint(number),
	}
}

// StepName returns a set of environment variables containing
// only the step name.
func StepName(name string) map[string]string {
	return map[string]string{
		"DRONE_STEP_NAME": name,
	}
}

// Build returns a set of environment variables containing
// build metadata.
func Build(build *drone.Build) map[string]string {
	env := map[string]string{
		"DRONE_BRANCH":               build.Target,
		"DRONE_SOURCE_BRANCH":        build.Source,
		"DRONE_TARGET_BRANCH":        build.Target,
		"DRONE_COMMIT":               build.After,
		"DRONE_COMMIT_SHA":           build.After,
		"DRONE_COMMIT_BEFORE":        build.Before,
		"DRONE_COMMIT_AFTER":         build.After,
		"DRONE_COMMIT_REF":           build.Ref,
		"DRONE_COMMIT_BRANCH":        build.Target,
		"DRONE_COMMIT_LINK":          build.Link,
		"DRONE_COMMIT_MESSAGE":       build.Message,
		"DRONE_COMMIT_AUTHOR":        build.Author,
		"DRONE_COMMIT_AUTHOR_EMAIL":  build.AuthorEmail,
		"DRONE_COMMIT_AUTHOR_AVATAR": build.AuthorAvatar,
		"DRONE_COMMIT_AUTHOR_NAME":   build.AuthorName,
		"DRONE_BUILD_NUMBER":         fmt.Sprint(build.Number),
		"DRONE_BUILD_PARENT":         fmt.Sprint(build.Parent),
		"DRONE_BUILD_EVENT":          build.Event,
		"DRONE_BUILD_ACTION":         build.Action,
		"DRONE_BUILD_STATUS":         "success",
		"DRONE_BUILD_DEBUG":          fmt.Sprint(build.Debug),
		"DRONE_BUILD_CREATED":        fmt.Sprint(build.Created),
		"DRONE_BUILD_STARTED":        fmt.Sprint(build.Started),
		"DRONE_BUILD_FINISHED":       fmt.Sprint(build.Finished),
		"DRONE_DEPLOY_TO":            build.Deploy,
		"DRONE_DEPLOY_ID":            fmt.Sprint(build.DeployID),
		"DRONE_BUILD_TRIGGER":        build.Trigger,

		// these are legacy configuration parameters for backward
		// compatibility with drone 0.8. These are deprecated and
		// should not be relied upon going forward.
		"CI_BUILD_NUMBER":         fmt.Sprint(build.Number),
		"CI_PARENT_BUILD_NUMBER":  fmt.Sprint(build.Parent),
		"CI_BUILD_CREATED":        fmt.Sprint(build.Created),
		"CI_BUILD_STARTED":        fmt.Sprint(build.Started),
		"CI_BUILD_FINISHED":       fmt.Sprint(build.Finished),
		"CI_BUILD_STATUS":         build.Status,
		"CI_BUILD_EVENT":          build.Event,
		"CI_BUILD_LINK":           build.Link,
		"CI_BUILD_TARGET":         build.Deploy,
		"CI_COMMIT_SHA":           build.After,
		"CI_COMMIT_REF":           build.Ref,
		"CI_COMMIT_BRANCH":        build.Target,
		"CI_COMMIT_MESSAGE":       build.Message,
		"CI_COMMIT_AUTHOR":        build.Author,
		"CI_COMMIT_AUTHOR_NAME":   build.AuthorName,
		"CI_COMMIT_AUTHOR_EMAIL":  build.AuthorEmail,
		"CI_COMMIT_AUTHOR_AVATAR": build.AuthorAvatar,
	}
	if isBuildFailing(build) {
		env["DRONE_BUILD_STATUS"] = "failure"
		env["DRONE_FAILED_STAGES"] = strings.Join(failedStages(build), ",")
	}
	if build.Started == 0 {
		env["DRONE_BUILD_STARTED"] = fmt.Sprint(time.Now().Unix())
	}
	if build.Finished == 0 {
		env["DRONE_BUILD_FINISHED"] = fmt.Sprint(time.Now().Unix())
	}
	if build.Event == drone.EventPullRequest {
		env["DRONE_PULL_REQUEST"] = re.FindString(build.Ref)
		env["DRONE_PULL_REQUEST_TITLE"] = build.Title
	}
	if strings.HasPrefix(build.Ref, "refs/tags/") {
		tag := strings.TrimPrefix(build.Ref, "refs/tags/")
		env["DRONE_TAG"] = tag
		copyenv(versions(tag), env)
		copyenv(calversions(tag), env)
	}
	return env
}

// Link returns a set of environment variables containing
// resource urls to the build.
func Link(repo *drone.Repo, build *drone.Build, system *drone.System) map[string]string {
	return map[string]string{
		"DRONE_BUILD_LINK": fmt.Sprintf(
			"%s://%s/%s/%d",
			system.Proto,
			system.Host,
			repo.Slug,
			build.Number,
		),
	}
}

// Netrc returns a set of environment variables containing
// the netrc file and credentials.
func Netrc(netrc *drone.Netrc) map[string]string {
	env := map[string]string{}
	if netrc != nil && netrc.Machine != "" {
		env["DRONE_NETRC_MACHINE"] = netrc.Machine
		env["DRONE_NETRC_USERNAME"] = netrc.Login
		env["DRONE_NETRC_PASSWORD"] = netrc.Password
		env["DRONE_NETRC_FILE"] = fmt.Sprintf(
			"machine %s login %s password %s",
			netrc.Machine,
			netrc.Login,
			netrc.Password,
		)
	}
	return env
}

// Combine is a helper function combines one or more maps of
// environment variables into a single map.
func Combine(env ...map[string]string) map[string]string {
	c := map[string]string{}
	for _, e := range env {
		for k, v := range e {
			c[k] = v
		}
	}
	return c
}

// Slice is a helper function that converts a map of environment
// variables to a slice of string values in key=value format.
func Slice(env map[string]string) []string {
	s := []string{}
	for k, v := range env {
		s = append(s, k+"="+v)
	}
	sort.Strings(s)
	return s
}

// copyenv copies environment variables from the source map
// to the destination map.
func copyenv(src, dst map[string]string) {
	for k, v := range src {
		dst[k] = v
	}
}

// helper function returns true of the build is failing.
func isBuildFailing(build *drone.Build) bool {
	if build.Status == drone.StatusError ||
		build.Status == drone.StatusFailing ||
		build.Status == drone.StatusKilled {
		return true
	}
	for _, stage := range build.Stages {
		if stage.Status == drone.StatusError ||
			stage.Status == drone.StatusFailing ||
			stage.Status == drone.StatusKilled {
			return true
		}
	}
	return false
}

// helper function returns true of the stage is failing.
func isStageFailing(stage *drone.Stage) bool {
	if stage.Status == drone.StatusError ||
		stage.Status == drone.StatusFailing ||
		stage.Status == drone.StatusKilled {
		return true
	}
	for _, step := range stage.Steps {
		if step.ErrIgnore && step.Status == drone.StatusFailing {
			continue
		}
		if step.Status == drone.StatusError ||
			step.Status == drone.StatusFailing ||
			step.Status == drone.StatusKilled {
			return true
		}
	}
	return false
}

// helper function returns the failed steps.
func failedSteps(stage *drone.Stage) []string {
	var steps []string
	for _, step := range stage.Steps {
		if step.ErrIgnore && step.Status == drone.StatusFailing {
			continue
		}
		if step.Status == drone.StatusError ||
			step.Status == drone.StatusFailing ||
			step.Status == drone.StatusKilled {
			steps = append(steps, step.Name)
		}
	}
	return steps
}

// helper function returns the failed stages.
func failedStages(build *drone.Build) []string {
	var stages []string
	for _, stage := range build.Stages {
		if stage.Status == drone.StatusError ||
			stage.Status == drone.StatusFailing ||
			stage.Status == drone.StatusKilled {
			stages = append(stages, stage.Name)
		}
	}
	return stages
}
