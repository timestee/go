// Tideland Go Library - Together - Cells - Behaviors
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// COLLECTOR BEHAVIOR
//--------------------

// collectorBehavior collects events for debugging.
type collectorBehavior struct {
	id      string
	emitter mesh.Emitter
	max     int
	sink    *event.Sink
	process event.SinkProcessor
}

// NewCollectorBehavior creates a collector behavior. It collects
// a maximum number of events, each event is passed through. If the
// maximum number is 0 it collects until the topic "reset". After
// receiving the topic "process" the processor will be called and
// the collected events will be reset afterwards.
func NewCollectorBehavior(id string, max int, processor event.SinkProcessor) mesh.Behavior {
	return &collectorBehavior{
		id:      id,
		max:     max,
		process: processor,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *collectorBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *collectorBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	b.sink = event.NewSink(b.max)
	return nil
}

// Terminate the behavior.
func (b *collectorBehavior) Terminate() error {
	b.sink.Clear()
	return nil
}

// Process collects, processes, and re-emits events.
func (b *collectorBehavior) Process(evt *event.Event) error {
	switch evt.Topic() {
	case event.TopicProcess:
		pl, err := b.process(b.sink)
		if err != nil {
			return err
		}
		return b.emitter.Broadcast(event.New(event.TopicResult, pl))
	case event.TopicReset:
		return b.sink.Clear()
	default:
		b.sink.Push(evt)
		return b.emitter.Broadcast(evt)
	}
}

// Recover from an error.
func (b *collectorBehavior) Recover(err interface{}) error {
	return b.sink.Clear()
}

// EOF
