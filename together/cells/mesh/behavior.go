// Tideland Go Library - Together - Cells - Mesh
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package mesh // import "tideland.dev/go/together/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/together/cells/event"
)

//--------------------
// EMITTER
//--------------------

// Emitter describes a behavior to emit events to subscribers. An instance
// is passed during initialization.
type Emitter interface {
	// Mesh returns the mesh of the emitter.
	Mesh() *Mesh

	// Subscribers returns the the IDs of the subscriber cells.
	Subscribers() []string

	// Emit emits the given event to the given subscriber if it exists.
	Emit(id string, evt *event.Event) error

	// Broadcast emits the given event to all subscribers.
	Broadcast(evt *event.Event) error

	// Self emits the given event back to the cell itself.
	Self(evt *event.Event) error
}

//--------------------
// BEHAVIOR
//--------------------

// Behavior is the interface that has to be implemented for event
// processing inside the cells.
type Behavior interface {
	// ID returns the individual identifier of a behavior instance.
	// Behaviors can be deployed multiple times as long as these return
	// different identifiers.
	ID() string

	// Init is called by the cells to initialize the behavior.
	// The passed emitter is for emitting events to subscribers.
	Init(emitter Emitter) error

	// Terminate is called when a cell is stopped.
	Terminate() error

	// Process is called to process the given event.
	Process(evt *event.Event) error

	// Recover is called in case of an error or panic during the processing
	// of an event. Here the behavior can check if it can recover and establish
	// a valid state. If it's not possible the implementation has to return
	// an error documenting the reason.
	Recover(err interface{}) error
}

// EOF
