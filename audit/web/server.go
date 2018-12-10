// Tideland Go Library - Audit - Web
//
// Copyright (C) 2009-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package web

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"net/http/httptest"
)

//--------------------
// TEST SERVER
//--------------------

// TestServer defines the test server with methods for requests
// and uploads.
type TestServer struct {
	server *httptest.Server
}

// StartServer starts a test server using the passed handler
func StartServer(handler http.Handler) *TestServer {
	return &TestServer{
		server: httptest.NewServer(handler),
	}
}

// Close shuts down the server and blocks until all outstanding
// requests have completed.
func (ts *TestServer) Close() {
	ts.server.Close()
}

// URL returns the local URL of the test server.
func (ts *TestServer) URL() string {
	return ts.server.URL
}

// DoRequest performs a request against the test server.
func (ts *TestServer) DoRequest(req *Request) *Response {
	restore := req.assert.IncrCallstackOffset()
	defer restore()
	// Create request and perform it.
	httpReq := req.httpRequest(ts.server.URL)
	c := http.Client{
		Transport: &http.Transport{},
	}
	httpResp, err := c.Do(httpReq)
	req.assert.Nil(err, "cannot perform test request")
	return newResponse(req.assert, httpResp)
}

// EOF
