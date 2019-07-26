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
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/trace/errors"
)

//--------------------
// CELL
//--------------------

// cell runs a behavior for the processing of events and emitting of
// resulting events.
type cell struct {
	behavior        Behavior
	subscribedCells map[string]*cell
	act             *actor.Actor
}

// newCell creates a new cell running the given behavior in a goroutine.
func newCell(behavior Behavior) (*cell, error) {
	c := &cell{
		behavior:        behavior,
		subscribedCells: map[string]*cell{},
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

// ID allows a behavior sub-component using an emitter to retrieve
// the identifier of the cell. Can only be called inside event processing
func (c *cell) ID() string {
	return c.behavior.ID()
}

// Emit allows a behavior to emit events to its subsribers.
func (c *cell) Emit(evt *event.Event) error {
	var serrs []error
	for _, subscriber := range c.subscribedCells {
		serrs = append(serrs, subscriber.process(evt))
	}
	return errors.Collect(serrs...)
}

// subscribers returns the subscriber IDs of the cell.
func (c *cell) subscribers() ([]string, error) {
	var subscriberIDs []string
	if aerr := c.act.DoSync(func() error {
		for subscriberID := range c.subscribedCells {
			subscriberIDs = append(subscriberIDs, subscriberID)
		}
		return nil
	}); aerr != nil {
		return nil, errors.Annotate(aerr, ErrCellBackend, msgCellBackend, c.behavior.ID())
	}
	return subscriberIDs, nil
}

// subscribe adds cells to the subscribers of this cell.
func (c *cell) subscribe(subscribers []*cell) error {
	if aerr := c.act.DoAsync(func() error {
		for _, subscriber := range subscribers {
			c.subscribedCells[subscriber.behavior.ID()] = subscriber
		}
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrCellBackend, msgCellBackend, c.behavior.ID())
	}
	return nil
}

// unsubscribe removes cells from the subscribers of this cell.
func (c *cell) unsubscribe(subscribers []*cell) error {
	if aerr := c.act.DoAsync(func() error {
		for _, subscriber := range subscribers {
			delete(c.subscribedCells, subscriber.behavior.ID())
		}
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrCellBackend, msgCellBackend, c.behavior.ID())
	}
	return nil
}

// process lets the cell behavior process the event asynchronously.
func (c *cell) process(evt *event.Event) error {
	var berr error
	if aerr := c.act.DoAsync(func() error {
		berr = c.behavior.Process(evt)
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrCellBackend, msgCellBackend, c.behavior.ID())
	}
	return berr
}

// stop terminates the cell and stops the actor.
func (c *cell) stop() error {
	var cerr error
	if aerr := c.act.DoSync(func() error {
		cerr = c.behavior.Terminate()
		c.behavior = &dummyBehavior{c.behavior.ID()}
		c.subscribedCells = map[string]*cell{}
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrCellBackend, msgCellBackend, c.behavior.ID())
	}
	// Stop actor with cell or given error.
	return c.act.Stop(cerr)
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

func (db *dummyBehavior) Process(evt *event.Event) error {
	return nil
}

func (db *dummyBehavior) Recover(r interface{}) error {
	return nil
}

// EOF
