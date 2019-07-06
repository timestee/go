// Tideland Go Library - Together - Mesh - Nodes - Unit Tests
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package nodes_test

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/mesh/nodes"
)

//--------------------
// TESTS
//--------------------

// TestSpawnNodes verifies starting the mesh, spawning some
// nodes, and stops the mesh.
func TestSpawnNodes(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	mesh := nodes.NewMesh()

	err := mesh.SpawnNodes(
		NewTestBehavior("foo", nil),
		NewTestBehavior("bar", nil),
		NewTestBehavior("baz", nil),
	)
	assert.NoError(err)
	ids := mesh.NodeIDs()
	assert.Length(ids, 3)
	assert.Contains(ids, "foo")
	assert.Contains(ids, "bar")
	assert.Contains(ids, "baz")

	err = mesh.Stop()
	assert.NoError(err)
}

// TestEmitEvents verifies emitting some events to a node.
func TestEmitEvents(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	mesh := nodes.NewMesh()

	datasC := make(chan []string, 1)
	err := mesh.SpawnNodes(NewTestBehavior("foo", datasC))
	assert.NoError(err)

	mesh.Emit("foo", nodes.Event{"add", "a"})
	mesh.Emit("foo", nodes.Event{"add", "b"})
	mesh.Emit("foo", nodes.Event{"add", "c"})
	mesh.Emit("foo", nodes.Event{"send", ""})

	datas := <-datasC

	assert.Length(datas, 3)

	err = mesh.Stop()
	assert.NoError(err)
}

// TestSubscription verifies the subscription mechanics.
func TestSubscription(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	mesh := nodes.NewMesh()

	fooDatasC := make(chan []string, 1)
	barDatasC := make(chan []string, 1)
	bazDatasC := make(chan []string, 1)
	err := mesh.SpawnNodes(
		NewTestBehavior("foo", fooDatasC),
		NewTestBehavior("bar", barDatasC),
		NewTestBehavior("baz", bazDatasC),
	)
	assert.NoError(err)
	mesh.Subscribe("foo", "bar")
	mesh.Subscribe("bar", "baz")

	mesh.Emit("foo", nodes.Event{"add", "a"})
	mesh.Emit("foo", nodes.Event{"add", "b"})
	mesh.Emit("foo", nodes.Event{"add", "c"})

	mesh.Emit("foo", nodes.Event{"length?", ""})

	<-fooDatasC

	mesh.Emit("bar", nodes.Event{"send", ""})

	datas := <-barDatasC

	assert.Length(datas, 1)

	mesh.Emit("bar", nodes.Event{"length?", ""})

	<-barDatasC

	mesh.Emit("baz", nodes.Event{"send", ""})

	datas = <-bazDatasC

	assert.Length(datas, 1)
	assert.Equal(datas[0], `length "bar" = 1`)

	err = mesh.Stop()
	assert.NoError(err)
}

//--------------------
// HELPERS
//--------------------

type TestBehavior struct {
	id      string
	datas   []string
	datasC  chan []string
	emitter nodes.Emitter
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

func (tb *TestBehavior) Init(emitter nodes.Emitter) error {
	tb.emitter = emitter
	return nil
}

func (tb *TestBehavior) Terminate() error {
	return nil
}

func (tb *TestBehavior) Process(event nodes.Event) {
	switch event.Topic {
	case "add":
		tb.datas = append(tb.datas, event.Data)
	case "clear":
		tb.datas = nil
	case "length?":
		tb.emitter.Emit(nodes.Event{
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
