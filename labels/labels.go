// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package labels

import (
	"fmt"
	"time"

	"github.com/drone/drone-go/drone"
)

// now returns the current time.
var now = time.Now

// FromRepo returns container labels derived from the
// Repository metadata.
func FromRepo(v *drone.Repo) map[string]string {
	return map[string]string{
		"io.drone.repo.namespace": v.Namespace,
		"io.drone.repo.name":      v.Name,
		"io.drone.repo.slug":      v.Slug,
	}
}

// FromBuild returns container labels derived from the
// Build metadata.
func FromBuild(v *drone.Build) map[string]string {
	return map[string]string{
		"io.drone.build.number": fmt.Sprint(v.Number),
	}
}

// FromStage returns container labels derived from the
// Stage metadata.
func FromStage(v *drone.Stage) map[string]string {
	return map[string]string{
		"io.drone.stage.name":   v.Name,
		"io.drone.stage.number": fmt.Sprint(v.Number),
	}
}

// FromStep returns container labels derived from the
// Step metadata.
func FromStep(v *drone.Step) map[string]string {
	return map[string]string{
		"io.drone.step.number": fmt.Sprint(v.Number),
		"io.drone.step.name":   v.Name,
	}
}

// FromSystem returns container labels derived from the
// System metadata.
func FromSystem(v *drone.System) map[string]string {
	return map[string]string{
		"io.drone":                "true",
		"io.drone.protected":      "false",
		"io.drone.system.host":    v.Host,
		"io.drone.system.proto":   v.Proto,
		"io.drone.system.version": v.Version,
	}
}

// WithTimeout returns container labels that define
// timeout and expiration values.
func WithTimeout(v *drone.Repo) map[string]string {
	return map[string]string{
		"io.drone.ttl":     fmt.Sprint(time.Duration(v.Timeout) * time.Minute),
		"io.drone.expires": fmt.Sprint(now().Add(time.Duration(v.Timeout)*time.Minute + time.Hour).Unix()),
		"io.drone.created": fmt.Sprint(now().Unix()),
	}
}

// Combine is a helper function combines one or more maps of
// labels into a single map.
func Combine(labels ...map[string]string) map[string]string {
	c := map[string]string{}
	for _, e := range labels {
		for k, v := range e {
			c[k] = v
		}
	}
	return c
}
