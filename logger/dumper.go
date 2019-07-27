// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package logger

import (
	"io"
	"net/http"
	"net/http/httputil"
	"os"
)

// Dumper dumps the http.Request and http.Response
// message payload for debugging purposes.
type Dumper interface {
	DumpRequest(*http.Request)
	DumpResponse(*http.Response)
}

// DiscardDumper returns a no-op dumper.
func DiscardDumper() Dumper {
	return new(discardDumper)
}

type discardDumper struct{}

func (*discardDumper) DumpRequest(*http.Request)   {}
func (*discardDumper) DumpResponse(*http.Response) {}

// StandardDumper returns a standard dumper.
func StandardDumper(body bool) Dumper {
	return &standardDumper{out: os.Stdout, body: body}
}

type standardDumper struct {
	body bool
	out  io.Writer
}

func (s *standardDumper) DumpRequest(req *http.Request) {
	dump, _ := httputil.DumpRequestOut(req, s.body)
	s.out.Write(dump)
}

func (s *standardDumper) DumpResponse(res *http.Response) {
	dump, _ := httputil.DumpResponse(res, s.body)
	s.out.Write(dump)
}
