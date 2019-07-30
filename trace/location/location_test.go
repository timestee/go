// Tideland Go Library - Trace - Location - Unit Tests
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package location_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/trace/location"
)

//--------------------
// TESTS
//--------------------

// TestHere tests retrieving the location in a detailed
// way and as ID.
func TestHere(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	pkg, file, fn, line := location.Here(0)

	assert.Equal(pkg, "tideland.dev/go/trace/location_test")
	assert.Equal(file, "location_test.go")
	assert.Equal(fn, "TestHere")
	assert.Equal(line, 30)

	id := location.HereID(0)

	assert.Equal(id, "(tideland.dev/go/trace/location_test:location_test.go:TestHere:37)")

	code := location.HereCode("E", 0)

	assert.Equal(code, "ETGTLL41")
}

// TestOffset tests retrieving the location with an offset.
func TestOffset(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	id := there()

	assert.Equal(id, "(tideland.dev/go/trace/location_test:location_test.go:TestOffset:49)")

	id = nestedThere()

	assert.Equal(id, "(tideland.dev/go/trace/location_test:location_test.go:TestOffset:53)")

	id = nameless()

	assert.Equal(id, "(tideland.dev/go/trace/location_test:location_test.go:nameless.func1:97)")

	id = location.HereID(-5)

	assert.Equal(id, "(tideland.dev/go/trace/location_test:location_test.go:TestOffset:61)")
}

// TestCache tests the caching of locations.
func TestCache(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	for i := 0; i < 100; i++ {
		id := nameless()

		assert.Equal(id, "(tideland.dev/go/trace/location_test:location_test.go:nameless.func1:97)")
	}
}

//--------------------
// HELPER
//--------------------

// there returns the id at the location of the caller.
func there() string {
	return location.HereID(1)
}

// nestedThere returns the id at the location of the caller but inside a local func.
func nestedThere() string {
	where := func() string {
		return location.HereID(2)
	}
	return where()
}

// nameless returns the id from calling a nested nameless function w/o an offset.
func nameless() string {
	noname := func() string {
		return location.HereID(0)
	}
	return noname()
}

// EOF
