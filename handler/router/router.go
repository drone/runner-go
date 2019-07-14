// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Parity Public License
// that can be found in the LICENSE file.

package router

import (
	"net/http"

	"github.com/drone/runner-go/handler"
	"github.com/drone/runner-go/handler/static"
	hook "github.com/drone/runner-go/logger/history"
	"github.com/drone/runner-go/pipeline/history"

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
	// middleware to require basic authentication.
	auth := basicauth.New(config.Realm, map[string][]string{
		config.Username: {config.Password},
	})

	// handler to serve static assets for the dashboard.
	fs := http.FileServer(static.New())

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/healthz", handler.HandleHealth(tracer))
	mux.Handle("/logs", auth(handler.HandleLogHistory(history)))
	mux.Handle("/view", auth(handler.HandleStage(tracer)))
	mux.Handle("/", auth(handler.HandleIndex(tracer)))
	return mux
}
