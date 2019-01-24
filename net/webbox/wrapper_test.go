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

// TestInvalidMethodWrapper tests the panic if the past handler for the
// MethodWrapper is invalid.
func TestInvalidMethodWrapper(t *testing.T) {
	assert := asserts.NewTesting(t, true)

	assert.Panics(func() {
		webbox.NewMethodWrapper(nil)
	}, "webbox: nil handler")
}

// TestMethodWrapper tests the wrapping of a handler for the dispatching
// of HTTP methods.
func TestMethodWrapper(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	wa := StartWebAsserter(assert)
	defer wa.Close()

	wa.Handle("/mwrap/", webbox.NewMethodWrapper(MethodHandler{}))

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

// TestNestedWrapperNoHandler tests the mapping of requests to a
// nested wrapper w/o sub-handlers.
func TestNestedWrapperNoHandler(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	wa := StartWebAsserter(assert)
	defer wa.Close()

	nw := webbox.NewNestedWrapper()

	wa.Handle("/foo", nw)

	wreq := wa.CreateRequest(http.MethodGet, "/foo")
	wresp := wreq.Do()

	wresp.AssertStatusCodeEquals(http.StatusNotFound)
	wresp.AssertBodyMatches("")
}

// TestNestedWrapper tests the mapping of requests to a number of
// nested individual handlers.
func TestNestedWrapper(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	wa := StartWebAsserter(assert)
	defer wa.Close()

	nw := webbox.NewNestedWrapper()

	nw.AppendFunc(func(w http.ResponseWriter, r *http.Request) {
		reply := ""
		if f, ok := webbox.PathField(r, 0); ok {
			reply = f
		}
		if f, ok := webbox.PathField(r, 1); ok {
			reply += "/" + f
		}
		w.Header().Add(environments.HeaderContentType, environments.ContentTypeTextPlain)
		w.Write([]byte(reply))
		w.WriteHeader(http.StatusOK)
	})
	nw.AppendFunc(func(w http.ResponseWriter, r *http.Request) {
		reply := ""
		if f, ok := webbox.PathField(r, 2); ok {
			reply = f
		}
		if f, ok := webbox.PathField(r, 3); ok {
			reply += "/" + f
		}
		w.Header().Add(environments.HeaderContentType, environments.ContentTypeTextPlain)
		w.Write([]byte(reply))
		w.WriteHeader(http.StatusOK)
	})

	wa.Handle("/orders/", nw)
	wa.Handle("/", nw)

	tests := []struct {
		path       string
		statusCode int
		body       string
	}{
		{
			path:       "/",
			statusCode: http.StatusOK,
			body:       "",
		}, {
			path:       "/orders/",
			statusCode: http.StatusOK,
			body:       "orders",
		}, {
			path:       "/orders/4711",
			statusCode: http.StatusOK,
			body:       "orders/4711",
		}, {
			path:       "/orders/4711/items",
			statusCode: http.StatusOK,
			body:       "items",
		}, {
			path:       "/orders/4711/items/1",
			statusCode: http.StatusOK,
			body:       "items/1",
		}, {
			path:       "/orders/4711/items/1/nothingelse",
			statusCode: http.StatusNotFound,
			body:       "",
		},
	}
	for i, test := range tests {
		assert.Logf("test case #%d: %s", i, test.path)
		wreq := wa.CreateRequest(http.MethodGet, test.path)
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
