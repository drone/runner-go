// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package logger

import (
	"bytes"
	"net/http"
	"os"
	"testing"
)

func TestStandardDumper(t *testing.T) {
	d := StandardDumper(true)
	if s, ok := d.(*standardDumper); !ok {
		t.Errorf("Expect standard dumper")
	} else if s.out != os.Stdout {
		t.Errorf("Expect standard dumper set to stdout")
	}
}

func TestDiscardDumper(t *testing.T) {
	d := DiscardDumper()
	if _, ok := d.(*discardDumper); !ok {
		t.Errorf("Expect discard dumper")
	}
}

func TestStandardDumper_DumpRequest(t *testing.T) {
	buf := new(bytes.Buffer)
	r, _ := http.NewRequest("GET", "http://example.com", nil)
	d := StandardDumper(true).(*standardDumper)
	d.out = buf
	d.DumpRequest(r)

	want := "GET / HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n"
	got := buf.String()
	if got != want {
		t.Errorf("Got dumped request %q", got)
	}
}

func TestStandardDumper_DumpResponse(t *testing.T) {
	buf := new(bytes.Buffer)
	r := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
	}
	d := StandardDumper(true).(*standardDumper)
	d.out = buf
	d.DumpResponse(r)

	want := "HTTP/1.0 200 OK\r\nContent-Length: 0\r\n\r\n"
	got := buf.String()
	if got != want {
		t.Errorf("Got dumped request %q", got)
	}
}
