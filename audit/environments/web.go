// Tideland Go Library - Audit - Environments
//
// Copyright (C) 2012-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package environments

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"net/http/httptest"

	"tideland.one/go/audit/asserts"
)

//--------------------
// CONSTANTS
//--------------------

const (
	HeaderAccept      = "Accept"
	HeaderContentType = "Content-Type"

	ContentTypeTextPlain       = "text/plain"
	ContentTypeApplicationJSON = "application/json"
	ContentTypeApplicationXML  = "application/xml"
)

//--------------------
// VALUES
//--------------------

// Values wraps header, cookie, query, and form values.
type Values struct {
	wa   *WebAsserter
	data map[string][]string
}

// newValues creates a new values instance.
func newValues(wa *WebAsserter) *Values {
	vs := &Values{
		wa:   wa,
		data: make(map[string][]string),
	}
	return vs
}

// Add adds or appends a value to a named field.
func (vs *Values) Add(key, value string) {
	kd := append(vs.data[key], value)
	vs.data[key] = kd
}

// Set sets value of a named field.
func (vs *Values) Set(key, value string) {
	vs.data[key] = []string{value}
}

// Get returns the values for the passed key. May be nil.
func (vs *Values) Get(key string) []string {
	return vs.data[key]
}

// AssertKeyExists tests if the values contain the passed key.
func (vs *Values) AssertKeyExists(key string, msgs ...string) {
	restore := vs.wa.assert.IncrCallstackOffset()
	defer restore()
	_, ok := vs.data[key]
	vs.wa.assert.True(ok, msgs...)
}

// AssertKeyContainsValue tests if the values contain the passed key
// and that the passed value.
func (vs *Values) AssertKeyContainsValue(key, expected string, msgs ...string) {
	restore := vs.wa.assert.IncrCallstackOffset()
	defer restore()
	kd, ok := vs.data[key]
	vs.wa.assert.True(ok, msgs...)
	vs.wa.assert.Contents(expected, kd, msgs...)
}

// AssertKeyValueEquals tests if the first value for a key equals the expected value.
func (vs *Values) AssertKeyValueEquals(key, expected string, msgs ...string) {
	restore := vs.wa.assert.IncrCallstackOffset()
	defer restore()
	values, ok := vs.data[key]
	vs.wa.assert.True(ok, msgs...)
	vs.wa.assert.NotEmpty(values, msgs...)
	vs.wa.assert.Equal(values[0], expected, msgs...)
}

// applyHeader applies its values to the HTTP request header.
func (vs *Values) applyHeader(r *http.Request) {
	for key, values := range vs.data {
		for _, value := range values {
			r.Header.Add(key, value)
		}
	}
}

// applyCookies applies its values to the HTTP request cookies.
func (vs *Values) applyCookies(r *http.Request) {
	restore := vs.wa.assert.IncrCallstackOffset()
	defer restore()
	for key, kd := range vs.data {
		vs.wa.assert.NotEmpty(kd, "cookie must not be empty")
		cookie := &http.Cookie{
			Name:  key,
			Value: kd[0],
		}
		r.AddCookie(cookie)
	}
}

//--------------------
// WEB RESPONSE
//--------------------

// WebResponse provides simplified access to a response in context of
// a web asserter.
type WebResponse struct {
}

//--------------------
// WEB REQUEST
//--------------------

// WebRequest provides simplified access to a request in context of
// a web asserter.
type WebRequest struct {
	wa     *WebAsserter
	path   string
	header *Values
}

// Header
func (wr *WebRequest) Header() *Values {
	if wr.header == nil {
		wr.header = newValues(wr.wa)
	}
	return wr.header
}

// Do performes the web request with the passed method.
func (wr *WebRequest) Do(method string) *WebResponse {
	return &WebResponse{}
}

//--------------------
// WEB ASSERTER
//--------------------

// WebMultiplexer functions shall analyse requests and return the ID of
// the handler registered at the WebAsserter where to map the request to.
type WebMultiplexer func(r *http.Request) (string, error)

// WebAsserter defines the test server with methods for requests
// and uploads.
type WebAsserter struct {
	assert    *asserts.Asserts
	server    *httptest.Server
	registry  map[string]http.Handler
	multiplex WebMultiplexer
}

// NewWebAsserter creates a web test server for the tests of own handler
// or the mocking of external systems.
func NewWebAsserter(assert *asserts.Asserts, mux WebMultiplexer) *WebAsserter {
	wa := &WebAsserter{
		assert:    assert,
		multiplex: mux,
		registry:  make(map[string]http.Handler),
	}
	wa.server = httptest.NewServer(http.HandlerFunc(wa.dispatch))
	return wa
}

// Register assigns a http.HandlerFunc to an ID. That ID has to be returned by
// the mapper to address the function.
func (wa *WebAsserter) Register(id string, handler http.Handler) {
	wa.registry[id] = handler
}

// URL returns the local URL of the internal test server.
func (wa *WebAsserter) URL() string {
	return wa.server.URL
}

// Close shuts down the internal test server and blocks until all
// outstanding requests have completed.
func (wa *WebAsserter) Close() {
	wa.server.Close()
}

// CreateRequest prepares a web request to be performed
// against this web asserter.
func (wa *WebAsserter) CreateRequest(path string) *WebRequest {
	return &WebRequest{
		wa:   wa,
		path: path,
	}
}

// dispatch is the handler of the internal test server. It uses
// the request mapper function to retrieve the ID of the handler
// to use and passes response writer and request to those.
func (wa *WebAsserter) dispatch(rw http.ResponseWriter, r *http.Request) {
	id, err := wa.multiplex(r)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	handler, ok := wa.registry[id]
	if !ok {
		http.Error(rw, "mapper returned invalid handler ID: "+id, http.StatusInternalServerError)
		return
	}
	handler.ServeHTTP(rw, r)
}

// consumeHeader consumes its values from the HTTP response header.
func (wa *WebAsserter) consumeHeader(r *http.Response) *Values {
	vs := newValues(wa)
	for key, values := range r.Header {
		for _, value := range values {
			vs.Add(key, value)
		}
	}
	return vs
}

// consumeCookies consumes its values from the HTTP response cookies.
func (wa *WebAsserter) consumeCookies(r *http.Response) *Values {
	vs := newValues(wa)
	for _, cookie := range r.Cookies() {
		vs.Add(cookie.Name, cookie.Value)
	}
	return vs
}

// EOF
