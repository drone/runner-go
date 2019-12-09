// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package replacer provides helper functions to mask
// secrets in data streams.
package replacer

import (
	"io"
	"strings"

	"github.com/drone/runner-go/pipeline/runtime/driver"
)

const maskedF = "[secret:%s]"

// Replacer is an io.Writer that finds and masks sensitive data.
type Replacer struct {
	w io.WriteCloser
	r *strings.Replacer
}

// New returns a replacer that wraps io.Writer w.
func New(w io.WriteCloser, secrets []driver.Secret) io.WriteCloser {
	var oldnew []string
	for _, secret := range secrets {
		if len(secret.GetValue()) == 0 || secret.IsMasked() == false {
			continue
		}
		// name := strings.ToLower(secret.GetName())
		// masked := fmt.Sprintf(maskedF, name)

		// TODO temporarily revert back to masking secrets
		// using the asterisk symbol due to confusion when
		// masking with [secret:name]
		masked := "******"
		oldnew = append(oldnew, string(secret.GetValue()))
		oldnew = append(oldnew, masked)
	}
	if len(oldnew) == 0 {
		return w
	}
	return &Replacer{
		w: w,
		r: strings.NewReplacer(oldnew...),
	}
}

// Write writes p to the base writer. The method scans for any
// sensitive data in p and masks before writing.
func (r *Replacer) Write(p []byte) (n int, err error) {
	_, err = r.w.Write([]byte(r.r.Replace(string(p))))
	return len(p), err
}

// Close closes the base writer.
func (r *Replacer) Close() error {
	return r.w.Close()
}
