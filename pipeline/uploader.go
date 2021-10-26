package pipeline

import (
	"context"
	"io"
)

type Uploader interface {
	UploadCard(context.Context, io.ReadCloser, *State, string) error
}

func NopUploader() Uploader {
	return new(nopUploader)
}

type nopUploader struct{}

func (*nopUploader) UploadCard(context.Context, io.ReadCloser, *State, string) error { return nil }
