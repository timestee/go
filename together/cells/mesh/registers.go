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
	"tideland.dev/go/trace/errors"
)

//--------------------
// REGISTRIES
//--------------------

// subscribedRegistry manages the cell IDs a cell subscribed to.
type subscribedRegistry map[string]struct{}

// add a cell ID to the registry.
func (sr subscribedRegistry) add(id string) {
	sr[id] = struct{}{}
}

// remove a cell ID from the registry.
func (sr subscribedRegistry) remove(id string) {
	delete(sr, id)
}

// cellEntry containes cells and the IDs of the cells it subscribed to.
type cellEntry struct {
	id           string
	cell         *cell
	subscribedTo subscribedRegistry
}

// cellRegistry manages a number of cells and provides some convenience.
type cellRegistry map[string]*cellEntry

// contains checks if the registry contains the given cell.
func (cr cellRegistry) contains(id string) bool {
	_, ok := cr[id]
	return ok
}

// add registers a cell.
func (cr cellRegistry) add(id string, c *cell) {
	cr[id] = &cellEntry{
		cell:         c,
		subscribedTo: subscribedRegistry{},
	}
}

// removes deregisters a cell.
func (cr cellRegistry) remove(id string) {
	delete(cr, id)
}

// cells retirieves the identified cells.
func (cr cellRegistry) cells(ids []string) ([]*cell, error) {
	var cells []*cell
	for _, id := range ids {
		entry, ok := cr[id]
		if !ok {
			return nil, errors.New(ErrCellNotFound, msgCellNotFound, id)
		}
		cells = append(cells, entry.cell)
	}
	return cells, nil
}

// subscribe subscribes the wanted subscribers to a cell and also tells
// those where they subscribed to.
func (cr cellRegistry) subscribe(id string, subscriberIDs []string) error {
	// Retrieve cells.
	entry, ok := cr[id]
	if !ok {
		return errors.New(ErrCellNotFound, msgCellNotFound, id)
	}
	var subscribers []*cell
	for _, subscriberID := range subscriberIDs {
		subscriber, ok := cr[subscriberID]
		if !ok {
			return errors.New(ErrCellNotFound, msgCellNotFound, subscriberID)
		}
		subscribers = append(subscribers, subscriber.cell)
	}
	// Got all, now subscribe.
	if err := entry.cell.subscribe(subscribers); err != nil {
		return err
	}
	// Tell subscribers where they subscribed to.
	for _, subscriberID := range subscriberIDs {
		cr[subscriberID].subscribedTo.add(id)
	}
	return nil
}

// unsubscribe unsubscribes the wanted subscribers from a cell and also tells
// those that they should drop their knowledge about the subscriber.
func (cr cellRegistry) unsubscribe(id string, unsubscriberIDs []string) error {
	// Retrieve cells.
	entry, ok := cr[id]
	if !ok {
		return errors.New(ErrCellNotFound, msgCellNotFound, id)
	}
	var unsubscribers []*cell
	for _, unsubscriberID := range unsubscriberIDs {
		unsubscriber, ok := cr[unsubscriberID]
		if !ok {
			return errors.New(ErrCellNotFound, msgCellNotFound, unsubscriberID)
		}
		unsubscribers = append(unsubscribers, unsubscriber.cell)
	}
	// Got all, now unsubscribe.
	if err := entry.cell.unsubscribe(unsubscribers); err != nil {
		return err
	}
	// Tell unsubscribers that they aren't subscribed anymore.
	for _, unsubscriberID := range unsubscriberIDs {
		cr[unsubscriberID].subscribedTo.remove(id)
	}
	return nil
}

// EOF
