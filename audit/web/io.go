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
	"encoding/json"
	"encoding/xml"
	"html/template"
	"net/http"
	"net/url"
	"regexp"

	"tideland.one/go/audit/asserts"
)

//--------------------
// REQUEST
//--------------------

// RequestProcessor is for pre-processing HTTP requests.
type RequestProcessor func(req *http.Request) *http.Request

// Request wraps all infos for a test request.
type Request struct {
	assert           *asserts.Asserts
	method           string
	path             string
	header           *Values
	cookies          *Values
	body             []byte
	requestProcessor RequestProcessor
}

// NewRequest creates a new test request with the given method
// and path.
func NewRequest(assert *asserts.Asserts, method, path string) *Request {
	return &Request{
		assert: assert,
		method: method,
		path:   path,
	}
}

// AddHeader adds or appends a request header.
func (r *Request) AddHeader(key, value string) *Request {
	if r.header == nil {
		r.header = NewValues(r.assert, url.Values{})
	}
	r.header.Add(key, value)
	return r
}

// AddCookie adds or overwrites a request header.
func (r *Request) AddCookie(key, value string) *Request {
	if r.header == nil {
		r.cookies = NewValues(r.assert, url.Values{})
	}
	r.cookies.Add(key, value)
	return r
}

// SetContentType sets the header Content-Type.
func (r *Request) SetContentType(contentType string) *Request {
	return r.AddHeader(HeaderContentType, contentType)
}

// SetAccept sets the header Accept.
func (r *Request) SetAccept(contentType string) *Request {
	return r.AddHeader(HeaderAccept, contentType)
}

// SetRequestProcessor sets the pre-processor.
func (r *Request) SetRequestProcessor(processor RequestProcessor) *Request {
	r.requestProcessor = processor
	return r
}

// MarshalBody sets the request body based on the set content type and
// the marshalled data.
func (r *Request) MarshalBody(data interface{}) *Request {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	// Marshal the passed data into the request body.
	var contentType string
	if r.header != nil {
		contentType = r.header.Values().Get(HeaderContentType)
	}
	switch contentType {
	case ApplicationJSON:
		body, err := json.Marshal(data)
		r.assert.Nil(err, "cannot marshal data to JSON")
		r.body = body
		r.AddHeader(HeaderContentType, ApplicationJSON)
		r.AddHeader(HeaderAccept, ApplicationJSON)
	case ApplicationXML:
		body, err := xml.Marshal(data)
		r.assert.Nil(err, "cannot marshal data to XML")
		r.body = body
		r.AddHeader(HeaderContentType, ApplicationXML)
		r.AddHeader(HeaderAccept, ApplicationXML)
	}
	return r
}

// RenderTemplate renders the passed data into the template
// and assigns it to the request body. The content type
// will be set too.
func (r *Request) RenderTemplate(templateSource string, data interface{}) *Request {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	// Render template.
	t, err := template.New(r.path).Parse(templateSource)
	r.assert.Nil(err, "cannot parse template")
	body := &bytes.Buffer{}
	err = t.Execute(body, data)
	r.assert.Nil(err, "cannot render template")
	r.body = body.Bytes()
	return r
}

//--------------------
// RESPONSE
//--------------------

// Response wraps all infos of a test response.
type Response struct {
	assert     *asserts.Asserts
	statusCode int
	header     Values
	cookies    Values
	body       []byte
}

// NewResponse creates a test response wrapper.
func NewResponse(assert *asserts.Asserts, statusCode int) *Response {
	return &Response{
		assert:     assert,
		statusCode: statusCode,
	}
}

// AssertStatusCodeEquals checks if the status is the expected one.
func (r *Response) AssertStatusCodeEquals(expected int) {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	r.assert.Equal(r.statusCode, expected, "response status differs")
}

// AssertHeaderExists checks if a header exists and retrieves it.
func (r *Response) AssertHeaderExists(key string) string {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	r.assert.NotEmpty(r.header, "response contains no header")
	r.header.AssertKeyExists(key, "header '"+key+"' not found")
	return r.header.Values().Get(key)
}

// AssertHeaderEquals checks if a header exists and compares
// it to an expected one.
func (r *Response) AssertHeaderEquals(key, expected string) {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	value := r.AssertHeaderExists(key)
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

// AssertCookieExists checks if a cookie exists and retrieves it.
func (r *Response) AssertCookieExists(key string) string {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	r.assert.NotEmpty(r.cookies, "response contains no cookies")
	value, ok := r.cookies[key]
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
	contentType, ok := r.header[HeaderContentType]
	r.assert.True(ok)
	switch contentType {
	case ApplicationJSON:
		err := json.Unmarshal(r.body, data)
		r.assert.Nil(err, "cannot unmarshal JSON body")
	case ApplicationXML:
		err := xml.Unmarshal(r.body, data)
		r.assert.Nil(err, "cannot unmarshal XML body")
	default:
		r.assert.Fail("unknown content type: " + contentType)
	}
}

// AssertBodyMatches checks if the body matches a regular expression.
func (r *Response) AssertBodyMatches(pattern string) {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	ok, err := regexp.MatchString(pattern, string(r.body))
	r.assert.Nil(err, "illegal content match pattern")
	r.assert.True(ok, "body doesn't match pattern")
}

// AssertBodyGrep greps content out of the body.
func (r *Response) AssertBodyGrep(pattern string) []string {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	expr, err := regexp.Compile(pattern)
	r.assert.Nil(err, "illegal content grep pattern")
	return expr.FindAllString(string(r.body), -1)
}

// AssertBodyContains checks if the body contains a string.
func (r *Response) AssertBodyContains(expected string) {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	r.assert.Contents(expected, r.body, "body doesn't contains expected")
}

// EOF
