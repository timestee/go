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
	"sync"

	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/trace/errors"
)

//--------------------
// MESH
//--------------------

// Mesh operates a set of interacting cells.
type Mesh struct {
	mu    sync.RWMutex
	cells map[string]*cell
}

// New creates a new event processing mesh.
func New() *Mesh {
	m := &Mesh{
		cells: map[string]*cell{},
	}
	return m
}

// SpawnCells starts cells to work as parts of the runtime.
func (m *Mesh) SpawnCells(behaviors ...Behavior) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, behavior := range behaviors {
		_, ok := m.cells[behavior.ID()]
		if ok {
			continue
		}
		c, err := newCell(behavior)
		if err != nil {
			return err
		}
		m.cells[behavior.ID()] = c
	}
	return nil
}

// CellIDs returns the identifiers of the spawned cells.
func (m *Mesh) CellIDs() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var ids []string
	for id := range m.cells {
		ids = append(ids, id)
	}
	return ids
}

// Subscribe subscribes the subscriber processors to the given processor.
func (m *Mesh) Subscribe(cellID string, subscriberIDs ...string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Retrieve all needed cells.
	c, ok := m.cells[cellID]
	if !ok {
		return errors.New(ErrCellNotFound, msgCellNotFound, cellID)
	}
	var subscribers []*cell
	for _, subscriberID := range subscriberIDs {
		subscriber, ok := m.cells[subscriberID]
		if !ok {
			return errors.New(ErrCellNotFound, msgCellNotFound, subscriberID)
		}
		subscribers = append(subscribers, subscriber)
	}
	// Got them, now subscribe.
	return c.subscribe(subscribers)
}

// Emit sends an event to the given cell.
func (m *Mesh) Emit(cellID string, evt *event.Event) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Retrieve the needed cell.
	c, ok := m.cells[cellID]
	if !ok {
		return errors.New(ErrCellNotFound, msgCellNotFound, cellID)
	}
	return c.process(evt)
}

// Stop terminates the behaviors, stops the cells, and cleans up.
func (m *Mesh) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	nerrs := make([]error, len(m.cells))
	idx := 0
	// Terminate.
	for _, n := range m.cells {
		nerrs[idx] = n.terminate()
		idx++
	}
	// Stop.
	idx = 0
	for _, n := range m.cells {
		nerrs[idx] = n.stop(nerrs[idx])
	}
	// Drop nil errors.
	var errs []error
	for _, err := range nerrs {
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Collect(errs...)
	}
	return nil
}

// EOF
