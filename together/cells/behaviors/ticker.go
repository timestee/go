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
	"time"

	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
	"tideland.dev/go/together/loop"
	"tideland.dev/go/together/notifier"
)

//--------------------
// TICKER BEHAVIOR
//--------------------

// tickerBehavior chronologically emits events.
type tickerBehavior struct {
	id       string
	emitter  mesh.Emitter
	duration time.Duration
	loop     *loop.Loop
}

// NewTickerBehavior creates a ticker behavior for the emitting of
// "tick" events every given duration.
func NewTickerBehavior(id string, duration time.Duration) mesh.Behavior {
	return &tickerBehavior{
		id:       id,
		duration: duration,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *tickerBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *tickerBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	b.loop = loop.New(b.tickerLoop).Go()
	return nil
}

// Terminate the behavior.
func (b *tickerBehavior) Terminate() error {
	return b.loop.Stop(nil)
}

// Process emits a ticker event each time the defined duration elapsed.
func (b *tickerBehavior) Process(evt *event.Event) error {
	if evt.Topic() == TopicTick {
		b.emitter.Broadcast(event.New(TopicTick, "id", b.id))
	}
	return nil
}

// Recover from an error. Counter will be set back to the initial counter.
func (b *tickerBehavior) Recover(err interface{}) error {
	return nil
}

// tickerLoop is the sending a tick event to itself. It acts there to
// avoid races when subscribers are updated.
func (b *tickerBehavior) tickerLoop(c *notifier.Closer) error {
	ticker := time.NewTicker(b.duration)
	defer ticker.Stop()
	for {
		select {
		case <-c.Done():
			return nil
		case <-ticker.C:
			b.emitter.Self(event.New(TopicTick))
		}
	}
}

// EOF
