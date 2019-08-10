// Tideland Go Library - Together - Cells - Behaviors
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
	"tideland.dev/go/trace/failure"
)

//--------------------
// MESH ROUTER BEHAVIOR
//--------------------

// MeshRouter is a function type determining which cells shall get the event.
type MeshRouter func(evt *event.Event) []string

// meshRouterBehavior check for each received event which subscriber will
// get it based on the router function.
type meshRouterBehavior struct {
	id      string
	routeTo MeshRouter
	msh     *mesh.Mesh
}

// NewMeshRouterBehavior creates a mesh router behavior using the passed function
// to determine to which cell the received event shall be re-emitted.
func NewMeshRouterBehavior(id string, router MeshRouter, msh *mesh.Mesh) mesh.Behavior {
	return &meshRouterBehavior{
		id:      id,
		routeTo: router,
		msh:     msh,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *meshRouterBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *meshRouterBehavior) Init(emitter mesh.Emitter) error {
	return nil
}

// Terminate the behavior.
func (b *meshRouterBehavior) Terminate() error {
	return nil
}

// Process emits the event to those ids returned by the router function.
func (b *meshRouterBehavior) Process(evt *event.Event) error {
	ids := b.routeTo(evt)
	var errs []error
	for _, id := range ids {
		errs = append(errs, b.msh.Emit(id, evt))
	}
	return failure.Collect(errs...)
}

// Recover from an error.
func (b *meshRouterBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
