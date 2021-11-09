package extractor

import (
	"bytes"
	"encoding/base64"
	"io"
	"regexp"
)

var (
	prefix = []byte("\u001B]1338;")
	suffix = []byte("\u001B]0m")
	re     = regexp.MustCompilePOSIX("\u001B]1338;((.*?)\u001B]0m)")
)

type Writer struct {
	base    io.Writer
	file    []byte
	chunked bool
}

func New(w io.Writer) *Writer {
	return &Writer{w, nil, false}
}

func (e *Writer) Write(p []byte) (n int, err error) {
	if bytes.HasPrefix(p, prefix) == false && e.chunked == false {
		return e.base.Write(p)
	}
	n = len(p)

	// if the data does not include the ansi suffix,
	// it exceeds the size of the buffer and is chunked.
	e.chunked = !bytes.Contains(p, suffix)

	// trim the ansi prefix and suffix from the data,
	// and also trim any spacing or newlines that could
	// cause confusion.
	p = bytes.TrimSpace(p)
	p = bytes.TrimPrefix(p, prefix)
	p = bytes.TrimSuffix(p, suffix)

	e.file = append(e.file, p...)
	return n, nil
}

func (e *Writer) File() ([]byte, bool) {
	if len(e.file) == 0 {
		return nil, false
	}
	data, err := base64.StdEncoding.DecodeString(string(e.file))
	if err != nil {
		return nil, false
	}
	return data, true
}
