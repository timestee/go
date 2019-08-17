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
	"tideland.dev/go/trace/failure"
)

//--------------------
// MESH
//--------------------

// Mesh operates a set of interacting cells.
type Mesh struct {
	mu    sync.RWMutex
	cells cellRegistry
}

// New creates a new event processing mesh.
func New() *Mesh {
	m := &Mesh{
		cells: cellRegistry{},
	}
	return m
}

// SpawnCells starts cells running the passed behaviors to work as parts
// of the mesh.
func (m *Mesh) SpawnCells(behaviors ...Behavior) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, behavior := range behaviors {
		id := behavior.ID()
		if m.cells.contains(id) {
			// No double deployment.
			continue
		}
		cell, err := newCell(m, behavior)
		if err != nil {
			return err
		}
		m.cells.add(id, cell)
	}
	return nil
}

// StopCells terminates the given cells.
func (m *Mesh) StopCells(ids ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, id := range ids {
		if err := m.cells.unsubscribeFromAll(id); err != nil {
			return err
		}
		return m.cells.remove(id)
	}
	return nil
}

// Cells returns the identifiers of the spawned cells.
func (m *Mesh) Cells() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var ids []string
	for id := range m.cells {
		ids = append(ids, id)
	}
	return ids
}

// Subscribers retrieves the subscriber IDs of a cell.
func (m *Mesh) Subscribers(id string) ([]string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Retrieve all needed cells.
	entry, ok := m.cells[id]
	if !ok {
		return nil, failure.New("cannot find cell %q", id)
	}
	return entry.cell.subscribers()
}

// Subscribe connects cells to the given cell.
func (m *Mesh) Subscribe(id string, subscriberIDs ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.cells.subscribe(id, subscriberIDs)
}

// Unsubsribe disconnect cells from the given cell.
func (m *Mesh) Unsubscribe(id string, unsubscriberIDs ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.cells.unsubscribe(id, unsubscriberIDs)
}

// Emit sends an event to the given cell.
func (m *Mesh) Emit(id string, evt *event.Event) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Retrieve the needed cell.
	entry, ok := m.cells[id]
	if !ok {
		return failure.New("cannot find cell %q", id)
	}
	return entry.cell.process(evt)
}

// Broadcast sends an event to all cells.
func (m *Mesh) Broadcast(evt *event.Event) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cerrs := make([]error, len(m.cells))
	idx := 0
	// Broadcast.
	for _, entry := range m.cells {
		cerrs[idx] = entry.cell.process(evt)
		idx++
	}
	// Return collected errors.
	return failure.Collect(cerrs...)
}

// Stop terminates the cells and cleans up.
func (m *Mesh) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	cerrs := make([]error, len(m.cells))
	idx := 0
	// Terminate.
	for _, entry := range m.cells {
		cerrs[idx] = entry.cell.stop()
		idx++
	}
	m.cells = cellRegistry{}
	// Return collected errors.
	return failure.Collect(cerrs...)
}

// EOF
