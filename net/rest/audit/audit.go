// Tideland Go Library - Network - REST - Audit
//
// Copyright (C) 2009-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package audit

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"

	"tideland.one/go/audit/asserts"
	"tideland.one/go/net/rest/core"
)

//--------------------
// CONSTENTS
//--------------------

// Request and response header fields and values for testing purposes.
const (
	HeaderAccept      = "Accept"
	HeaderContentType = "Content-Type"

	ApplicationGOB  = "application/vnd.tideland.gob"
	ApplicationJSON = "application/json"
	ApplicationXML  = "application/xml"
)

//--------------------
// TEST TYPES
//--------------------

// KeyValues handles keys and values for request headers and cookies.
type KeyValues map[string]string

//--------------------
// REQUEST
//--------------------

// RequestProcessor is for pre-processing HTTP requests.
type RequestProcessor func(req *http.Request) *http.Request

// Request wraps all infos for a test request.
type Request struct {
	Method           string
	Path             string
	Header           KeyValues
	Cookies          KeyValues
	Body             []byte
	RequestProcessor RequestProcessor
}

// NewRequest creates a new test request with the given method
// and path.
func NewRequest(method, path string) *Request {
	return &Request{
		Method: method,
		Path:   path,
	}
}

// AddHeader adds or overwrites a request header.
func (r *Request) AddHeader(key, value string) *Request {
	if r.Header == nil {
		r.Header = KeyValues{}
	}
	r.Header[key] = value
	return r
}

// AddCookie adds or overwrites a request header.
func (r *Request) AddCookie(key, value string) *Request {
	if r.Cookies == nil {
		r.Cookies = KeyValues{}
	}
	r.Cookies[key] = value
	return r
}

// SetRequestProcessor sets the pre-processor.
func (r *Request) SetRequestProcessor(processor RequestProcessor) *Request {
	r.RequestProcessor = processor
	return r
}

// MarshalBody sets the request body based on the type and
// the marshalled data.
func (r *Request) MarshalBody(
	assert asserts.Asserts,
	contentType string,
	data interface{},
) *Request {
	switch contentType {
	case ApplicationGOB:
		body := &bytes.Buffer{}
		enc := gob.NewEncoder(body)
		err := enc.Encode(data)
		assert.Nil(err, "cannot encode data to GOB")
		r.Body = body.Bytes()
		r.AddHeader(HeaderContentType, ApplicationGOB)
		r.AddHeader(HeaderAccept, ApplicationGOB)
	case ApplicationJSON:
		body, err := json.Marshal(data)
		assert.Nil(err, "cannot marshal data to JSON")
		r.Body = body
		r.AddHeader(HeaderContentType, ApplicationJSON)
		r.AddHeader(HeaderAccept, ApplicationJSON)
	case ApplicationXML:
		body, err := xml.Marshal(data)
		assert.Nil(err, "cannot marshal data to XML")
		r.Body = body
		r.AddHeader(HeaderContentType, ApplicationXML)
		r.AddHeader(HeaderAccept, ApplicationXML)
	}
	return r
}

// RenderTemplate renders the passed data into the template
// and assigns it to the request body. The content type
// will be set too.
func (r *Request) RenderTemplate(
	assert asserts.Asserts,
	contentType string,
	templateSource string,
	data interface{},
) *Request {
	// Render template.
	t, err := template.New(r.Path).Parse(templateSource)
	assert.Nil(err, "cannot parse template")
	body := &bytes.Buffer{}
	err = t.Execute(body, data)
	assert.Nil(err, "cannot render template")
	r.Body = body.Bytes()
	// Set content type.
	r.AddHeader(HeaderContentType, contentType)
	r.AddHeader(HeaderAccept, contentType)
	return r
}

//--------------------
// RESPONSE
//--------------------

// Response wraps all infos of a test response.
type Response struct {
	assert  asserts.Asserts
	Status  int
	Header  KeyValues
	Cookies KeyValues
	Body    []byte
}

// AssertStatusEquals checks if the status is the expected one.
func (r *Response) AssertStatusEquals(expected int) {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	r.assert.Equal(r.Status, expected, "response status differs")
}

// AssertHeader checks if a header exists and retrieves it.
func (r *Response) AssertHeader(key string) string {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	r.assert.NotEmpty(r.Header, "response contains no header")
	value, ok := r.Header[key]
	r.assert.True(ok, "header '"+key+"' not found")
	return value
}

// AssertHeaderEquals checks if a header exists and compares
// it to an expected one.
func (r *Response) AssertHeaderEquals(key, expected string) {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	value := r.AssertHeader(key)
	r.assert.Equal(value, expected, "header value is not equal to expected")
}

// AssertHeaderContains checks if a header exists and looks for
// an expected part.
func (r *Response) AssertHeaderContains(key, expected string) {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	value := r.AssertHeader(key)
	r.assert.Substring(expected, value, "header value does not contain expected")
}

// AssertCookie checks if a cookie exists and retrieves it.
func (r *Response) AssertCookie(key string) string {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	r.assert.NotEmpty(r.Cookies, "response contains no cookies")
	value, ok := r.Cookies[key]
	r.assert.True(ok, "cookie '"+key+"' not found")
	return value
}

// AssertCookieEquals checks if a cookie exists and compares
// it to an expected one.
func (r *Response) AssertCookieEquals(key, expected string) {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	value := r.AssertCookie(key)
	r.assert.Equal(value, expected, "cookie value is not equal to expected")
}

// AssertCookieContains checks if a cookie exists and looks for
// an expected part.
func (r *Response) AssertCookieContains(key, expected string) {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	value := r.AssertCookie(key)
	r.assert.Substring(expected, value, "cookie value does not contain expected")
}

// AssertUnmarshalledBody retrieves the body based on the content type
// and unmarshals it accordingly.
func (r *Response) AssertUnmarshalledBody(data interface{}) {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	contentType, ok := r.Header[HeaderContentType]
	r.assert.True(ok)
	switch contentType {
	case ApplicationGOB:
		body := bytes.NewBuffer(r.Body)
		dec := gob.NewDecoder(body)
		err := dec.Decode(data)
		r.assert.Nil(err, "cannot decode GOB body")
	case ApplicationJSON:
		err := json.Unmarshal(r.Body, data)
		r.assert.Nil(err, "cannot unmarshal JSON body")
	case ApplicationXML:
		err := xml.Unmarshal(r.Body, data)
		r.assert.Nil(err, "cannot unmarshal XML body")
	default:
		r.assert.Fail("unknown content type: " + contentType)
	}
}

// AssertUnmarshalledFeedback retrieves a rest.Feedback as body out of
// response and returns it for further tests.
func (r *Response) AssertUnmarshalledFeedback() core.Feedback {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	fb := core.Feedback{}
	r.AssertUnmarshalledBody(&fb)
	return fb
}

// AssertBodyMatches checks if the body matches a regular expression.
func (r *Response) AssertBodyMatches(pattern string) {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	ok, err := regexp.MatchString(pattern, string(r.Body))
	r.assert.Nil(err, "illegal content match pattern")
	r.assert.True(ok, "body doesn't match pattern")
}

// AssertBodyGrep greps content out of the body.
func (r *Response) AssertBodyGrep(pattern string) []string {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	expr, err := regexp.Compile(pattern)
	r.assert.Nil(err, "illegal content grep pattern")
	return expr.FindAllString(string(r.Body), -1)
}

// AssertBodyContains checks if the body contains a string.
func (r *Response) AssertBodyContains(expected string) {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	r.assert.Contents(expected, r.Body, "body doesn't contains expected")
}

//--------------------
// TEST SERVER
//--------------------

// TestServer defines the test server with methods for requests
// and uploads.
type TestServer struct {
	server *httptest.Server
	assert asserts.Asserts
}

// StartServer starts a test server using the passed handler
func StartServer(handler http.Handler, assert asserts.Asserts) *TestServer {
	return &TestServer{
		server: httptest.NewServer(handler),
		assert: assert,
	}
}

// Close shuts down the server and blocks until all outstanding
// requests have completed.
func (ts *TestServer) Close() {
	ts.server.Close()
}

// DoRequest performs a request against the test server.
func (ts *TestServer) DoRequest(req *Request) *Response {
	// First prepare it.
	transport := &http.Transport{}
	c := &http.Client{Transport: transport}
	url := ts.server.URL + req.Path
	var bodyReader io.Reader
	if req.Body != nil {
		bodyReader = ioutil.NopCloser(bytes.NewBuffer(req.Body))
	}
	httpReq, err := http.NewRequest(req.Method, url, bodyReader)
	ts.assert.Nil(err, "cannot prepare request")
	for key, value := range req.Header {
		httpReq.Header.Set(key, value)
	}
	for key, value := range req.Cookies {
		cookie := &http.Cookie{
			Name:  key,
			Value: value,
		}
		httpReq.AddCookie(cookie)
	}
	// Check if request shall be pre-processed before performed.
	if req.RequestProcessor != nil {
		httpReq = req.RequestProcessor(httpReq)
	}
	// Now do it.
	resp, err := c.Do(httpReq)
	ts.assert.Nil(err, "cannot perform test request")
	return ts.response(resp)
}

// DoUpload is a special request for uploading a file.
func (ts *TestServer) DoUpload(path, fieldname, filename, data string) *Response {
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
func (ts *TestServer) response(hr *http.Response) *Response {
	respHeader := KeyValues{}
	for key, values := range hr.Header {
		respHeader[key] = strings.Join(values, ", ")
	}
	respCookies := KeyValues{}
	for _, cookie := range hr.Cookies() {
		respCookies[cookie.Name] = cookie.Value
	}
	respBody, err := ioutil.ReadAll(hr.Body)
	ts.assert.Nil(err, "cannot read response")
	defer hr.Body.Close()
	return &Response{
		assert:  ts.assert,
		Status:  hr.StatusCode,
		Header:  respHeader,
		Cookies: respCookies,
		Body:    respBody,
	}
}

// EOF
