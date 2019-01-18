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
	"tideland.one/go/net/webbox"
)

//--------------------
// TESTS
//--------------------

// TestMethodMultiplexer tests the multiplexing of methods to different handler.
func TestMethodMultiplexer(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	wa := StartWebAsserter(assert)
	defer wa.Close()

	mmux := webbox.NewMethodMux()

	mmux.HandleFunc(http.MethodGet, MakeMethodEcho(assert))
	mmux.HandleFunc(http.MethodPatch, MakeMethodEcho(assert))
	mmux.HandleFunc(http.MethodOptions, MakeMethodEcho(assert))

	wa.Register("/mmux/", mmux)

	tests := []struct {
		method     string
		statusCode int
		body       string
	}{
		{
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			body:       "METHOD: GET!",
		}, {
			method:     http.MethodHead,
			statusCode: http.StatusMethodNotAllowed,
			body:       "",
		}, {
			method:     http.MethodPost,
			statusCode: http.StatusMethodNotAllowed,
			body:       "no matching method handler found",
		}, {
			method:     http.MethodPut,
			statusCode: http.StatusMethodNotAllowed,
			body:       "no matching method handler found",
		}, {
			method:     http.MethodPatch,
			statusCode: http.StatusOK,
			body:       "METHOD: PATCH!",
		}, {
			method:     http.MethodDelete,
			statusCode: http.StatusMethodNotAllowed,
			body:       "no matching method handler found",
		}, {
			method:     http.MethodConnect,
			statusCode: http.StatusMethodNotAllowed,
			body:       "no matching method handler found",
		}, {
			method:     http.MethodOptions,
			statusCode: http.StatusOK,
			body:       "METHOD: OPTIONS!",
		}, {
			method:     http.MethodTrace,
			statusCode: http.StatusMethodNotAllowed,
			body:       "no matching method handler found",
		},
	}
	for i, test := range tests {
		assert.Logf("test case #%d: %s", i, test.method)
		wreq := wa.CreateRequest(test.method, "/mmux/")
		wresp := wreq.Do()
		wresp.AssertStatusCodeEquals(test.statusCode)
		wresp.AssertBodyMatches(test.body)
	}
}

// EOF
