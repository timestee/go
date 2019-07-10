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

// TestPayloadConv verifies the converting of the string payloads
// into wanted target values.
func TestPayloadConv(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	// First some bool.
	evt := event.New("topic", "a", "true", "b", "false", "c")

	b, err := evt.PayloadBool("a")
	assert.NoError(err)
	assert.True(b)
	b, err = evt.PayloadBool("b")
	assert.NoError(err)
	assert.False(b)
	b, err = evt.PayloadBool("c")
	assert.NoError(err)
	assert.True(b)
	b, err = evt.PayloadBool("d")
	assert.NoError(err)
	assert.False(b)

	// Next some float.
	evt = event.New("topic", "a", "1.1", "b", "47.11", "c")

	f, err := evt.PayloadFloat("a")
	assert.NoError(err)
	assert.Equal(f, 1.1)
	f, err = evt.PayloadFloat("b")
	assert.NoError(err)
	assert.Equal(f, 47.11)
	f, err = evt.PayloadFloat("c")
	assert.ErrorMatch(err, `.*parsing "true": invalid syntax.*`)
	f, err = evt.PayloadFloat("d")
	assert.ErrorMatch(err, `.*parsing "false": invalid syntax.*`)

	// Next some int.
	evt = event.New("topic", "a", "1", "b", "-4711", "c")

	i, err := evt.PayloadInt("a")
	assert.NoError(err)
	assert.Equal(i, int64(1))
	i, err = evt.PayloadInt("b")
	assert.NoError(err)
	assert.Equal(i, int64(-4711))
	i, err = evt.PayloadInt("c")
	assert.ErrorMatch(err, `.*parsing "true": invalid syntax.*`)
	i, err = evt.PayloadInt("d")
	assert.ErrorMatch(err, `.*parsing "false": invalid syntax.*`)

	// And finally some uint.
	evt = event.New("topic", "a", "1", "b", "-4711", "c")

	ui, err := evt.PayloadUint("a")
	assert.NoError(err)
	assert.Equal(ui, uint64(1))
	ui, err = evt.PayloadUint("b")
	assert.ErrorMatch(err, `.*parsing "-4711": invalid syntax.*`)
	ui, err = evt.PayloadUint("c")
	assert.ErrorMatch(err, `.*parsing "true": invalid syntax.*`)
	ui, err = evt.PayloadUint("d")
	assert.ErrorMatch(err, `.*parsing "false": invalid syntax.*`)
}

// EOF
