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
	"regexp"

	"tideland.one/go/audit/asserts"
)

//--------------------
// CONSTANTS
//--------------------

const (
	HeaderAccept      = "Accept"
	HeaderContentType = "Content-Type"

	ApplicationJSON = "application/json"
	ApplicationXML  = "application/xml"
)

//--------------------
// VALUES
//--------------------

// Values wraps header, cookie, query, and form values.
type Values struct {
	assert *asserts.Asserts
	data   map[string][]string
}

// NewValues creates a new values instance.
func NewValues(assert *asserts.Asserts) *Values {
	vs := &Values{
		assert: assert,
		data:   make(map[string][]string),
	}
	return vs
}

// ConsumeHeader consumes its values from the HTTP response header.
func ConsumeHeader(assert *asserts.Asserts, resp *http.Response) *Values {
	vs := NewValues(assert)
	for key, values := range resp.Header {
		for _, value := range values {
			vs.Add(key, value)
		}
	}
	return vs
}

// ConsumeCookies consumes its values from the HTTP response cookies.
func ConsumeCookies(assert *asserts.Asserts, resp *http.Response) *Values {
	vs := NewValues(assert)
	for _, cookie := range resp.Cookies() {
		vs.Add(cookie.Name, cookie.Value)
	}
	return vs
}

// Add adds or appends a value to a named field.
func (vs *Values) Add(key, value string) {
	kd := append(vs.data[key], value)
	vs.data[key] = kd
}

// Get returns the values for the passed key. May be nil.
func (vs *Values) Get(key string) []string {
	return vs.data[key]
}

// AssertKeyExists tests if the values contain the passed key.
func (vs *Values) AssertKeyExists(key string, msgs ...string) {
	restore := vs.assert.IncrCallstackOffset()
	defer restore()
	_, ok := vs.data[key]
	vs.assert.True(ok, msgs...)
}

// AssertKeyContainsValue tests if the values contain the passed key
// and that the passed value.
func (vs *Values) AssertKeyContainsValue(key, expected string, msgs ...string) {
	restore := vs.assert.IncrCallstackOffset()
	defer restore()
	kd, ok := vs.data[key]
	vs.assert.True(ok, msgs...)
	vs.assert.Contents(expected, kd, msgs...)
}

// AssertKeyValueEquals tests if the first value for a key equals the expected value.
func (vs *Values) AssertKeyValueEquals(key, expected string, msgs ...string) {
	restore := vs.assert.IncrCallstackOffset()
	defer restore()
	values, ok := vs.data[key]
	vs.assert.True(ok, msgs...)
	vs.assert.NotEmpty(values, msgs...)
	vs.assert.Equal(values[0], expected, msgs...)
}

// applyHeader applies its values to the HTTP request header.
func (vs *Values) applyHeader(req *http.Request) {
	for key, values := range vs.data {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
}

// applyCookies applies its values to the HTTP request cookies.
func (vs *Values) applyCookies(req *http.Request) {
	restore := vs.assert.IncrCallstackOffset()
	defer restore()
	for key, kd := range vs.data {
		vs.assert.NotEmpty(kd, "cookie must not be empty")
		cookie := &http.Cookie{
			Name:  key,
			Value: kd[0],
		}
		req.AddCookie(cookie)
	}
}

//--------------------
// REQUEST
//--------------------

// RequestProcessor is for pre-processing HTTP requests.
type RequestProcessor func(req *http.Request) *http.Request

// Request provides a convenient way to create a manual request to be handled by
// the TestServer and the registered handler there.
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
		assert:  assert,
		method:  method,
		path:    path,
		header:  NewValues(assert),
		cookies: NewValues(assert),
	}
}

// AddHeader adds or appends a request header.
func (r *Request) AddHeader(key, value string) *Request {
	r.header.Add(key, value)
	return r
}

// AddCookie adds or overwrites a request header.
func (r *Request) AddCookie(key, value string) *Request {
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

// AssertMarshalBody sets the request body based on the set content type and
// the marshalled data and asserts that everything works fine.
func (r *Request) AssertMarshalBody(data interface{}) *Request {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	// Marshal the passed data into the request body.
	contentType := r.header.Get(HeaderContentType)
	r.assert.NotEmpty(contentType, "content type must be set for marshalling")
	switch contentType[0] {
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

// AssertRenderTemplate renders the passed data into the template and
// assigns it to the request body. It asserts that everything works fine.
func (r *Request) AssertRenderTemplate(templateSource string, data interface{}) *Request {
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
	header     *Values
	cookies    *Values
	body       []byte
}

// NewResponse creates a test response wrapper.
func NewResponse(assert *asserts.Asserts, statusCode int) *Response {
	return &Response{
		assert:     assert,
		statusCode: statusCode,
	}
}

// Header returns the header values of the response.
func (r *Response) Header() *Values {
	return r.header
}

// Cookies returns the cookie values of the response.
func (r *Response) Cookies() *Values {
	return r.cookies
}

// AssertStatusCodeEquals checks if the status is the expected one.
func (r *Response) AssertStatusCodeEquals(expected int) {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	r.assert.Equal(r.statusCode, expected, "response status differs")
}

// AssertUnmarshalledBody retrieves the body based on the content type
// and unmarshals it accordingly. It asserts that everything works fine.
func (r *Response) AssertUnmarshalledBody(data interface{}) {
	restore := r.assert.IncrCallstackOffset()
	defer restore()
	contentType := r.header.Get(HeaderContentType)
	r.assert.NotEmpty(contentType)
	switch contentType[0] {
	case ApplicationJSON:
		err := json.Unmarshal(r.body, data)
		r.assert.Nil(err, "cannot unmarshal JSON body")
	case ApplicationXML:
		err := xml.Unmarshal(r.body, data)
		r.assert.Nil(err, "cannot unmarshal XML body")
	default:
		r.assert.Fail("unknown content type: " + contentType[0])
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
