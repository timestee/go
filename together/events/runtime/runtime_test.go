// Tideland Go Library - Together - Events - Runtime - Unit Tests
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package runtime_test

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/events/runtime"
)

//--------------------
// TESTS
//--------------------

// TestSpawnProcessors verifies starting the runtime, spawning some
// processors, and stops the runtime.
func TestSpawnProcessors(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	rt := runtime.New()

	err := rt.SpawnProcessors(
		NewTestProcessor("foo", nil),
		NewTestProcessor("bar", nil),
		NewTestProcessor("baz", nil),
	)
	assert.NoError(err)
	ids := rt.ProcessorIDs()
	assert.Length(ids, 3)
	assert.Contains(ids, "foo")
	assert.Contains(ids, "bar")
	assert.Contains(ids, "baz")

	err = rt.Stop()
	assert.NoError(err)
}

// TestEmitEvents verifies emitting some events to an
// individual
func TestEmitEvents(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	rt := runtime.New()

	datasC := make(chan []string, 1)
	err := rt.SpawnProcessors(NewTestProcessor("foo", datasC))
	assert.NoError(err)

	rt.Emit("foo", runtime.Event{"add", "a"})
	rt.Emit("foo", runtime.Event{"add", "b"})
	rt.Emit("foo", runtime.Event{"add", "c"})
	rt.Emit("foo", runtime.Event{"send", ""})

	datas := <-datasC

	assert.Length(datas, 3)

	err = rt.Stop()
	assert.NoError(err)
}

// TestSubscription verifies the subscription mechanis.
func TestSubscription(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	rt := runtime.New()

	fooDatasC := make(chan []string, 1)
	barDatasC := make(chan []string, 1)
	bazDatasC := make(chan []string, 1)
	err := rt.SpawnProcessors(
		NewTestProcessor("foo", fooDatasC),
		NewTestProcessor("bar", barDatasC),
		NewTestProcessor("baz", bazDatasC),
	)
	assert.NoError(err)
	rt.Subscribe("foo", "bar")
	rt.Subscribe("bar", "baz")

	rt.Emit("foo", runtime.Event{"add", "a"})
	rt.Emit("foo", runtime.Event{"add", "b"})
	rt.Emit("foo", runtime.Event{"add", "c"})

	rt.Emit("foo", runtime.Event{"length?", ""})

	<-fooDatasC

	rt.Emit("bar", runtime.Event{"send", ""})

	datas := <-barDatasC

	assert.Length(datas, 1)

	rt.Emit("bar", runtime.Event{"length?", ""})

	<-barDatasC

	rt.Emit("baz", runtime.Event{"send", ""})

	datas = <-bazDatasC

	assert.Length(datas, 1)
	assert.Equal(datas[0], `length "bar" = 1`)

	err = rt.Stop()
	assert.NoError(err)
}

//--------------------
// HELPERS
//--------------------

type TestProcessor struct {
	id      string
	datas   []string
	datasC  chan []string
	emitter runtime.Emitter
}

func NewTestProcessor(id string, datasC chan []string) *TestProcessor {
	return &TestProcessor{
		id:     id,
		datasC: datasC,
	}
}

func (tp *TestProcessor) ID() string {
	return tp.id
}

func (tp *TestProcessor) Init(emitter runtime.Emitter) error {
	tp.emitter = emitter
	return nil
}

func (tp *TestProcessor) Terminate() error {
	return nil
}

func (tp *TestProcessor) Process(event runtime.Event) {
	switch event.Topic {
	case "add":
		tp.datas = append(tp.datas, event.Data)
	case "clear":
		tp.datas = nil
	case "length?":
		tp.emitter.Emit(runtime.Event{
			Topic: "add",
			Data:  fmt.Sprintf("length %q = %d", tp.id, len(tp.datas)),
		})
		close(tp.datasC)
	case "send":
		tp.datasC <- tp.datas
	}
}

func (tp *TestProcessor) Recover(r interface{}) error {
	return nil
}

// EOF
