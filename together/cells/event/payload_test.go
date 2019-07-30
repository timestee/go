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

// TestSimplePayload verifies creation of a payload with key/value pairs.
func TestSimplePayload(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	now := time.Now()
	pl := event.NewPayload("a", 1, "b", "2", "c", now, "d")

	va := pl.At("a").AsInt(0)
	assert.Equal(va, 1)
	vb1 := pl.At("b").AsInt(0)
	assert.Equal(vb1, 2)
	vb2 := pl.At("b").AsString("0")
	assert.Equal(vb2, "2")
	vc := pl.At("c").AsTime(time.Time{})
	assert.Equal(vc, now)
	vd := pl.At("d").AsBool(false)
	assert.True(vd)
	assert.True(pl.At("d").IsDefined())
	assert.True(pl.At("e").IsUndefined())
}

// TestPayloadDefaults verifies retrieving default values from payloads.
func TestPayloadDefaults(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	now := time.Now()
	pl := event.NewPayload()

	vs := pl.At("a").AsString("foo")
	assert.Equal(vs, "foo")
	vi := pl.At("a").AsInt(1234)
	assert.Equal(vi, 1234)
	vf := pl.At("a").AsFloat64(12.34)
	assert.Equal(vf, 12.34)
	vb := pl.At("a").AsBool(true)
	assert.Equal(vb, true)
	vt := pl.At("a").AsTime(now)
	assert.Equal(vt, now)
	vd := pl.At("a").AsDuration(time.Second)
	assert.Equal(vd, time.Second)
}

// TestNestedPayloads verifies retrieving values from nested payloads.
func TestNestedPayloads(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	plc := event.NewPayload("ca", 100, "cb", 200)
	plb := event.NewPayload("ba", 10, "bb", plc)
	pla := event.NewPayload("aa", 1, "ab", plb)

	vaa := pla.At("aa").AsInt(0)
	assert.Equal(vaa, 1)
	vba := pla.At("ab/ba").AsInt(0)
	assert.Equal(vba, 10)
	vca := pla.At("ab/bb/ca").AsInt(0)
	assert.Equal(vca, 100)
	vcb := pla.At("ab/bb/cb").AsInt(0)
	assert.Equal(vcb, 200)
}

// TestPayloadClone verifies the cloning of payloads together with
// modification of individual ones.
func TestPayloadClone(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	pla := event.NewPayload("a", 1, "b", "two", "c", 3.0)

	vaa := pla.At("a").AsInt(0)
	assert.Equal(vaa, 1)
	vba := pla.At("b").AsString("zero")
	assert.Equal(vba, "two")
	vca := pla.At("c").AsFloat64(0.0)
	assert.Equal(vca, 3.0)

	plb := pla.Clone("a", "4711", "d", "foo")

	vab := plb.At("a").AsInt(0)
	assert.Equal(vab, 4711)
	vbb := plb.At("b").AsString("zero")
	assert.Equal(vbb, "two")
	vcb := plb.At("c").AsFloat64(0.0)
	assert.Equal(vcb, 3.0)
	vdb := plb.At("d").AsString("bar")
	assert.Equal(vdb, "foo")

	assert.Length(plb.Keys(), 4)
}

// EOF
