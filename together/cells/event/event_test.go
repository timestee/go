// Tideland Go Library - Together - Cells - Event - Unit Tests
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package event_test // import "tideland.dev/go/together/cells/event"

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/event"
)

//--------------------
// TESTS
//--------------------

// TestEmptyEvent verifies creation of an empty event with default
// topic and return of default value for random key.
func TestEmptyEvent(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evt := event.New()

	assert.Equal(evt.Topic(), event.DefaultTopic)
	assert.Equal(evt.Payload("a"), event.NonExistingValue)
	assert.Equal(evt.Payload("b"), event.NonExistingValue)
	assert.Equal(evt.Payload("c"), event.NonExistingValue)
}

// TestTopicOnly verifies creation of an event with only a topic.
func TestTopicOnly(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evt := event.New("topic")

	assert.Equal(evt.Topic(), "topic")
	assert.Equal(evt.Payload("a"), event.NonExistingValue)
	assert.Equal(evt.Payload("b"), event.NonExistingValue)
	assert.Equal(evt.Payload("c"), event.NonExistingValue)
}

// TestKeyValues verifies creation of an event with a topic
// and matching key/value pairs.
func TestKeyValues(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evt := event.New("topic", "a", "1", "b", "2")

	assert.Equal(evt.Topic(), "topic")
	assert.Equal(evt.Payload("a"), "1")
	assert.Equal(evt.Payload("b"), "2")
	assert.Equal(evt.Payload("c"), event.NonExistingValue)
}

// TestDefaultValue verifies creation of an event with a topic
// and a valueless final key.
func TestDefaultValue(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evt := event.New("topic", "a")

	assert.Equal(evt.Topic(), "topic")
	assert.Equal(evt.Payload("a"), event.DefaultValue)
	assert.Equal(evt.Payload("b"), event.NonExistingValue)

	evt = event.New("topic", "a", "1", "b", "2", "c")

	assert.Equal(evt.Topic(), "topic")
	assert.Equal(evt.Payload("a"), "1")
	assert.Equal(evt.Payload("b"), "2")
	assert.Equal(evt.Payload("c"), event.DefaultValue)
	assert.Equal(evt.Payload("d"), event.NonExistingValue)
}

// EOF
