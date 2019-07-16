// Tideland Go Library - Together - Cells - Mesh - Unit Tests
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package mesh_test // import "tideland.dev/go/together/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestSpawnCells verifies starting the mesh, spawning some
// cells, and stops the mesh.
func TestSpawnCells(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	err := msh.SpawnCells(
		NewTestBehavior("foo", nil),
		NewTestBehavior("bar", nil),
		NewTestBehavior("baz", nil),
	)
	assert.NoError(err)

	ids := msh.CellIDs()
	assert.Length(ids, 3)
	assert.Contains(ids, "foo")
	assert.Contains(ids, "bar")
	assert.Contains(ids, "baz")

	err = msh.Stop()
	assert.NoError(err)
}

// TestEmitEvents verifies emitting some events to a node.
func TestEmitEvents(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	lenC := asserts.MakeWaitChan()

	err := msh.SpawnCells(NewTestBehavior("foo", lenC))
	assert.NoError(err)

	msh.Emit("foo", event.New("add", "x", "a"))
	msh.Emit("foo", event.New("add", "x", "b"))
	msh.Emit("foo", event.New("add", "x", "c"))
	msh.Emit("foo", event.New("send"))

	dataLen := <-lenC

	assert.Equal(dataLen, 3)

	err = msh.Stop()
	assert.NoError(err)
}

// TestSubscribe verifies the subscription of cells.
func TestSubscribe(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	oneTest := func(data interface{}) error {
		if !assert.Equal(data, 1) {
			return errors.New("not 1")
		}
		return nil
	}
	msh := mesh.New()

	fooC := asserts.MakeWaitChan()
	barC := asserts.MakeWaitChan()
	bazC := asserts.MakeWaitChan()

	err := msh.SpawnCells(
		NewTestBehavior("foo", fooC),
		NewTestBehavior("bar", barC),
		NewTestBehavior("baz", bazC),
	)
	assert.NoError(err)

	msh.Subscribe("foo", "bar")
	msh.Subscribe("bar", "baz")

	msh.Emit("foo", event.New("add", "x"))
	msh.Emit("foo", event.New("length"))
	msh.Emit("foo", event.New("send"))

	<-fooC

	msh.Emit("bar", event.New("length"))
	msh.Emit("bar", event.New("send"))
	assert.WaitTested(barC, oneTest, 1*time.Second)

	msh.Emit("baz", event.New("send"))
	assert.WaitTested(bazC, oneTest, 1*time.Second)

	err = msh.Stop()
	assert.NoError(err)
}

// TestUnsubscribe verifies the unsubscription of cells.
func TestUnsubscribe(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	zeroTest := func(data interface{}) error {
		if !assert.Equal(data, 0) {
			return errors.New("not 0")
		}
		return nil
	}
	oneTest := func(data interface{}) error {
		if !assert.Equal(data, 1) {
			return errors.New("not 1")
		}
		return nil
	}
	msh := mesh.New()

	fooC := asserts.MakeWaitChan()
	barC := asserts.MakeWaitChan()
	bazC := asserts.MakeWaitChan()

	err := msh.SpawnCells(
		NewTestBehavior("foo", fooC),
		NewTestBehavior("bar", barC),
		NewTestBehavior("baz", bazC),
	)
	assert.NoError(err)

	// Subscribe bar and baz, test both.
	msh.Subscribe("foo", "bar", "baz")

	msh.Emit("foo", event.New("add", "x"))
	msh.Emit("foo", event.New("length"))
	msh.Emit("foo", event.New("send"))

	<-fooC

	msh.Emit("bar", event.New("send"))
	assert.WaitTested(barC, oneTest, 1*time.Second)

	msh.Emit("baz", event.New("send"))
	assert.WaitTested(bazC, oneTest, 1*time.Second)

	// Unsubscribe baz, test both, expect zero in baz.
	msh.Unsubscribe("foo", "baz")

	msh.Emit("foo", event.New("add", "x"))
	msh.Emit("foo", event.New("length"))
	msh.Emit("foo", event.New("send"))

	<-fooC

	msh.Emit("bar", event.New("send"))
	assert.WaitTested(barC, oneTest, 1*time.Second)

	msh.Emit("baz", event.New("send"))
	assert.WaitTested(bazC, zeroTest, 1*time.Second)

	err = msh.Stop()
	assert.NoError(err)
}

//--------------------
// HELPERS
//--------------------

type TestBehavior struct {
	id      string
	datas   []interface{}
	dataC   chan interface{}
	emitter mesh.Emitter
}

func NewTestBehavior(id string, dataC chan interface{}) *TestBehavior {
	return &TestBehavior{
		id:    id,
		dataC: dataC,
	}
}

func (tb *TestBehavior) ID() string {
	return tb.id
}

func (tb *TestBehavior) Init(emitter mesh.Emitter) error {
	tb.emitter = emitter
	return nil
}

func (tb *TestBehavior) Terminate() error {
	return nil
}

func (tb *TestBehavior) Process(evt *event.Event) {
	switch evt.Topic() {
	case "add":
		tb.datas = append(tb.datas, evt.Payload("x"))
	case "length":
		tb.emitter.Emit(event.New(
			"add",
			"x", fmt.Sprintf("length %q = %d", tb.id, len(tb.datas)),
		))
	case "send":
		tb.dataC <- len(tb.datas)
		tb.datas = nil
	}
}

func (tb *TestBehavior) Recover(r interface{}) error {
	return nil
}

// EOF
