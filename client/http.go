// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/logger"
)

const (
	endpointNode   = "/rpc/v2/nodes/%s"
	endpointPing   = "/rpc/v2/ping"
	endpointStages = "/rpc/v2/stage"
	endpointStage  = "/rpc/v2/stage/%d"
	endpointStep   = "/rpc/v2/step/%d"
	endpointWatch  = "/rpc/v2/build/%d/watch"
	endpointBatch  = "/rpc/v2/step/%d/logs/batch"
	endpointUpload = "/rpc/v2/step/%d/logs/upload"
	endpointCard   = "/rpc/v2/step/%d/card"
)

var _ Client = (*HTTPClient)(nil)

// defaultClient is the default http.Client.
var defaultClient = &http.Client{
	CheckRedirect: func(*http.Request, []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// New returns a new runner client.
func New(endpoint, secret string, skipverify bool) *HTTPClient {
	client := &HTTPClient{
		Endpoint:   endpoint,
		Secret:     secret,
		SkipVerify: skipverify,
	}
	if skipverify {
		client.Client = &http.Client{
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
	}
	return client
}

// An HTTPClient manages communication with the runner API.
type HTTPClient struct {
	Client     *http.Client
	Logger     logger.Logger
	Dumper     logger.Dumper
	Endpoint   string
	Secret     string
	SkipVerify bool
}

// Join notifies the server the runner is joining the cluster.
func (p *HTTPClient) Join(ctx context.Context, machine string) error {
	return nil
}

// Leave notifies the server the runner is leaving the cluster.
func (p *HTTPClient) Leave(ctx context.Context, machine string) error {
	return nil
}

// Ping sends a ping message to the server to test connectivity.
func (p *HTTPClient) Ping(ctx context.Context, machine string) error {
	_, err := p.do(ctx, endpointPing, "POST", nil, nil)
	return err
}

// Request requests the next available build stage for execution.
func (p *HTTPClient) Request(ctx context.Context, args *Filter) (*drone.Stage, error) {
	src := args
	dst := new(drone.Stage)
	_, err := p.retry(ctx, endpointStages, "POST", src, dst)
	return dst, err
}

// Accept accepts the build stage for execution.
func (p *HTTPClient) Accept(ctx context.Context, stage *drone.Stage) error {
	uri := fmt.Sprintf(endpointStage+"?machine=%s", stage.ID, url.QueryEscape(stage.Machine))
	src := stage
	dst := new(drone.Stage)
	_, err := p.retry(ctx, uri, "POST", nil, dst)
	if dst != nil {
		src.Updated = dst.Updated
		src.Version = dst.Version
	}
	return err
}

// Detail gets the build stage details for execution.
func (p *HTTPClient) Detail(ctx context.Context, stage *drone.Stage) (*Context, error) {
	uri := fmt.Sprintf(endpointStage, stage.ID)
	dst := new(Context)
	_, err := p.retry(ctx, uri, "GET", nil, dst)
	return dst, err
}

// Update updates the build stage.
func (p *HTTPClient) Update(ctx context.Context, stage *drone.Stage) error {
	uri := fmt.Sprintf(endpointStage, stage.ID)
	src := stage
	dst := new(drone.Stage)
	for i, step := range src.Steps {
		// a properly implemented runner should never encounter
		// input errors. these checks are included to help
		// developers creating new runners.
		if step.Number == 0 {
			return fmt.Errorf("step[%d] missing number", i)
		}
		if step.StageID == 0 {
			return fmt.Errorf("step[%d] missing stage id", i)
		}
		if step.Status == drone.StatusRunning &&
			step.Started == 0 {
			return fmt.Errorf("step[%d] missing start time", i)
		}
	}
	_, err := p.retry(ctx, uri, "PUT", src, dst)
	if dst != nil {
		src.Updated = dst.Updated
		src.Version = dst.Version

		set := map[int]*drone.Step{}
		for _, step := range dst.Steps {
			set[step.Number] = step
		}
		for _, step := range src.Steps {
			from, ok := set[step.Number]
			if ok {
				step.ID = from.ID
				step.StageID = from.StageID
				step.Started = from.Started
				step.Stopped = from.Stopped
				step.Version = from.Version
			}
		}
	}
	return err
}

// UpdateStep updates the build step.
func (p *HTTPClient) UpdateStep(ctx context.Context, step *drone.Step) error {
	uri := fmt.Sprintf(endpointStep, step.ID)
	src := step
	dst := new(drone.Step)
	_, err := p.retry(ctx, uri, "PUT", src, dst)
	if dst != nil {
		src.Version = dst.Version
	}
	return err
}

// Watch watches for build cancellation requests.
func (p *HTTPClient) Watch(ctx context.Context, build int64) (bool, error) {
	uri := fmt.Sprintf(endpointWatch, build)
	res, err := p.retry(ctx, uri, "POST", nil, nil)
	if err != nil {
		return false, err
	}
	if res.StatusCode == 200 {
		return true, nil
	}
	return false, nil
}

// Batch batch writes logs to the build logs.
func (p *HTTPClient) Batch(ctx context.Context, step int64, lines []*drone.Line) error {
	uri := fmt.Sprintf(endpointBatch, step)
	_, err := p.do(ctx, uri, "POST", &lines, nil)
	return err
}

// Upload uploads the full logs to the server.
func (p *HTTPClient) Upload(ctx context.Context, step int64, lines []*drone.Line) error {
	uri := fmt.Sprintf(endpointUpload, step)
	_, err := p.retry(ctx, uri, "POST", &lines, nil)
	return err
}

// UploadCard uploads a card to drone server.
func (p *HTTPClient) UploadCard(ctx context.Context, step int64, card *drone.CardInput) error {
	uri := fmt.Sprintf(endpointCard, step)
	_, err := p.retry(ctx, uri, "POST", &card, nil)
	return err
}

func (p *HTTPClient) retry(ctx context.Context, path, method string, in, out interface{}) (*http.Response, error) {
	for {
		res, err := p.do(ctx, path, method, in, out)
		// do not retry on Canceled or DeadlineExceeded
		if err := ctx.Err(); err != nil {
			p.logger().Tracef("http: context canceled")
			return res, err
		}
		// do not retry on optimisitic lock errors
		if err == ErrOptimisticLock {
			return res, err
		}
		if res != nil {
			// Check the response code. We retry on 500-range
			// responses to allow the server time to recover, as
			// 500's are typically not permanent errors and may
			// relate to outages on the server side.
			if res.StatusCode > 501 {
				p.logger().Tracef("http: server error: re-connect and re-try: %s", err)
				time.Sleep(time.Second * 10)
				continue
			}
			// We also retry on 204 no content response codes,
			// used by the server when a long-polling request
			// is intentionally disconnected and should be
			// automatically reconnected.
			if res.StatusCode == 204 {
				p.logger().Tracef("http: no content returned: re-connect and re-try")
				time.Sleep(time.Second * 10)
				continue
			}
		} else if err != nil {
			p.logger().Tracef("http: request error: %s", err)
			time.Sleep(time.Second * 10)
			continue
		}
		return res, err
	}
}

// do is a helper function that posts a signed http request with
// the input encoded and response decoded from json.
func (p *HTTPClient) do(ctx context.Context, path, method string, in, out interface{}) (*http.Response, error) {
	var buf bytes.Buffer

	// marshal the input payload into json format and copy
	// to an io.ReadCloser.
	if in != nil {
		json.NewEncoder(&buf).Encode(in)
	}

	endpoint := p.Endpoint + path
	req, err := http.NewRequest(method, endpoint, &buf)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	// the request should include the secret shared between
	// the agent and server for authorization.
	req.Header.Add("X-Drone-Token", p.Secret)

	if p.Dumper != nil {
		p.Dumper.DumpRequest(req)
	}

	res, err := p.client().Do(req)
	if res != nil {
		defer func() {
			// drain the response body so we can reuse
			// this connection.
			io.Copy(ioutil.Discard, io.LimitReader(res.Body, 4096))
			res.Body.Close()
		}()
	}
	if err != nil {
		return res, err
	}

	if p.Dumper != nil {
		p.Dumper.DumpResponse(res)
	}

	// if the response body return no content we exit
	// immediately. We do not read or unmarshal the response
	// and we do not return an error.
	if res.StatusCode == 204 {
		return res, nil
	}

	// Check the response for a 409 conflict. This indicates an
	// optimistic lock error, in which case multiple clients may
	// be attempting to update the same record. Convert this error
	// code to a proper error.
	if res.StatusCode == 409 {
		return nil, ErrOptimisticLock
	}

	// else read the response body into a byte slice.
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res, err
	}

	if res.StatusCode > 299 {
		// if the response body includes an error message
		// we should return the error string.
		if len(body) != 0 {
			return res, errors.New(
				string(body),
			)
		}
		// if the response body is empty we should return
		// the default status code text.
		return res, errors.New(
			http.StatusText(res.StatusCode),
		)
	}
	if out == nil {
		return res, nil
	}
	return res, json.Unmarshal(body, out)
}

// client is a helper funciton that returns the default client
// if a custom client is not defined.
func (p *HTTPClient) client() *http.Client {
	if p.Client == nil {
		return defaultClient
	}
	return p.Client
}

// logger is a helper funciton that returns the default logger
// if a custom logger is not defined.
func (p *HTTPClient) logger() logger.Logger {
	if p.Logger == nil {
		return logger.Discard()
	}
	return p.Logger
}
