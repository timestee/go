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
	"tideland.dev/go/together/actor"
	"tideland.dev/go/trace/errors"
)

//--------------------
// CELL
//--------------------

// cell runs a behavior for the processing of events and emitting of
// resulting events.
type cell struct {
	behavior    Behavior
	subscribers map[string]*cell
	act         *actor.Actor
}

// newCell creates a new cell running the given behavior in a goroutine.
func newCell(behavior Behavior) (*cell, error) {
	c := &cell{
		behavior:    behavior,
		subscribers: map[string]*cell{},
		act: actor.New(
			actor.WithQueueLen(32),
			actor.WithRecoverer(behavior.Recover),
		).Go(),
	}
	err := c.behavior.Init(c)
	if err != nil {
		// Stop the actor with the annotated error.
		return nil, c.act.Stop(errors.Annotate(err, ErrCellInit, msgCellInit, behavior.ID()))
	}
	return c, nil
}

// Emit allows a behavior to emit events to its subsribers.
func (c *cell) Emit(event Event) error {
	if aerr := c.act.DoAsync(func() error {
		for _, subscriber := range c.subscribers {
			subscriber.process(event)
		}
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrCellBackend, msgCellBackend, c.behavior.ID())
	}
	return nil
}

// subscribe adds processor engines to the subscribers of this engine.
func (c *cell) subscribe(subscribers []*cell) error {
	if aerr := c.act.DoAsync(func() error {
		for _, subscriber := range subscribers {
			c.subscribers[subscriber.behavior.ID()] = subscriber
		}
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrCellBackend, msgCellBackend, c.behavior.ID())
	}
	return nil
}

// process lets the cell behavior process the event asynchronously.
func (c *cell) process(event Event) error {
	if aerr := c.act.DoAsync(func() error {
		c.behavior.Process(event)
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrCellBackend, msgCellBackend, c.behavior.ID())
	}
	return nil
}

// terminate tells the processor to end and replaces it with a dummy.
func (c *cell) terminate() error {
	var err error
	if aerr := c.act.DoSync(func() error {
		err = c.behavior.Terminate()
		c.behavior = &dummyBehavior{c.behavior.ID()}
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrCellBackend, msgCellBackend, c.behavior.ID())
	}
	return err
}

// stop ends the actor.
func (c *cell) stop(err error) error {
	return c.act.Stop(err)
}

//--------------------
// DUMMY BEHAVIOR
//--------------------

// dummyBehavior will be used by a cell while it's shutting down.
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
