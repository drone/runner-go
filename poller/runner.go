// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package runtime

import (
	"context"

	"github.com/drone/drone-go/drone"
)

// Runner is responsible for running the pipeline. It is invoked
// by the poller when a pipeline is received by the remote system.
type Runner interface {
	Run(context.Context, *drone.Stage) error
}
