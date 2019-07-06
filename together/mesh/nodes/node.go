// Tideland Go Library - Together - Mesh - Nodes
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package nodes // import "tideland.dev/go/together/mesh/nodes"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/together/actor"
	"tideland.dev/go/trace/errors"
)

//--------------------
// NODE
//--------------------

// node runs one behavior for the processing of events and emitting of
// resulting events.
type node struct {
	behavior    Behavior
	subscribers map[string]*node
	act         *actor.Actor
}

// newNode creates a new node running the given behavior in a goroutine.
func newNode(behavior Behavior) (*node, error) {
	n := &node{
		behavior:    behavior,
		subscribers: map[string]*node{},
		act: actor.New(
			actor.WithQueueLen(32),
			actor.WithRecoverer(behavior.Recover),
		).Go(),
	}
	err := n.behavior.Init(n)
	if err != nil {
		// Stop the actor with the annotated error.
		return nil, n.act.Stop(errors.Annotate(err, ErrNodeInit, msgNodeInit, behavior.ID()))
	}
	return n, nil
}

// Emit allows a behavior to emit events to its subsribers.
func (n *node) Emit(event Event) error {
	if aerr := n.act.DoAsync(func() error {
		for _, subscriber := range n.subscribers {
			subscriber.process(event)
		}
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrNodeBackend, msgNodeBackend, n.behavior.ID())
	}
	return nil
}

// subscribe adds processor engines to the subscribers of this engine.
func (n *node) subscribe(subscribers []*node) error {
	if aerr := n.act.DoAsync(func() error {
		for _, subscriber := range subscribers {
			n.subscribers[subscriber.behavior.ID()] = subscriber
		}
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrNodeBackend, msgNodeBackend, n.behavior.ID())
	}
	return nil
}

// process lets the node behavior process the event asynchronously.
func (n *node) process(event Event) error {
	if aerr := n.act.DoAsync(func() error {
		n.behavior.Process(event)
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrNodeBackend, msgNodeBackend, n.behavior.ID())
	}
	return nil
}

// terminate tells the processor to end and replaces it with a dummy.
func (n *node) terminate() error {
	var err error
	if aerr := n.act.DoSync(func() error {
		err = n.behavior.Terminate()
		n.behavior = &dummyBehavior{n.behavior.ID()}
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrNodeBackend, msgNodeBackend, n.behavior.ID())
	}
	return err
}

// stop ends the actor.
func (n *node) stop(err error) error {
	return n.act.Stop(err)
}

//--------------------
// DUMMY BEHAVIOR
//--------------------

// dummyBehavior will be used by a node while it's shutting down.
type dummyBehavior struct {
	id string
}

func (db *dummyBehavior) ID() string {
	return db.id
}

func (db *dummyBehavior) Init(emitter Emitter) error {
	return nil
}

func (db *dummyBehavior) Terminate() error {
	return nil
}

func (db *dummyBehavior) Process(event Event) {
}

func (db *dummyBehavior) Recover(r interface{}) error {
	return nil
}

// EOF
