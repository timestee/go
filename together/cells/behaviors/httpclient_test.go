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

// TestHTTPClientBehaviorGet tests the HTTP client behavior, here
// the GET method.
func TestHTTPClientBehaviorGet(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	sigc := asserts.MakeWaitChan()
	url := "https://api.github.com/repos/tideland/go/events"
	msh := mesh.New()
	defer msh.Stop()

	processor := func(accessor event.SinkAccessor) (*event.Payload, error) {
		ghEventTypes := []string{}
		accessor.Do(func(index int, evt *event.Event) error {
			return evt.Payload().At("data").AsPayload().Do(func(key string, value *event.Value) error {
				ghEventType := value.AsPayload().At("Type").AsString("<unknown>")
				ghEventTypes = append(ghEventTypes, key+"/"+ghEventType)
				return nil
			})
		})
		sigc <- ghEventTypes
		return nil, nil
	}

	msh.SpawnCells(
		behaviors.NewHTTPClientBehavior("github"),
		behaviors.NewCollectorBehavior("collector", 10, processor),
	)

	msh.Emit("", event.New(behaviors.TopicHTTPGet, "id", "test", "url", url))

	assert.WaitTested(sigc, func(v interface{}) error {
		ghEventTypes := v.([]string)
		for i, ghEventType := range ghEventTypes {
			assert.Logf("%d) %s", i, ghEventType)
		}
		return nil
	}, time.Second)
}

// EOF
