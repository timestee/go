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
// COUNTDOWN BEHAVIOR
//--------------------

// Zeroer is called when the countdown reaches zero. The collected
// events are passed, the returned event will be emitted, and the
// returned number sets a new countdown.
type Zeroer func(accessor event.SinkAccessor) (*event.Event, int, error)

// countdownBehavior counts events based on the counter function.
type countdownBehavior struct {
	id      string
	emitter mesh.Emitter
	sink    *event.Sink
	t       int
	zeroer  Zeroer
}

// NewCountdownBehavior creates a countdown behavior based on the given
// t value and zeroer function.
func NewCountdownBehavior(id string, t int, zeroer Zeroer) mesh.Behavior {
	return &countdownBehavior{
		id:     id,
		sink:   event.NewSink(t),
		t:      t,
		zeroer: zeroer,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *countdownBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *countdownBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *countdownBehavior) Terminate() error {
	b.sink.Clear()
	return nil
}

// Process puts the received events into a sink. When reaching t the
// zeroer will be called with access to the sink. Its returned event
// will be emitted, the returned t will be set, and the sink cleared.
func (b *countdownBehavior) Process(evt *event.Event) error {
	switch evt.Topic() {
	case event.TopicReset:
		t := evt.Payload().At("t").AsInt(b.t)
		b.sink.Clear()
		b.t = t
	default:
		if b.t <= 0 {
			return nil
		}
		sl, err := b.sink.Push(evt)
		if err != nil {
			return err
		}
		if sl == b.t {
			// T-0, call the zeroer, set t, and emit event.
			zevt, t, err := b.zeroer(b.sink)
			if err != nil {
				return err
			}
			b.sink.Clear()
			b.t = t
			return b.emitter.Broadcast(zevt)
		}
	}
	return nil
}

// Recover from an error.
func (b *countdownBehavior) Recover(err interface{}) error {
	return b.sink.Clear()
}

// EOF
