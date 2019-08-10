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
// CONSTANTS
//--------------------

const waitTimeout = 20 * time.Millisecond

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

	ids := msh.Cells()
	assert.Length(ids, 3)
	assert.Contains(ids, "foo")
	assert.Contains(ids, "bar")
	assert.Contains(ids, "baz")

	err = msh.Stop()
	assert.NoError(err)
}

// TestStopCells verifies stopping some cells.
func TestStopCells(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	// Initial environment with subscriptions.
	err := msh.SpawnCells(
		NewTestBehavior("foo", nil),
		NewTestBehavior("bar", nil),
		NewTestBehavior("baz", nil),
	)
	assert.NoError(err)

	ids := msh.Cells()
	assert.Length(ids, 3)
	assert.Contains(ids, "foo")
	assert.Contains(ids, "bar")
	assert.Contains(ids, "baz")

	msh.Subscribe("foo", "bar", "baz")

	fooS, err := msh.Subscribers("foo")
	assert.NoError(err)
	assert.Length(fooS, 2)
	assert.Contains(fooS, "bar")
	assert.Contains(fooS, "baz")

	// Stopping shall unsubscribe too.
	err = msh.StopCells("baz")

	ids = msh.Cells()
	assert.Length(ids, 2)
	assert.Contains(ids, "foo")
	assert.Contains(ids, "bar")

	fooS, err = msh.Subscribers("foo")
	assert.NoError(err)
	assert.Length(fooS, 1)
	assert.Contains(fooS, "bar")

	err = msh.Stop()
	assert.NoError(err)
}

// TestEmitEvents verifies emitting some events to a node.
func TestEmitEvents(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	threeTest := mkLenTest(assert, 3)
	msh := mesh.New()

	fooC := asserts.MakeWaitChan()

	err := msh.SpawnCells(NewTestBehavior("foo", fooC))
	assert.NoError(err)

	msh.Emit("foo", event.New("add", "x", "a"))
	msh.Emit("foo", event.New("add", "x", "b"))
	msh.Emit("foo", event.New("add", "x", "c"))
	msh.Emit("foo", event.New("send"))

	assert.WaitTested(fooC, threeTest, waitTimeout)

	err = msh.Stop()
	assert.NoError(err)
}

// TestBroadcastEvents verifies broadcasting some events to a node.
func TestBroadcastEvents(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	threeTest := mkLenTest(assert, 3)
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

	msh.Broadcast(event.New("add", "x", "a"))
	msh.Broadcast(event.New("add", "x", "b"))
	msh.Broadcast(event.New("add", "x", "c"))
	msh.Broadcast(event.New("send"))

	assert.WaitTested(fooC, threeTest, waitTimeout)
	assert.WaitTested(barC, threeTest, waitTimeout)
	assert.WaitTested(bazC, threeTest, waitTimeout)

	err = msh.Stop()
	assert.NoError(err)
}

// TestSubscribe verifies the subscription of cells.
func TestSubscribe(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	oneTest := mkLenTest(assert, 1)
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

	fooS, err := msh.Subscribers("foo")
	assert.NoError(err)
	assert.Length(fooS, 1)
	assert.Contains(fooS, "bar")

	msh.Emit("foo", event.New("add", "x"))
	msh.Emit("foo", event.New("length"))
	msh.Emit("foo", event.New("send"))
	assert.WaitTested(fooC, oneTest, waitTimeout)

	msh.Emit("bar", event.New("length"))
	msh.Emit("bar", event.New("send"))
	assert.WaitTested(barC, oneTest, waitTimeout)

	msh.Emit("baz", event.New("send"))
	assert.WaitTested(bazC, oneTest, waitTimeout)

	err = msh.Stop()
	assert.NoError(err)
}

// TestUnsubscribe verifies the unsubscription of cells.
func TestUnsubscribe(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	zeroTest := mkLenTest(assert, 0)
	oneTest := mkLenTest(assert, 1)
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

	fooS, err := msh.Subscribers("foo")
	assert.NoError(err)
	assert.Length(fooS, 2)
	assert.Contains(fooS, "bar")
	assert.Contains(fooS, "baz")

	msh.Emit("foo", event.New("add", "x"))
	msh.Emit("foo", event.New("length"))
	msh.Emit("foo", event.New("send"))
	assert.WaitTested(fooC, oneTest, waitTimeout)

	msh.Emit("bar", event.New("send"))
	assert.WaitTested(barC, oneTest, waitTimeout)

	msh.Emit("baz", event.New("send"))
	assert.WaitTested(bazC, oneTest, waitTimeout)

	// Unsubscribe baz, test both, expect zero in baz.
	msh.Unsubscribe("foo", "baz")

	fooS, err = msh.Subscribers("foo")
	assert.NoError(err)
	assert.Length(fooS, 1)
	assert.Contains(fooS, "bar")

	msh.Emit("foo", event.New("add", "x"))
	msh.Emit("foo", event.New("length"))
	msh.Emit("foo", event.New("send"))
	assert.WaitTested(fooC, oneTest, waitTimeout)

	msh.Emit("bar", event.New("send"))
	assert.WaitTested(barC, oneTest, waitTimeout)

	msh.Emit("baz", event.New("send"))
	assert.WaitTested(bazC, zeroTest, waitTimeout)

	err = msh.Stop()
	assert.NoError(err)
}

// TestInvalidSubscriptions verifies the invalid (un)subscriptions of cells.
func TestInvalidSubscriptions(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	err := msh.SpawnCells(
		NewTestBehavior("foo", nil),
		NewTestBehavior("bar", nil),
	)
	assert.NoError(err)

	err = msh.Subscribe("foo", "bar", "baz")
	assert.ErrorMatch(err, ".*cannot find cell.*")

	err = msh.Subscribe("foo", "bar")
	assert.NoError(err)

	err = msh.Unsubscribe("foo", "bar", "baz")
	assert.ErrorMatch(err, ".*cannot find cell.*")

	err = msh.Unsubscribe("foo", "bar")
	assert.NoError(err)

	err = msh.Unsubscribe("foo", "bar")
	assert.NoError(err)

	err = msh.Stop()
	assert.NoError(err)
}

// TestSubscriberIDs verifies the retrieval of subscriber IDs.
func TestSubscriberIDs(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	err := msh.SpawnCells(
		NewTestBehavior("foo", nil),
		NewTestBehavior("bar", nil),
		NewTestBehavior("baz", nil),
	)
	assert.NoError(err)

	err = msh.Subscribe("foo", "bar", "baz")
	assert.NoError(err)

	subscriberIDs, err := msh.Subscribers("foo")
	assert.NoError(err)
	assert.Length(subscriberIDs, 2)

	subscriberIDs, err = msh.Subscribers("bar")
	assert.NoError(err)
	assert.Length(subscriberIDs, 0)

	err = msh.Unsubscribe("foo", "baz")
	assert.NoError(err)

	subscriberIDs, err = msh.Subscribers("foo")
	assert.NoError(err)
	assert.Length(subscriberIDs, 1)

	err = msh.Stop()
	assert.NoError(err)
}

//--------------------
// HELPERS
//--------------------

func mkLenTest(assert *asserts.Asserts, l int) func(interface{}) error {
	return func(data interface{}) error {
		if !assert.Equal(data, l) {
			return errors.New("not 3")
		}
		return nil
	}
}

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

func (tb *TestBehavior) Process(evt *event.Event) error {
	switch evt.Topic() {
	case "add":
		x := evt.Payload().At("x").AsString("-")
		tb.datas = append(tb.datas, x)
	case "length":
		tb.emitter.EmitAll(event.New(
			"add",
			"x", fmt.Sprintf("length %q = %d", tb.id, len(tb.datas)),
		))
	case "send":
		tb.dataC <- len(tb.datas)
		tb.datas = nil
	}
	return nil
}

func (tb *TestBehavior) Recover(r interface{}) error {
	return nil
}

// EOF
