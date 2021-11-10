package pipeline

import (
	"context"
)

type Uploader interface {
	UploadCard(context.Context, []byte, *State, string) error
}

func NopUploader() Uploader {
	return new(nopUploader)
}

type nopUploader struct{}

func (*nopUploader) UploadCard(context.Context, []byte, *State, string) error { return nil }
