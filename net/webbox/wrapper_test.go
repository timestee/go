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

// TestMethodWrapper tests the wrapping of a handler for the dispatching
// of HTTP methods.
func TestMethodWrapper(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	wa := StartWebAsserter(assert)
	defer wa.Close()

	wa.Register("/mwrap/", webbox.NewMethodWrapper(MethodHandler{}))

	tests := []struct {
		method     string
		statusCode int
		body       string
	}{
		{
			method:     http.MethodGet,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		}, {
			method:     http.MethodHead,
			statusCode: http.StatusBadRequest,
			body:       "",
		}, {
			method:     http.MethodPost,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		}, {
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			body:       "METHOD: PUT!",
		}, {
			method:     http.MethodPatch,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		}, {
			method:     http.MethodDelete,
			statusCode: http.StatusNoContent,
			body:       "",
		}, {
			method:     http.MethodConnect,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		}, {
			method:     http.MethodOptions,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		}, {
			method:     http.MethodTrace,
			statusCode: http.StatusBadRequest,
			body:       "bad request",
		},
	}
	for i, test := range tests {
		assert.Logf("test case #%d: %s", i, test.method)
		wreq := wa.CreateRequest(test.method, "/mwrap/")
		wresp := wreq.Do()
		wresp.AssertStatusCodeEquals(test.statusCode)
		wresp.AssertBodyMatches(test.body)
	}
}

//--------------------
// HELPING HANDLER
//--------------------

// MethodHelper provides some of the methods for the MethodWrapper.
type MethodHandler struct{}

func (mh MethodHandler) ServePut(w http.ResponseWriter, r *http.Request) {
	reply := "METHOD: " + r.Method + "!"
	w.Header().Add(environments.HeaderContentType, environments.ContentTypeTextPlain)
	w.Write([]byte(reply))
	w.WriteHeader(http.StatusOK)
}

func (mh MethodHandler) ServeDelete(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "", http.StatusNoContent)
}

func (mh MethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "bad request", http.StatusBadRequest)
}

// EOF
