// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package runtime

import (
	"context"
	"io"

	"github.com/drone/drone-go/drone"
	"github.com/drone/runner-go/manifest"
	"github.com/drone/runner-go/secret"
)

type (
	// CompilerArgs provides compiler arguments.
	CompilerArgs struct {
		// Manifest provides the parsed manifest.
		Manifest *manifest.Manifest

		// Pipeline provides the parsed pipeline. This pipeline is
		// the compiler source and is converted to the intermediate
		// representation by the Compile method.
		Pipeline manifest.Resource

		// Build provides the compiler with stage information that
		// is converted to environment variable format and passed to
		// each pipeline step. It is also used to clone the commit.
		Build *drone.Build

		// Stage provides the compiler with stage information that
		// is converted to environment variable format and passed to
		// each pipeline step.
		Stage *drone.Stage

		// Repo provides the compiler with repo information. This
		// repo information is converted to environment variable
		// format and passed to each pipeline step. It is also used
		// to clone the repository.
		Repo *drone.Repo

		// System provides the compiler with system information that
		// is converted to environment variable format and passed to
		// each pipeline step.
		System *drone.System

		// Netrc provides netrc parameters that can be used by the
		// default clone step to authenticate to the remote
		// repository.
		Netrc *drone.Netrc

		// Secret returns a named secret value that can be injected
		// into the pipeline step.
		Secret secret.Provider
	}

	// Compiler compiles the Yaml configuration file to an
	// intermediate representation optimized for simple execution.
	Compiler interface {
		Compile(context.Context, CompilerArgs) Spec
	}

	// LinterArgs provides linting arguments.
	LinterArgs struct {
		Build *drone.Build
		Repo  *drone.Repo
	}

	// Linter lints the Yaml configuration file and returns an
	// error if one or more linting rules fails.
	Linter interface {
		Lint(context.Context, LinterArgs) error
	}

	// Engine is the interface that must be implemented by a
	// pipeline execution engine.
	Engine interface {
		// Setup the pipeline environment.
		Setup(context.Context, Spec) error

		// Destroy the pipeline environment.
		Destroy(context.Context, Spec) error

		// Run runs the pipeline step.
		Run(context.Context, Spec, Step, io.Writer) (*State, error)
	}

	// Spec is an interface that must be implemented by all
	// pipeline specifications.
	Spec interface {
		// StepAt returns the step at the specified index.
		StepAt(int) Step

		// StepLen returns the number of steps.
		StepLen() int
	}

	// Step is an interface that must be implemented by all
	// pipeline steps.
	Step interface {
		// GetName returns the step name.
		GetName() string

		// GetDependencies returns the step dependencies
		// used to define an execution graph.
		GetDependencies() []string

		// GetEnviron returns the step environment variables.
		GetEnviron() map[string]string

		// SetEnviron updates the step environment variables.
		SetEnviron(map[string]string)

		// GetErrPolicy returns the step error policy.
		GetErrPolicy() ErrPolicy

		// GetRunPolicy returns the step run policy.
		GetRunPolicy() RunPolicy

		// GetSecretAt returns the secret at the specified
		// index.
		GetSecretAt(int) Secret

		// GetSecretLen returns the number of secrets.
		GetSecretLen() int

		// IsDetached returns true if the step is detached
		// and executed in the background.
		IsDetached() bool

		// Clone returns a copy of the Step.
		Clone() Step

		// GetImage returns the image used in the step.
		GetImage() string
	}

	// State reports the step state.
	State struct {
		// ExitCode returns the exit code of the exited step.
		ExitCode int

		// GetExited reports whether the step has exited.
		Exited bool

		// OOMKilled reports whether the step has been
		// killed by the process manager.
		OOMKilled bool
	}

	// Secret is an interface that must be implemented
	// by all pipeline secrets.
	Secret interface {
		// GetName returns the secret name.
		GetName() string

		// GetValue returns the secret value.
		GetValue() string

		// IsMasked returns true if the secret value should
		// be masked. If true the secret value is masked in
		// the logs.
		IsMasked() bool
	}
)
