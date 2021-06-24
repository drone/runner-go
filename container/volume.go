// Copyright 2021 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package container

import (
	"path/filepath"
	"strings"
)

// IsRestrictedVolume is helper function that
// returns true if mounting the volume is restricted for un-trusted containers.
func IsRestrictedVolume(path string) bool {
	path, err := filepath.Abs(path)
	if err != nil {
		return true
	}

	path = strings.ToLower(path)

	switch {
	case path == "/":
	case path == "/etc":
	case path == "/etc/docker" || strings.HasPrefix(path, "/etc/docker/"):
	case path == "/var":
	case path == "/var/run" || strings.HasPrefix(path, "/var/run/"):
	case path == "/proc" || strings.HasPrefix(path, "/proc/"):
	case path == "/usr/local/bin" || strings.HasPrefix(path, "/usr/local/bin/"):
	case path == "/usr/local/sbin" || strings.HasPrefix(path, "/usr/local/sbin/"):
	case path == "/usr/bin" || strings.HasPrefix(path, "/usr/bin/"):
	case path == "/bin" || strings.HasPrefix(path, "/bin/"):
	case path == "/mnt" || strings.HasPrefix(path, "/mnt/"):
	case path == "/mount" || strings.HasPrefix(path, "/mount/"):
	case path == "/media" || strings.HasPrefix(path, "/media/"):
	case path == "/sys" || strings.HasPrefix(path, "/sys/"):
	case path == "/dev" || strings.HasPrefix(path, "/dev/"):
	default:
		return false
	}

	return true
}
