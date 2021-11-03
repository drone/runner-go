// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package livelog

import (
	"bufio"
	"encoding/base64"
	"io"
	"io/ioutil"
	"regexp"
)

var re = regexp.MustCompile("#((.*?)#)")

// Copy copies from src to dst and removes until either EOF
// is reached on src or an error occurs.
func Copy(dst io.Writer, src io.ReadCloser) error {
	r := bufio.NewReader(src)
	for {
		bytes, err := r.ReadBytes('\n')
		// check logs for card data
		card := re.FindStringSubmatch(string(bytes))
		if card != nil {
			data, err := base64.StdEncoding.DecodeString(card[len(card)-1:][0])
			if err == nil {
				_ = ioutil.WriteFile("/tmp/card.json", data, 0644)
			}
			continue
		}
		if _, err := dst.Write(bytes); err != nil {
			return err
		}
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
	}
}
