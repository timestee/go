// Tideland Go Library - Together - Events - Runtime
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package runtime // import "tideland.dev/go/together/events/runtime"

//--------------------
// IMPORTS
//--------------------

import (
	"sync"

	"tideland.dev/go/trace/errors"
)

//--------------------
// EVENT
//--------------------

// Event is some data with a topic describing it.
type Event struct {
	Topic string
	Data  string
}

//--------------------
// EMITTER
//--------------------

// Emitter describes a type to emit events to subscribers. An instance
// is passed to a Processor at initialization.
type Emitter interface {
	Emit(event Event) error
}

//--------------------
// PROCESSOR
//--------------------

// Processor is the interface that has to be implemented for event
// processing inside the runtime.
type Processor interface {
	// ID returns the individual identifier of a processor instance.
	// Processors can be deployed multiple times as long as these return
	// different identifiers.
	ID() string

	// Init is called by the runtime to initialize the processor.
	// Events can be sent to subscribers by emitter.Emit().
	Init(emitter Emitter) error

	// Terminate is called when a processor is stopped.
	Terminate() error

	// Process is called to process the given event.
	Process(event Event)

	// Recover is called in case of an error or panic during the processing
	// of an event. Here the behavior can check if it can recover and establish
	// a valid state. If it's not possible the implementation has to return
	// an error documenting the reason.
	Recover(r interface{}) error
}

//--------------------
// RUNTIME
//--------------------

// Runtime operates a set of interacting cells.
type Runtime struct {
	mu         sync.RWMutex
	processors map[string]*processorEngine
}

// New creates a new event processing runtime.
func New() *Runtime {
	r := &Runtime{
		processors: map[string]*processorEngine{},
	}
	return r
}

// SpawnProcessors starts processors to work as parts of the runtime.
func (r *Runtime) SpawnProcessors(processors ...Processor) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, processor := range processors {
		_, ok := r.processors[processor.ID()]
		if ok {
			continue
		}
		engine, err := newProcessorEngine(processor)
		if err != nil {
			return err
		}
		r.processors[processor.ID()] = engine
	}
	return nil
}

// ProcessorIDs returns the identifiers of the spawned processors.
func (r *Runtime) ProcessorIDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var ids []string
	for id := range r.processors {
		ids = append(ids, id)
	}
	return ids
}

// Subscribe subscribes the subscriber processors to the given processor.
func (r *Runtime) Subscribe(processorID string, subscriberIDs ...string) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Retrieve all needed processor engines.
	processor, ok := r.processors[processorID]
	if !ok {
		return errors.New(ErrProcessorNotFound, msgProcessorNotFound, processorID)
	}
	var subscribers []*processorEngine
	for _, subscriberID := range subscriberIDs {
		subscriber, ok := r.processors[subscriberID]
		if !ok {
			return errors.New(ErrProcessorNotFound, msgProcessorNotFound, subscriberID)
		}
		subscribers = append(subscribers, subscriber)
	}
	// Got them, now subscribe.
	return processor.subscribe(subscribers)
}

// Emit sends an event to the given processor.
func (r *Runtime) Emit(processorID string, event Event) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	// Retrieve the needed processor engine.
	processor, ok := r.processors[processorID]
	if !ok {
		return errors.New(ErrProcessorNotFound, msgProcessorNotFound, processorID)
	}
	return processor.process(event)
}

// Stop terminates the processors, stops the engines, and cleans up.
func (r *Runtime) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	perrs := make([]error, len(r.processors))
	idx := 0
	// Terminate.
	for _, processor := range r.processors {
		perrs[idx] = processor.terminate()
		idx++
	}
	// Stop.
	idx = 0
	for _, processor := range r.processors {
		perrs[idx] = processor.stop(perrs[idx])
	}
	// Drop nil errors.
	var errs []error
	for _, err := range perrs {
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
