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
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/event"
)

//--------------------
// TESTS
//--------------------

// TestTopicOnly verifies creation of an event with only a topic.
func TestTopicOnly(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evt := event.New("test")

	assert.Equal(evt.Topic(), "test")
	pl, err := evt.Payload().String("foo")
	assert.ErrorMatch(err, ".*ENOVAL.*")
	assert.Equal(pl, "")
}

// TestKeyValues verifies creation of an event with a topic
// and key/value pairs.
func TestKeyValues(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evt := event.New("topic", "a", 1, "b", "2")

	assert.Equal(evt.Topic(), "topic")
	plai, err := evt.Payload().Int("a")
	assert.NoError(err)
	assert.Equal(plai, 1)
	plbi, err := evt.Payload().Int("b")
	assert.NoError(err)
	assert.Equal(plbi, 2)
	plbs, err := evt.Payload().String("b")
	assert.NoError(err)
	assert.Equal(plbs, "2")
}

// TestWithPayload verifies creation of an event with a topic
// and an external created payload.
func TestWithPayload(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	p := event.NewPayload("a", 1, "b", "2")
	evt := event.WithPayload("topic", p)

	assert.Equal(evt.Topic(), "topic")
	plai, err := evt.Payload().Int("a")
	assert.NoError(err)
	assert.Equal(plai, 1)
	plbi, err := evt.Payload().Int("b")
	assert.NoError(err)
	assert.Equal(plbi, 2)
	plbs, err := evt.Payload().String("b")
	assert.NoError(err)
	assert.Equal(plbs, "2")
}

// TestDefaultValue verifies creation of an event with a topic
// and a valueless final key.
func TestDefaultValue(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	evt := event.New("topic", "a")

	assert.Equal(evt.Topic(), "topic")
	plab, err := evt.Payload().Bool("a")
	assert.NoError(err)
	assert.True(plab)

	evt = event.New("topic", "a", 1, "b", 2, "c")

	assert.Equal(evt.Topic(), "topic")
	plai, err := evt.Payload().Int("a")
	assert.NoError(err)
	assert.Equal(plai, 1)
	plbi, err := evt.Payload().Int("b")
	assert.NoError(err)
	assert.Equal(plbi, 2)
	plcb, err := evt.Payload().Bool("c")
	assert.NoError(err)
	assert.True(plcb)
}

// TestPayloadConv verifies the converting of the string payloads
// into wanted target values.
func TestPayloadConv(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	// First some bool.
	evt := event.New("topic", "a", "true", "b", "false", "c")

	b, err := evt.Payload().Bool("a")
	assert.NoError(err)
	assert.True(b)
	b, err = evt.Payload().Bool("b")
	assert.NoError(err)
	assert.False(b)
	b, err = evt.Payload().Bool("c")
	assert.NoError(err)
	assert.True(b)
	b, err = evt.Payload().Bool("d")
	assert.ErrorMatch(err, ".*ENOVAL.*")

	// Next some float.
	evt = event.New("topic", "a", "1.1", "b", "47.11", "c")

	f, err := evt.Payload().Float64("a")
	assert.NoError(err)
	assert.Equal(f, 1.1)
	f, err = evt.Payload().Float64("b")
	assert.NoError(err)
	assert.Equal(f, 47.11)
	f, err = evt.Payload().Float64("c")
	assert.ErrorMatch(err, `.*parsing "true": invalid syntax.*`)
	f, err = evt.Payload().Float64("d")
	assert.ErrorMatch(err, ".*ENOVAL.*")

	// Now some int.
	evt = event.New("topic", "a", 1, "b", "-4711", "c")

	i, err := evt.Payload().Int("a")
	assert.NoError(err)
	assert.Equal(i, 1)
	i, err = evt.Payload().Int("b")
	assert.NoError(err)
	assert.Equal(i, -4711)
	i, err = evt.Payload().Int("c")
	assert.ErrorMatch(err, `.*parsing "true": invalid syntax.*`)
	i, err = evt.Payload().Int("d")
	assert.ErrorMatch(err, ".*ENOVAL.*")

	// Next some int64.
	evt = event.New("topic", "a", "1", "b", "-4711", "c")

	ii, err := evt.Payload().Int64("a")
	assert.NoError(err)
	assert.Equal(ii, int64(1))
	ii, err = evt.Payload().Int64("b")
	assert.NoError(err)
	assert.Equal(ii, int64(-4711))
	ii, err = evt.Payload().Int64("c")
	assert.ErrorMatch(err, `.*parsing "true": invalid syntax.*`)
	ii, err = evt.Payload().Int64("d")
	assert.ErrorMatch(err, ".*ENOVAL.*")

	// Now some uint64.
	evt = event.New("topic", "a", "1", "b", "-4711", "c")

	ui, err := evt.Payload().Uint64("a")
	assert.NoError(err)
	assert.Equal(ui, uint64(1))
	ui, err = evt.Payload().Uint64("b")
	assert.ErrorMatch(err, `.*parsing "-4711": invalid syntax.*`)
	ui, err = evt.Payload().Uint64("c")
	assert.ErrorMatch(err, `.*parsing "true": invalid syntax.*`)
	ui, err = evt.Payload().Uint64("d")
	assert.ErrorMatch(err, ".*ENOVAL.*")

	// Next some time.
	now := time.Now()
	then := time.Date(2019, time.July, 23, 20, 0, 0, 0, time.UTC)
	evt = event.New("topic", "a", now, "b", then.Format(time.RFC3339))

	ta, err := evt.Payload().Time("a")
	assert.NoError(err)
	assert.Equal(ta, now)
	tb, err := evt.Payload().Time("b")
	assert.NoError(err)
	assert.Equal(tb, then)

	// And some duration.
	fiveSecs := 5 * time.Second
	evt = event.New("topic", "a", fiveSecs, "b", "5s")

	da, err := evt.Payload().Duration("a")
	assert.NoError(err)
	assert.Equal(da, fiveSecs)
	db, err := evt.Payload().Duration("b")
	assert.NoError(err)
	assert.Equal(db, fiveSecs)
}

// TestPayloadClone verifies the converting of the string payloads
// into wanted target values.
func TestPayloadClone(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	pA := event.NewPayload("a", 1, "b", "two", "c", 3.0)

	ia, err := pA.Int("a")
	assert.NoError(err)
	assert.Equal(ia, 1)
	sb, err := pA.String("b")
	assert.NoError(err)
	assert.Equal(sb, "two")

	pB := pA.Clone("a", "4711", "d", "foo")

	ia, err = pB.Int("a")
	assert.NoError(err)
	assert.Equal(ia, 4711)
	assert.Length(pB.Keys(), 4)
}

// EOF
