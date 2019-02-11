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
	"tideland.one/go/net/jwt/token"
	"tideland.one/go/net/webbox"
)

//--------------------
// TESTS
//--------------------

// TestPathFields tests the splitting of request paths into fields.
func TestPathFields(t *testing.T) {
	assert := asserts.NewTesting(t, true)

	r, err := http.NewRequest(http.MethodGet, "http://localhost/", nil)
	assert.NoError(err)
	fs := webbox.PathFields(r)
	assert.Length(fs, 0)

	r, err = http.NewRequest(http.MethodGet, "http://localhost/foo", nil)
	assert.NoError(err)
	fs = webbox.PathFields(r)
	assert.Length(fs, 1)

	r, err = http.NewRequest(http.MethodGet, "http://localhost/foo/bar", nil)
	assert.NoError(err)
	fs = webbox.PathFields(r)
	assert.Length(fs, 2)

	r, err = http.NewRequest(http.MethodGet, "http://localhost/foo/bar?yadda=1", nil)
	assert.NoError(err)
	fs = webbox.PathFields(r)
	assert.Length(fs, 2)
}

// TestPathField tests the checking and extrecting of a field out of
// a request path.
func TestPathField(t *testing.T) {
	assert := asserts.NewTesting(t, true)

	r, err := http.NewRequest(http.MethodGet, "http://localhost/foo", nil)
	assert.NoError(err)
	f, ok := webbox.PathField(r, 0)
	assert.True(ok)
	assert.Equal(f, "foo")
	f, ok = webbox.PathField(r, 1)
	assert.False(ok)
	assert.Equal(f, "")

	r, err = http.NewRequest(http.MethodGet, "http://localhost/orders/4711/items/1", nil)
	assert.NoError(err)
	f, ok = webbox.PathField(r, 1)
	assert.True(ok)
	assert.Equal(f, "4711")
	f, ok = webbox.PathField(r, 3)
	assert.True(ok)
	assert.Equal(f, "1")
	f, ok = webbox.PathField(r, 5)
	assert.False(ok)
	assert.Equal(f, "")
}

//--------------------
// WEB ASSERTER AND HELPERS
//--------------------

// StartTestServer initialises and starts the asserter for the tests.
func startWebAsserter(assert *asserts.Asserts) *environments.WebAsserter {
	wa := environments.NewWebAsserter(assert)
	return wa
}

// makeMethodEcho creates a handler echoing the HTTP method.
func makeMethodEcho(assert *asserts.Asserts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reply := "METHOD: " + r.Method + "!"
		w.Header().Add(environments.HeaderContentType, environments.ContentTypeTextPlain)
		w.Write([]byte(reply))
		w.WriteHeader(http.StatusOK)
	}
}

// createToken creates a test JWT with the passed access claim.
func createToken(assert *asserts.Asserts, access string) *token.JWT {
	claims := token.NewClaims()
	claims.Set("access", access)
	jwt, err := token.Encode(claims, []byte("secret"), token.HS512)
	assert.NoError(err)
	return jwt
}

// data is used in marshalling tests.
type data struct {
	Number int
	Name   string
	Tags   []string
}

// EOF
