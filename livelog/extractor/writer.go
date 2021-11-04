package extractor

import (
	"encoding/base64"
	"io"
	"regexp"
)

var re = regexp.MustCompile("#((.*?)#)")

type Writer struct {
	base io.Writer
	file []byte
}

func New(w io.Writer) *Writer {
	return &Writer{w, nil}
}

func (e *Writer) Write(p []byte) (n int, err error) {
	card := re.FindStringSubmatch(string(p))
	if card == nil {
		return e.base.Write(p)
	}

	data, err := base64.StdEncoding.DecodeString(card[len(card)-1:][0])
	if err == nil {
		e.file = data
	}
	// remove encoded string for logs
	return e.base.Write([]byte(""))
}

func (e *Writer) File() ([]byte, bool) {
	if len(e.file) > 0 {
		return e.file, true
	} else {
		return nil, false
	}
}
