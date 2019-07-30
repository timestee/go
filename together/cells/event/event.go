// Tideland Go Library - Together - Cells - Event
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package event // import "tideland.dev/go/together/cells/event"

//--------------------
// IMPORTS
//--------------------

import (
	"time"
)

//--------------------
// CONSTANTS
//--------------------

// Standard topics.
const (
	TopicCollected = "collected"
	TopicCounted   = "counted"
	TopicProcess   = "process"
	TopicProcessed = "processed"
	TopicReset     = "reset"
	TopicResult    = "result"
	TopicStatus    = "status"
	TopicTick      = "tick"
)

//--------------------
// EVENT
//--------------------

// Event describes an event of the cells. It contains a topic as well as
// a possible number of key/value pairs as payload.
type Event struct {
	timestamp time.Time
	topic     string
	payload   *Payload
}

// New creates a new event. The arguments after the topic are taken
// to create a new payload.
func New(topic string, kvs ...interface{}) *Event {
	return &Event{
		timestamp: time.Now(),
		topic:     topic,
		payload:   NewPayload(kvs...),
	}
}

// WithPayload creates a new event with a given external payload.
func WithPayload(topic string, pl *Payload) *Event {
	return &Event{
		timestamp: time.Now(),
		topic:     topic,
		payload:   pl,
	}
}

// Timestamp returns the event timestamp.
func (e *Event) Timestamp() time.Time {
	return e.timestamp
}

// Topic returns the event topic.
func (e *Event) Topic() string {
	return e.topic
}

// Payload returns the event payload.
func (e *Event) Payload() *Payload {
	return e.payload
}

// EOF
