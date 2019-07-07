// Tideland Go Library - Together - Cells - Event
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package event // import "tideland.dev/go/together/cells/event"

//--------------------
// DEFAULTS
//--------------------

const (
	DefaultTopic     = "signal"
	DefaultValue     = "true"
	NonExistingValue = "false"
)

//--------------------
// EVENT
//--------------------

type Event struct {
	topic   string
	payload map[string]string
}

// New creates a new event. First data is the topic, the rest are key/value
// pairs for the payload. In case of no data the topic is "signal", in case
// of a final key without a value it's value will be set to "true".
func New(datas ...string) *Event {
	e := &Event{
		topic:   DefaultTopic,
		payload: map[string]string{},
	}
	var key string
	for i, data := range datas {
		switch {
		case i == 0:
			e.topic = data
			continue
		case i%2 == 1:
			key = data
			e.payload[key] = DefaultValue
			continue
		default:
			e.payload[key] = data
		}
	}
	return e
}

// Topic returns the event topic.
func (e *Event) Topic() string {
	return e.topic
}

// Payload returns the payload for the given key. If that doesn't
// exist the value "false" will be returned.
func (e *Event) Payload(key string) string {
	p, ok := e.payload[key]
	if !ok {
		return NonExistingValue
	}
	return p
}

// EOF
