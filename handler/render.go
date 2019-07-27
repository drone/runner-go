// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/drone/runner-go/handler/template"
)

// renderJSON writes the json-encoded representation of v to
// the response body.
func renderJSON(w http.ResponseWriter, v interface{}) {
	for k, v := range noCacheHeaders {
		w.Header().Set(k, v)
	}
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

// render writes the template to the response body.
func render(w http.ResponseWriter, t string, v interface{}) {
	w.Header().Set("Content-Type", "text/html")
	template.T.ExecuteTemplate(w, t, v)
}
