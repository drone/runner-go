// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package pipeline

import "context"

// A Reporter reports the pipeline status.
type Reporter interface {
	// ReportStage reports the stage status.
	ReportStage(context.Context, *State) error

	// ReportStep reports the named step status.
	ReportStep(context.Context, *State, string) error
}

// NopReporter returns a noop reporter.
func NopReporter() Reporter {
	return new(nopReporter)
}

type nopReporter struct{}

func (*nopReporter) ReportStage(context.Context, *State) error        { return nil }
func (*nopReporter) ReportStep(context.Context, *State, string) error { return nil }
