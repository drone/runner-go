// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

// Package manifest provides definitions for the Yaml schema.
package manifest

// Resource enums.
const (
	KindApproval   = "approval"
	KindDeployment = "deployment"
	KindPipeline   = "pipeline"
	KindSecret     = "secret"
	KindSignature  = "signature"
)

type (
	// Manifest is a collection of Drone resources.
	Manifest struct {
		Resources []Resource
	}

	// Resource represents a Drone resource.
	Resource interface {
		GetVersion() string
		GetKind() string
		GetType() string
		GetName() string
	}

	// DependantResource is a resoure with runtime dependencies.
	DependantResource interface {
		Resource
		GetDependsOn() []string
	}

	// PlatformResource is a resoure with platform requirements.
	PlatformResource interface {
		Resource
		GetPlatform() Platform
	}

	// TriggeredResource is a resoure with trigger rules.
	TriggeredResource interface {
		Resource
		GetTrigger() Conditions
	}

	// RawResource is a raw encoded resource with the common
	// metadata extracted.
	RawResource struct {
		Version  string
		Kind     string
		Type     string
		Name     string
		Deps     []string `yaml:"depends_on"`
		Platform Platform
		Data     []byte `yaml:"-"`
	}
)
