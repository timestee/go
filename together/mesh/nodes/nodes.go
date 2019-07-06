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
	"sync"

	"tideland.dev/go/trace/errors"
)

//--------------------
// MESH
//--------------------

// Mesh operates a set of interacting nodes.
type Mesh struct {
	mu    sync.RWMutex
	nodes map[string]*node
}

// NewMesh creates a new event processing mesh.
func NewMesh() *Mesh {
	m := &Mesh{
		nodes: map[string]*node{},
	}
	return m
}

// SpawnNodes starts nodes to work as parts of the runtime.
func (m *Mesh) SpawnNodes(behaviors ...Behavior) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, behavior := range behaviors {
		_, ok := m.nodes[behavior.ID()]
		if ok {
			continue
		}
		node, err := newNode(behavior)
		if err != nil {
			return err
		}
		m.nodes[behavior.ID()] = node
	}
	return nil
}

// NodeIDs returns the identifiers of the spawned nodes.
func (m *Mesh) NodeIDs() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var ids []string
	for id := range m.nodes {
		ids = append(ids, id)
	}
	return ids
}

// Subscribe subscribes the subscriber processors to the given processor.
func (m *Mesh) Subscribe(nodeID string, subscriberIDs ...string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Retrieve all needed nodes.
	n, ok := m.nodes[nodeID]
	if !ok {
		return errors.New(ErrNodeNotFound, msgNodeNotFound, nodeID)
	}
	var subscribers []*node
	for _, subscriberID := range subscriberIDs {
		subscriber, ok := m.nodes[subscriberID]
		if !ok {
			return errors.New(ErrNodeNotFound, msgNodeNotFound, subscriberID)
		}
		subscribers = append(subscribers, subscriber)
	}
	// Got them, now subscribe.
	return n.subscribe(subscribers)
}

// Emit sends an event to the given node.
func (m *Mesh) Emit(nodeID string, event Event) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Retrieve the needed node.
	n, ok := m.nodes[nodeID]
	if !ok {
		return errors.New(ErrNodeNotFound, msgNodeNotFound, nodeID)
	}
	return n.process(event)
}

// Stop terminates the behaviors, stops the nodes, and cleans up.
func (m *Mesh) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	nerrs := make([]error, len(m.nodes))
	idx := 0
	// Terminate.
	for _, n := range m.nodes {
		nerrs[idx] = n.terminate()
		idx++
	}
	// Stop.
	idx = 0
	for _, n := range m.nodes {
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
