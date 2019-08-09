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
// ROUTER BEHAVIOR
//--------------------

// Router is a function type determining which cells shall get the event.
type Router func(evt *event.Event) []string

// routerBehavior check for each received event which subscriber will
// get it based on the router function.
type routerBehavior struct {
	id      string
	routeTo Router
	msh     *mesh.Mesh
}

// NewRouterBehavior creates a router behavior using the passed function
// to determine to which subscriber the received event will be emitted.
func NewRouterBehavior(id string, router Router, msh *mesh.Mesh) mesh.Behavior {
	return &routerBehavior{
		id:      id,
		routeTo: router,
		msh:     msh,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *routerBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *routerBehavior) Init(emitter mesh.Emitter) error {
	return nil
}

// Terminate the behavior.
func (b *routerBehavior) Terminate() error {
	return nil
}

// Process emits the event to those ids returned by the router function.
func (b *routerBehavior) Process(evt *event.Event) error {
	ids := b.routeTo(evt)
	var errs []error
	for _, id := range ids {
		errs = append(errs, b.msh.Emit(id, evt))
	}
	return failure.Collect(errs...)
}

// Recover from an error.
func (b *routerBehavior) Recover(err interface{}) error {
	return nil
}

// EOF
