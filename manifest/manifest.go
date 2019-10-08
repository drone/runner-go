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

	// ConcurrentResource is a resource with concurrency limits.
	ConcurrentResource interface {
		Resource
		GetConcurrency() Concurrency
	}

	// DependantResource is a resource with runtime dependencies.
	DependantResource interface {
		Resource
		GetDependsOn() []string
	}

	// PlatformResource is a resource with platform requirements.
	PlatformResource interface {
		Resource
		GetPlatform() Platform
	}

	// RoutedResource is a resource that can be routed to
	// specific build nodes.
	RoutedResource interface {
		Resource
		GetNodes() map[string]string
	}

	// TriggeredResource is a resource with trigger rules.
	TriggeredResource interface {
		Resource
		GetTrigger() Conditions
	}

	// RawResource is a raw encoded resource with the common
	// metadata extracted.
	RawResource struct {
		Version     string
		Kind        string
		Type        string
		Name        string
		Deps        []string `yaml:"depends_on"`
		Node        map[string]string
		Concurrency Concurrency
		Platform    Platform
		Data        []byte `yaml:"-"`
	}
)
