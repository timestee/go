// Tideland Go Library - Together - Cells - Mesh - Unit Tests
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package mesh_test

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"testing"

	"tideland.dev/go/audit/asserts"
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

	datasC := make(chan []string, 1)
	err := msh.SpawnCells(NewTestBehavior("foo", datasC))
	assert.NoError(err)

	msh.Emit("foo", mesh.Event{"add", "a"})
	msh.Emit("foo", mesh.Event{"add", "b"})
	msh.Emit("foo", mesh.Event{"add", "c"})
	msh.Emit("foo", mesh.Event{"send", ""})

	datas := <-datasC

	assert.Length(datas, 3)

	err = msh.Stop()
	assert.NoError(err)
}

// TestSubscription verifies the subscription mechanics.
func TestSubscription(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msh := mesh.New()

	fooDatasC := make(chan []string, 1)
	barDatasC := make(chan []string, 1)
	bazDatasC := make(chan []string, 1)
	err := msh.SpawnCells(
		NewTestBehavior("foo", fooDatasC),
		NewTestBehavior("bar", barDatasC),
		NewTestBehavior("baz", bazDatasC),
	)
	assert.NoError(err)

	msh.Subscribe("foo", "bar")
	msh.Subscribe("bar", "baz")

	msh.Emit("foo", mesh.Event{"add", "a"})
	msh.Emit("foo", mesh.Event{"add", "b"})
	msh.Emit("foo", mesh.Event{"add", "c"})
	msh.Emit("foo", mesh.Event{"length?", ""})

	<-fooDatasC

	msh.Emit("bar", mesh.Event{"send", ""})

	datas := <-barDatasC

	assert.Length(datas, 1)

	msh.Emit("bar", mesh.Event{"length?", ""})

	<-barDatasC

	msh.Emit("baz", mesh.Event{"send", ""})

	datas = <-bazDatasC

	assert.Length(datas, 1)
	assert.Equal(datas[0], `length "bar" = 1`)

	err = msh.Stop()
	assert.NoError(err)
}

//--------------------
// HELPERS
//--------------------

type TestBehavior struct {
	id      string
	datas   []string
	datasC  chan []string
	emitter mesh.Emitter
}

func NewTestBehavior(id string, datasC chan []string) *TestBehavior {
	return &TestBehavior{
		id:     id,
		datasC: datasC,
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

func (tb *TestBehavior) Process(event mesh.Event) {
	switch event.Topic {
	case "add":
		tb.datas = append(tb.datas, event.Data)
	case "clear":
		tb.datas = nil
	case "length?":
		tb.emitter.Emit(mesh.Event{
			Topic: "add",
			Data:  fmt.Sprintf("length %q = %d", tb.id, len(tb.datas)),
		})
		close(tb.datasC)
	case "send":
		tb.datasC <- tb.datas
	}
}

func (tb *TestBehavior) Recover(r interface{}) error {
	return nil
}

// EOF
