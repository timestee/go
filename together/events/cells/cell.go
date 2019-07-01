// Tideland Go Library - Together - Events - Cells
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package cells // import "tideland.dev/go/together/events/cells"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/together/actor"
	"tideland.dev/go/trace/errors"
)

//--------------------
// BEHAVIOR
//--------------------

// Behavior is the interface that has to be implemented
// for the usage inside of cells.
type Behavior interface {
	// Init is called to initialize the behavior inside the environment.
	// The passed context allows the behavior to interact with this
	// environment and to emit events to subscribers during ProcessEvent().
	// So if this is needed the context should be stored inside the behavior.
	Init(c *Cell) error

	// Terminate is called when a cell is stopped.
	Terminate() error

	// ProcessEvent is called to process the passed event. For emitting
	// events to subscribes during processing the cell passed with the
	// Init() method has to be taken.
	ProcessEvent(event Event)

	// Recover is called in case of an error or panic during the processing
	// of an event. Here the behavior can check if it can recover and establish
	// a valid state. If it's not possible the implementation has to return
	// an error documenting the reason.
	Recover(r interface{}) error
}

//--------------------
// CELL
//--------------------

// Cell runs one behavior for the processing of events and emitting of
// resulting events.
type Cell struct {
	runtime     *Runtime
	id          string
	behavior    Behavior
	subscribers map[string]*Cell
	act         *actor.Actor
}

// NewCell creates a new cell instance with the passed identifier and
// behavior.
func NewCell(id string, behavior Behavior) (*Cell, error) {
	c := &Cell{
		id:          id,
		behavior:    behavior,
		subscribers: map[string]*Cell{},
		act: actor.New(
			actor.WithQueueLen(32),
			actor.WithRecoverer(behavior),
		).Go(),
	}
	err := c.behavior.Init(c)
	if err != nil {
		// Stop the actor with the annotated error.
		return nil, c.act.Stop(errors.Annotate(err, ErrCellInit, msgCellInit, id))
	}
	return c, nil
}

// Runtime returns the runtime environment of the cell.
func (c *Cell) Runtime() *Runtime {
	return c.runtime
}

// ID returns the identifier of the cell.
func (c *Cell) ID() string {
	return c.id
}

// ProcessEvent lets the cell behavior process the event asynchronously
func (c *Cell) ProcessEvent(event Event) error {
	if aerr := c.act.DoSync(func() error {
		c.behavior.ProcessEvent(event)
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrCellBackend, msgCellBackend, c.id)
	}
	return nil
}

// Emit allows a behavior to emit events to the subsribers of the cell.
func (c *Cell) Emit(event Event) error {
	if aerr := c.act.DoAsync(func() error {
		for id, cell := range c.subscribers {
			cell.ProcessEvent(event)
		}
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrCellBackend, msgCellBackend, c.id)
	}
	return nil
}

// Stop terminates the cell.
func (c *Cell) Stop() error {
	var err error
	if aerr := c.act.DoSync(func() error {
		err = c.behavior.Terminate()
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrCellBackend, msgCellBackend, c.id)
	}
	return err
}

// connect integrates the cell into the runtime.
func (c *Cell) connect(runtime *Runtime) error {
	if aerr := c.act.DoSync(func() error {
		c.runtime = runtime
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrCellBackend, msgCellBackend, c.id)
	}
	return nil
}

// subscribe adds cells to the subscribers of the cell.
func (c *Cell) subscribe(cells []*Cell) error {
	if aerr := c.act.DoSync(func() error {
		for _, cell := range cells {
			c.subscribers[c.ID()] = cell
		}
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrCellBackend, msgCellBackend, c.id)
	}
	return nil
}

// EOF
