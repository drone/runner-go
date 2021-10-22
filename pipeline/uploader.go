package pipeline

import (
	"context"
	"io"
)

type Uploader interface {
	UploadCard(ctx context.Context, r io.ReadCloser, step int64) error
}

func NopUploader() Uploader {
	return new(nopUploader)
}

type nopUploader struct{}

func (*nopUploader) UploadCard(ctx context.Context, r io.ReadCloser, step int64) error { return nil }
