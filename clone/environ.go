// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package clone

// Config provides the Git Configuration parameters.
type Config struct {
	User       User
	Trace      bool
	SkipVerify bool
}

// User provides the Git user parameters.
type User struct {
	Name  string
	Email string
}

// Environ returns a set of global Git environment variables,
// from the configuration input.
func Environ(config Config) map[string]string {
	environ := map[string]string{
		"GIT_AUTHOR_NAME":     "drone",
		"GIT_AUTHOR_EMAIL":    "noreply@drone",
		"GIT_COMMITTER_NAME":  "drone",
		"GIT_COMMITTER_EMAIL": "noreply@drone",
		"GIT_TERMINAL_PROMPT": "0",
	}
	if s := config.User.Name; s != "" {
		environ["GIT_AUTHOR_NAME"] = s
		environ["GIT_COMMITTER_NAME"] = s
	}
	if s := config.User.Email; s != "" {
		environ["GIT_AUTHOR_EMAIL"] = s
		environ["GIT_COMMITTER_EMAIL"] = s
	}
	if config.Trace {
		environ["GIT_TRACE"] = "true"
	}
	if config.SkipVerify {
		environ["GIT_SSL_NO_VERIFY"] = "true"
	}
	return environ
}
