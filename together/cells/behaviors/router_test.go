// Tideland Go Library - Together - Cells - Behaviors - Unit Tests
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors_test // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"strings"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/cells/behaviors"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestRouterBehavior tests the router behavior.
func TestRouterBehavior(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sigc := asserts.MakeWaitChan()
	msh := mesh.New()
	defer msh.Stop()

	router := func(evt *event.Event) []string {
		return strings.Split(evt.Payload().At("ids").AsString(""), "/")
	}
	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		sigc <- accessor.Len()
		return nil, nil
	}

	msh.SpawnCells(
		behaviors.NewMeshRouterBehavior("router", router),
		behaviors.NewCollectorBehavior("test-1", 10, processor),
		behaviors.NewCollectorBehavior("test-2", 10, processor),
		behaviors.NewCollectorBehavior("test-3", 10, processor),
		behaviors.NewCollectorBehavior("test-4", 10, processor),
		behaviors.NewCollectorBehavior("test-5", 10, processor),
	)
	msh.Subscribe("router", "test-1", "test-2", "test-3", "test-4", "test-5")

	msh.Emit("router", event.New("route-it", "ids", "test-1/test-2"))
	msh.Emit("router", event.New("route-it", "ids", "test-1/test-2/test-3"))
	msh.Emit("router", event.New("route-it", "ids", "test-3/test-4/test-5"))
	msh.Emit("router", event.New("route-it", "ids", "unknown"))

	time.Sleep(100 * time.Millisecond)

	test := func(id string, l int) {
		msh.Emit(id, event.New(event.TopicProcess))
		assert.Wait(sigc, l, time.Second)
	}

	test("test-1", 2)
	test("test-2", 2)
	test("test-3", 2)
	test("test-4", 1)
	test("test-5", 1)
}

// EOF
