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
	"strings"
	"testing"

	"tideland.one/go/audit/asserts"
	"tideland.one/go/audit/web"
)

//--------------------
// TESTS
//--------------------

// TestSimpleRequests tests simple requests to individual handlers.
func TestSimpleRequests(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	ts := StartTestServer()
	defer ts.Close()

	tests := []struct {
		method     string
		path       string
		statusCode int
		body       string
	}{
		{
			method:     http.MethodGet,
			path:       "/hello/world",
			statusCode: http.StatusOK,
			body:       "Hello, World!",
		}, {
			method:     http.MethodGet,
			path:       "/hello/tester",
			statusCode: http.StatusOK,
			body:       "Hello, Tester!",
		}, {
			method:     http.MethodPost,
			path:       "/hello/postman",
			statusCode: http.StatusOK,
			body:       "Hello, Postman!",
		}, {
			method:     http.MethodOptions,
			path:       "/path/does/not/exist",
			statusCode: http.StatusInternalServerError,
		},
	}
	for _, test := range tests {
		req := web.NewRequest(assert, test.method, test.path)
		resp := ts.DoRequest(req)
		resp.AssertStatusCodeEquals(test.statusCode)
		if test.body != "" {
			resp.AssertBodyMatches(test.body)
		}
	}
}

//--------------------
// MUX MAPPER AND HANDLER
//--------------------

// StartTestServer initialises and starts the test server.
func StartTestServer() *web.TestServer {
	mux := web.NewMultiplexer(Mapper)
	mux.Register("get/hello/world", MakeHelloWorldHandler("World"))
	mux.Register("get/hello/tester", MakeHelloWorldHandler("Tester"))
	mux.Register("post/hello/postman", MakeHelloWorldHandler("Postman"))

	return web.StartServer(mux)
}

// Mapper returns the ID for the test handler to user.
func Mapper(r *http.Request) (string, error) {
	return strings.ToLower(r.Method + r.URL.Path), nil
}

// MakeHelloWorldHandler creates a "Hello, World" handler.
func MakeHelloWorldHandler(who string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reply := "Hello, " + who + "!"
		w.Write([]byte(reply))
		w.WriteHeader(http.StatusOK)
	}
}

// EOF
