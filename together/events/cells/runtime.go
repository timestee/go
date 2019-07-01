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
	"sync"

	"tideland.dev/go/trace/errors"
)

//--------------------
// RUNTIME
//--------------------

// Runtime operates a set of interacting cells.
type Runtime struct {
	mu    sync.RWMutex
	cells map[string]*Cell
}

// NewRuntime creates a new cell runtime.
func NewRuntime() *Runtime {
	r := &Runtime{
		cells: map[string]*Cell{},
	}
	return r
}

// AddCells adds new cells to the runtime. If a cell is added
// multiple times it's no problem. But different cells with the
// same ID lead to an error.
func (r *Runtime) AddCells(cells ...*Cell) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	// Check cells and identifiers.
	for _, cell := range cells {
		rcell, ok := r.cells[cell.ID()]
		if ok && rcell != cell {
			return errors.New(ErrRuntimeAdd, msgRuntimeAdd, cell.ID())
		}
	}
	// All fine, add them.
	for _, cell := range cells {
		r.cells[cell.ID()] = cell
	}
	return nil
}

// EOF
