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
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"

	"tideland.one/go/audit/asserts"
)

//--------------------
// TEST SERVER
//--------------------

// TestServer defines the test server with methods for requests
// and uploads.
type TestServer struct {
	assert *asserts.Asserts
	server *httptest.Server
}

// StartServer starts a test server using the passed handler
func StartServer(assert *asserts.Asserts, handler http.Handler) *TestServer {
	return &TestServer{
		assert: assert,
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
	restore := ts.assert.IncrCallstackOffset()
	defer restore()
	// First prepare it.
	transport := &http.Transport{}
	c := &http.Client{Transport: transport}
	url := ts.server.URL + req.path
	var bodyReader io.Reader
	if req.body != nil {
		bodyReader = ioutil.NopCloser(bytes.NewBuffer(req.body))
	}
	httpReq, err := http.NewRequest(req.method, url, bodyReader)
	ts.assert.Nil(err, "cannot prepare request")
	req.header.applyHeader(httpReq)
	req.cookies.applyCookies(httpReq)
	// Check if request shall be pre-processed before performed.
	if req.requestProcessor != nil {
		httpReq = req.requestProcessor(httpReq)
	}
	// Now do it.
	resp, err := c.Do(httpReq)
	ts.assert.Nil(err, "cannot perform test request")
	return ts.response(resp)
}

// DoUpload is a special request for uploading a file.
func (ts *TestServer) DoUpload(path, fieldname, filename, data string) *Response {
	restore := ts.assert.IncrCallstackOffset()
	defer restore()
	// Prepare request.
	transport := &http.Transport{}
	c := &http.Client{Transport: transport}
	url := ts.server.URL + path
	buffer := &bytes.Buffer{}
	writer := multipart.NewWriter(buffer)
	part, err := writer.CreateFormFile(fieldname, filename)
	ts.assert.Nil(err, "cannot create form file")
	_, err = io.WriteString(part, data)
	ts.assert.Nil(err, "cannot write data")
	contentType := writer.FormDataContentType()
	err = writer.Close()
	ts.assert.Nil(err, "cannot close multipart writer")
	// And now do it.
	resp, err := c.Post(url, contentType, buffer)
	ts.assert.Nil(err, "cannot perform test upload")
	return ts.response(resp)
}

// response creates a Response instance out of the http.Response-
func (ts *TestServer) response(httpResp *http.Response) *Response {
	header := ConsumeHeader(ts.assert, httpResp)
	cookies := ConsumeCookies(ts.assert, httpResp)
	body, err := ioutil.ReadAll(httpResp.Body)
	ts.assert.Nil(err, "cannot read response")
	defer httpResp.Body.Close()
	return &Response{
		assert:     ts.assert,
		statusCode: httpResp.StatusCode,
		header:     header,
		cookies:    cookies,
		body:       body,
	}
}

// EOF
