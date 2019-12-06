// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package runtime

import (
	"fmt"
	"io"
	"strings"
)

const maskedF = "[secret:%s]"

// replacer is an io.Writer that finds and masks sensitive data.
type replacer struct {
	w io.WriteCloser
	r *strings.Replacer
}

// newReplacer returns a replacer that wraps writer w.
func newReplacer(w io.WriteCloser, secrets []Secret) io.WriteCloser {
	var oldnew []string
	for _, secret := range secrets {
		if len(secret.GetValue()) == 0 || secret.IsMasked() == false {
			continue
		}
		name := strings.ToLower(secret.GetName())
		masked := fmt.Sprintf(maskedF, name)
		oldnew = append(oldnew, string(secret.GetValue()))
		oldnew = append(oldnew, masked)
	}
	if len(oldnew) == 0 {
		return w
	}
	return &replacer{
		w: w,
		r: strings.NewReplacer(oldnew...),
	}
}

// Write writes p to the base writer. The method scans for any
// sensitive data in p and masks before writing.
func (r *replacer) Write(p []byte) (n int, err error) {
	_, err = r.w.Write([]byte(r.r.Replace(string(p))))
	return len(p), err
}

// Close closes the base writer.
func (r *replacer) Close() error {
	return r.w.Close()
}
