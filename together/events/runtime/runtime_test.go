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
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/events/runtime"
)

//--------------------
// TESTS
//--------------------

func TestSpawnProcessors(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	rt := runtime.New()

	err := rt.SpawnProcessors(
		NewTestProcessor("foo"),
		NewTestProcessor("bar"),
		NewTestProcessor("baz"),
	)
	assert.Nil(err)
	ids := rt.ProcessorIDs()
	assert.Length(ids, 3)
	assert.Contains(ids, "foo")
	assert.Contains(ids, "bar")
	assert.Contains(ids, "baz")

	err = rt.Stop()
	assert.Nil(err)
}

//--------------------
// HELPERS
//--------------------

type TestProcessor struct {
	id      string
	datas   []string
	emitter runtime.Emitter
}

func NewTestProcessor(id string) *TestProcessor {
	return &TestProcessor{
		id: id,
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

}

func (tp *TestProcessor) Recover(r interface{}) error {
	return nil
}

// EOF
