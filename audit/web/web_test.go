// Tideland Go Library - Audit - Web - Unit Test
//
// Copyright (C) 2009-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package web_test

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"testing"

	"tideland.one/go/audit/asserts"
	"tideland.one/go/audit/web"
)

//--------------------
// TESTS
//--------------------

// TestGetJSON tests the GET command with a JSON result.
func TestGetJSON(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	ts := StartTestServer()

	req1 := web.NewRequest(assert, http.MethodGet, "/hello/world")
	resp1 := ts.DoRequest(req1)
	resp1.AssertStatusCodeEquals(http.StatusOK)
}

//--------------------
// MUX MAPPER AND HANDLER
//--------------------

// StartTestServer initialises and starts the test server.
func StartTestServer() *web.TestServer {
	mux := web.NewMultiplexer(Mapper)
	mux.Register("hello-world", MakeHelloWorldHandler())

	return web.StartServer(mux)
}

// Mapper returns the ID for the test handler to user.
func Mapper(req *http.Request) (string, error) {
	switch req.URL.Path {
	case "/hello/world":
		return "hello-world", nil
	}
	return "", nil
}

// MakeHelloWorldHandler creates a "Hello, World" handler.
func MakeHelloWorldHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reply := "Hello, World (" + r.Method + ")"
		w.Write([]byte(reply))
		w.WriteHeader(http.StatusOK)
	}
}

// EOF
