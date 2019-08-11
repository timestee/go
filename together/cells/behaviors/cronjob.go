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
// CRONJOB BEHAVIOR
//--------------------

// Cronjob dynamically returns the event to be emitted by the cronjob behavior.
type Cronjob func() *event.Event

// cronjobBehavior chronologically emits events.
type cronjobBehavior struct {
	id       string
	emitter  mesh.Emitter
	duration time.Duration
	cronjob  Cronjob
	loop     *loop.Loop
}

// NewCronjobBehavior creates a ticker behavior for the emitting of
// "tick" events every given duration.
func NewCronjobBehavior(id string, duration time.Duration, cronjob Cronjob) mesh.Behavior {
	return &cronjobBehavior{
		id:       id,
		duration: duration,
		cronjob:  cronjob,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *cronjobBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *cronjobBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	b.loop = loop.New(b.tickerLoop).Go()
	return nil
}

// Terminate the behavior.
func (b *cronjobBehavior) Terminate() error {
	return b.loop.Stop(nil)
}

// Process emits a ticker event each time the defined duration elapsed.
func (b *cronjobBehavior) Process(evt *event.Event) error {
	if evt.Topic() == TopicTick {
		b.emitter.Broadcast(b.cronjob())
	}
	return nil
}

// Recover from an error. Counter will be set back to the initial counter.
func (b *cronjobBehavior) Recover(err interface{}) error {
	return nil
}

// tickerLoop is the sending a tick event to itself. It acts there to
// avoid races when subscribers are updated.
func (b *cronjobBehavior) tickerLoop(c *notifier.Closer) error {
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
