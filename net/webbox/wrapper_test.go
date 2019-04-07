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
	"errors"
	"net/http"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/audit/environments"
	"tideland.dev/go/net/jwt/token"
	"tideland.dev/go/net/webbox"
)

//--------------------
// TESTS
//--------------------

// TestInvalidMethodWrapper tests the panic if the past handler for the
// MethodWrapper is invalid.
func TestInvalidMethodWrapper(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	assert.Panics(func() {
		webbox.NewMethodWrapper(nil)
	}, "webbox: nil handler")
}

// TestMethodWrapper tests the wrapping of a handler for the dispatching
// of HTTP methods.
func TestMethodWrapper(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := startWebAsserter(assert)
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
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := startWebAsserter(assert)
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
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := startWebAsserter(assert)
	defer wa.Close()

	nw := webbox.NewNestedWrapper()

	nw.AppendFunc(func(w http.ResponseWriter, r *http.Request) {
		reply := ""
		if f, ok := webbox.PathField(r, 0); ok {
			reply = f
		}
		if f, ok := webbox.PathField(r, 1); ok {
			reply += " && " + f
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
			reply += " && " + f
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
			body:       "orders && 4711",
		}, {
			path:       "/orders/4711/items",
			statusCode: http.StatusOK,
			body:       "items",
		}, {
			path:       "/orders/4711/items/1",
			statusCode: http.StatusOK,
			body:       "items && 1",
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

// TestJWTWrapper tests access control by usage of the
// JWT wrapper.
func TestJWTWrapper(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	wa := startWebAsserter(assert)
	defer wa.Close()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(environments.HeaderContentType, environments.ContentTypeTextPlain)
		w.Write([]byte("request passed"))
		w.WriteHeader(http.StatusOK)
	})
	jwtWrapper := webbox.NewJWTWrapper(handler, &webbox.JWTWrapperConfig{
		Key: []byte("secret"),
		Gatekeeper: func(w http.ResponseWriter, r *http.Request, claims token.Claims) error {
			access, ok := claims.GetString("access")
			if !ok || access != "allowed" {
				return errors.New("access is not allowed")
			}
			return nil
		},
	})

	wa.Handle("/", jwtWrapper)

	tests := []struct {
		key         string
		accessClaim string
		statusCode  int
		body        string
	}{
		{
			key:         "",
			accessClaim: "",
			statusCode:  http.StatusUnauthorized,
			body:        "request contains no authorization header",
		}, {
			key:         "unknown",
			accessClaim: "allowed",
			statusCode:  http.StatusUnauthorized,
			body:        "cannot verify the signature",
		}, {
			key:         "secret",
			accessClaim: "allowed",
			statusCode:  http.StatusOK,
			body:        "request passed",
		}, {
			key:         "unknown",
			accessClaim: "forbidden",
			statusCode:  http.StatusUnauthorized,
			body:        "cannot verify the signature",
		}, {
			key:         "secret",
			accessClaim: "forbidden",
			statusCode:  http.StatusUnauthorized,
			body:        "access rejected by gatekeeper: access is not allowed",
		},
	}
	for i, test := range tests {
		assert.Logf("test case #%d: %s / %s", i, test.key, test.accessClaim)
		wreq := wa.CreateRequest(http.MethodGet, "/")
		if test.key != "" && test.accessClaim != "" {
			claims := token.NewClaims()
			claims.Set("access", test.accessClaim)
			jwt, err := token.Encode(claims, []byte(test.key), token.HS512)
			assert.NoError(err)
			wreq.Header().Set("Authorization", "Bearer "+jwt.String())
		}
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
