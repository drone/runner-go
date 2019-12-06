// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package router

import (
	"net/http"

	"github.com/drone/runner-go/handler"
	"github.com/drone/runner-go/handler/static"
	hook "github.com/drone/runner-go/logger/history"
	"github.com/drone/runner-go/pipeline/reporter/history"

	"github.com/99designs/basicauth-go"
)

// Config provides router configuration.
type Config struct {
	Username string
	Password string
	Realm    string
}

// New returns a new route handler.
func New(tracer *history.History, history *hook.Hook, config Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handler.HandleHealth(tracer))

	// omit dashboard handlers when no password configured.
	if config.Password == "" {
		return mux
	}

	// middleware to require basic authentication.
	auth := basicauth.New(config.Realm, map[string][]string{
		config.Username: {config.Password},
	})

	// handler to serve static assets for the dashboard.
	fs := http.FileServer(static.New())

	// dashboard handles.
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.Handle("/logs", auth(handler.HandleLogHistory(history)))
	mux.Handle("/view", auth(handler.HandleStage(tracer, history)))
	mux.Handle("/", auth(handler.HandleIndex(tracer)))
	return mux
}
