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
	ts := StartTestServer(assert)
	defer ts.Close()

	tests := []struct {
		method      string
		path        string
		statusCode  int
		contentType string
		body        string
	}{
		{
			method:      http.MethodGet,
			path:        "/hello/world",
			statusCode:  http.StatusOK,
			contentType: web.ContentTypeTextPlain,
			body:        "Hello, World!",
		}, {
			method:      http.MethodGet,
			path:        "/hello/tester",
			statusCode:  http.StatusOK,
			contentType: web.ContentTypeTextPlain,
			body:        "Hello, Tester!",
		}, {
			method:      http.MethodPost,
			path:        "/hello/postman",
			statusCode:  http.StatusOK,
			contentType: web.ContentTypeTextPlain,
			body:        "Hello, Postman!",
		}, {
			method:     http.MethodOptions,
			path:       "/path/does/not/exist",
			statusCode: http.StatusInternalServerError,
			body:       "mapper returned invalid handler ID",
		},
	}
	for i, test := range tests {
		assert.Logf("test case #%d: %s %s", i, test.method, test.path)
		req := web.NewRequest(assert, test.method, test.path)
		resp := ts.DoRequest(req)
		resp.AssertStatusCodeEquals(test.statusCode)
		if test.contentType != "" {
			resp.Header().AssertKeyValueEquals(web.HeaderContentType, test.contentType)
		}
		if test.body != "" {
			resp.AssertBodyMatches(test.body)
		}
	}
}

// TestHeaderCookies tests access to header and cookies.
func TestHeaderCookies(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	ts := StartTestServer(assert)
	defer ts.Close()

	tests := []struct {
		path   string
		header string
		cookie string
	}{
		{
			path:   "/header/cookies",
			header: "foo",
			cookie: "12345",
		}, {
			path:   "/header/cookies",
			header: "bar",
			cookie: "98765",
		},
	}
	for i, test := range tests {
		assert.Logf("test case #%d: GET %s", i, test.path)
		req := web.NewRequest(assert, http.MethodGet, test.path)
		req.AddHeader("Header-In", test.header)
		req.AddHeader("Cookie-In", test.cookie)
		resp := ts.DoRequest(req)
		resp.AssertStatusCodeEquals(http.StatusOK)
		resp.Header().AssertKeyValueEquals("Header-Out", test.header)
		resp.Cookies().AssertKeyValueEquals("Cookie-Out", test.cookie)
	}
}

//--------------------
// MUX MAPPER AND HANDLER
//--------------------

// StartTestServer initialises and starts the test server.
func StartTestServer(assert *asserts.Asserts) *web.TestServer {
	mux := web.NewMultiplexer(Mapper)
	mux.Register("get/hello/world", MakeHelloWorldHandler(assert, "World"))
	mux.Register("get/hello/tester", MakeHelloWorldHandler(assert, "Tester"))
	mux.Register("post/hello/postman", MakeHelloWorldHandler(assert, "Postman"))
	mux.Register("get/header/cookies", MakeHeaderCookiesHandler(assert))

	return web.StartServer(mux)
}

// Mapper returns the ID for the test handler to user.
func Mapper(r *http.Request) (string, error) {
	return strings.ToLower(r.Method + r.URL.Path), nil
}

// MakeHelloWorldHandler creates a "Hello, World" handler.
func MakeHelloWorldHandler(assert *asserts.Asserts, who string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reply := "Hello, " + who + "!"
		w.Header().Add(web.HeaderContentType, web.ContentTypeTextPlain)
		w.Write([]byte(reply))
		w.WriteHeader(http.StatusOK)
	}
}

// MakeHeaderCookiesHandler creates a handler for header and cookies.
func MakeHeaderCookiesHandler(assert *asserts.Asserts) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerOut := r.Header.Get("Header-In")
		cookieOut := r.Header.Get("Cookie-In")
		http.SetCookie(w, &http.Cookie{
			Name:  "Cookie-Out",
			Value: cookieOut,
		})
		w.Header().Set(web.HeaderContentType, web.ContentTypeTextPlain)
		w.Header().Set("Header-Out", headerOut)
		w.Write([]byte("Done!"))
		w.WriteHeader(http.StatusOK)
	}
}

// EOF
