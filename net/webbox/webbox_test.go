// Tideland Go Library - Network - Web Toolbox - Unit Tests
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package webbox_test

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"testing"

	"tideland.one/go/audit/asserts"
	"tideland.one/go/audit/environments"
	"tideland.one/go/net/webbox"
)

//--------------------
// TESTS
//--------------------

// TestPathField tests the checking and extrecting of a field out of
// a request path.
func TestPathField(t *testing.T) {
	assert := asserts.NewTesting(t, true)

	r, err := http.NewRequest(http.MethodGet, "http://localhost/foo", nil)
	assert.NoError(err)
	f, ok := webbox.PathField(r, 1)
	assert.True(ok)
	assert.Equal(f, "foo")
	f, ok = webbox.PathField(r, 2)
	assert.False(ok)
	assert.Equal(f, "")

	r, err = http.NewRequest(http.MethodGet, "http://localhost/orders/4711/items/1", nil)
	assert.NoError(err)
	f, ok = webbox.PathField(r, 2)
	assert.True(ok)
	assert.Equal(f, "4711")
	f, ok = webbox.PathField(r, 4)
	assert.True(ok)
	assert.Equal(f, "1")
	f, ok = webbox.PathField(r, 5)
	assert.False(ok)
	assert.Equal(f, "")

}

//--------------------
// WEB ASSERTER AND HELPING HANDLER
//--------------------

// StartTestServer initialises and starts the asserter for the tests.
func StartWebAsserter(assert *asserts.Asserts) *environments.WebAsserter {
	wa := environments.NewWebAsserter(assert, func(r *http.Request) (string, error) {
		return r.URL.Path, nil
	})
	return wa
}

// MakeMethodEcho creates a handler echoing the HTTP method.
func MakeMethodEcho(assert *asserts.Asserts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reply := "METHOD: " + r.Method + "!"
		w.Header().Add(environments.HeaderContentType, environments.ContentTypeTextPlain)
		w.Write([]byte(reply))
		w.WriteHeader(http.StatusOK)
	}
}

// EOF
