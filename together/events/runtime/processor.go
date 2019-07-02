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
	"tideland.dev/go/together/actor"
	"tideland.dev/go/trace/errors"
)

//--------------------
// PROCESSOR ENGINE
//--------------------

// cell runs one behavior for the processing of events and emitting of
// resulting events.
type processorEngine struct {
	processor   Processor
	subscribers map[string]*processorEngine
	act         *actor.Actor
}

// newProcessorEngine creates a new engine running the given processor
// in a goroutine.
func newProcessorEngine(processor Processor) (*processorEngine, error) {
	pe := &processorEngine{
		processor:   processor,
		subscribers: map[string]*processorEngine{},
		act: actor.New(
			actor.WithQueueLen(32),
			actor.WithRecoverer(processor.Recover),
		).Go(),
	}
	err := pe.processor.Init(pe)
	if err != nil {
		// Stop the actor with the annotated error.
		return nil, pe.act.Stop(errors.Annotate(err, ErrEngineInit, msgEngineInit, processor.ID()))
	}
	return pe, nil
}

// Emit allows a processor to emit events to its subsribers.
func (pe *processorEngine) Emit(event Event) error {
	if aerr := pe.act.DoAsync(func() error {
		for _, engine := range pe.subscribers {
			engine.process(event)
		}
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrEngineBackend, msgEngineBackend, pe.processor.ID())
	}
	return nil
}

// subscribe adds processor engines to the subscribers of this engine.
func (pe *processorEngine) subscribe(engines []*processorEngine) error {
	if aerr := pe.act.DoSync(func() error {
		for _, engine := range engines {
			pe.subscribers[engine.processor.ID()] = engine
		}
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrEngineBackend, msgEngineBackend, pe.processor.ID())
	}
	return nil
}

// process lets the processor engine process the event asynchronously.
func (pe *processorEngine) process(event Event) error {
	if aerr := pe.act.DoAsync(func() error {
		pe.processor.Process(event)
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrEngineBackend, msgEngineBackend, pe.processor.ID())
	}
	return nil
}

// stop terminates the processor and the engine.
func (pe *processorEngine) stop() error {
	var err error
	if aerr := pe.act.DoSync(func() error {
		err = pe.processor.Terminate()
		return nil
	}); aerr != nil {
		return errors.Annotate(aerr, ErrEngineBackend, msgEngineBackend, pe.processor.ID())
	}
	return err
}

// EOF
