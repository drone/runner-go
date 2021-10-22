package uploader

import (
	"context"
	"encoding/json"
	"io"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/client"
	"github.com/drone/runner-go/pipeline"
)

var _ pipeline.Uploader = (*Upload)(nil)

type Upload struct {
	client client.Client
}

func New(client client.Client) *Upload {
	return &Upload{
		client: client,
	}
}

func (s *Upload) UploadCard(ctx context.Context, r io.ReadCloser, step int64) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	card := drone.CardInput{}
	err = json.Unmarshal(bytes, &card)
	if err != nil {
		return err
	}
	err = s.client.UploadCard(ctx, step, &card)
	if err != nil {
		return err
	}
	return nil
}
