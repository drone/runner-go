package uploader

import (
	"context"
	"encoding/json"
	"io"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/client"
	"github.com/drone/runner-go/internal"
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

func (s *Upload) UploadCard(ctx context.Context, r io.ReadCloser, state *pipeline.State, stepName string) error {
	src := state.Find(stepName)
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	card := drone.CardInput{}
	err = json.Unmarshal(bytes, &card)
	if err != nil {
		return err
	}
	err = s.client.UploadCard(ctx, src.ID, &card)
	if err != nil {
		return err
	}
	// update step schema
	state.Lock()
	src.Schema = card.Schema
	cpy := internal.CloneStep(src)
	state.Unlock()
	err = s.client.UpdateStep(ctx, cpy)
	if err == nil {
		state.Lock()
		internal.MergeStep(cpy, src)
		state.Unlock()
	}
	return nil
}
