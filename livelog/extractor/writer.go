package extractor

import (
	"bytes"
	"encoding/base64"
	"io"
	"regexp"
)

const Esc = "\u001B"

var (
	prefix = Esc + "]1338;"
	re     = regexp.MustCompilePOSIX("\u001B]1338;((.*?)\u001B]0m)")
)

type Writer struct {
	base io.Writer
	file []byte
}

func New(w io.Writer) *Writer {
	return &Writer{w, nil}
}

func (e *Writer) Write(p []byte) (n int, err error) {
	if bytes.HasPrefix(p, []byte(prefix)) == false {
		return e.base.Write(p)
	}
	card := re.FindStringSubmatch(string(p))
	if len(card) != 0 {
		data, err := base64.StdEncoding.DecodeString(card[len(card)-1:][0])
		if err == nil {
			e.file = data
		}
		return e.base.Write([]byte(""))
	}
	return e.base.Write(p)
}

func (e *Writer) File() ([]byte, bool) {
	if len(e.file) > 0 {
		return e.file, true
	} else {
		return nil, false
	}
}
