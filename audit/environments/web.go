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
	"bytes"
	"encoding/json"
	"encoding/xml"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
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

// consumeHeader consumes its values from the HTTP response header.
func consumeHeader(wa *WebAsserter, resp *http.Response) *Values {
	vs := newValues(wa)
	for key, values := range resp.Header {
		for _, value := range values {
			vs.Add(key, value)
		}
	}
	return vs
}

// consumeCookies consumes its values from the HTTP response cookies.
func consumeCookies(wa *WebAsserter, resp *http.Response) *Values {
	vs := newValues(wa)
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
	wa      *WebAsserter
	resp    *http.Response
	header  *Values
	cookies *Values
	body    []byte
}

//--------------------
// WEB REQUEST
//--------------------

// WebRequest provides simplified access to a request in context of
// a web asserter.
type WebRequest struct {
	wa        *WebAsserter
	method    string
	path      string
	header    *Values
	cookies   *Values
	fieldname string
	filename  string
	body      []byte
}

// Header returns a values instance for request header.
func (wreq *WebRequest) Header() *Values {
	if wreq.header == nil {
		wreq.header = newValues(wreq.wa)
	}
	return wreq.header
}

// Cookies returns a values instance for request cookies.
func (wreq *WebRequest) Cookies() *Values {
	if wreq.cookies == nil {
		wreq.cookies = newValues(wreq.wa)
	}
	return wreq.cookies
}

// SetContentType sets the header Content-Type.
func (wreq *WebRequest) SetContentType(contentType string) {
	wreq.Header().Add(HeaderContentType, contentType)
}

// SetAccept sets the header Accept.
func (wreq *WebRequest) SetAccept(contentType string) {
	wreq.Header().Set(HeaderAccept, contentType)
}

// Upload sets the request as a file upload request.
func (wreq *WebRequest) Upload(fieldname, filename, data string) {
	wreq.fieldname = fieldname
	wreq.filename = filename
	wreq.body = []byte(data)
}

// AssertMarshalBody sets the request body based on the set content type and
// the marshalled data and asserts that everything works fine.
func (wreq *WebRequest) AssertMarshalBody(data interface{}) {
	restore := wreq.wa.assert.IncrCallstackOffset()
	defer restore()
	// Marshal the passed data into the request body.
	contentType := wreq.Header().Get(HeaderContentType)
	wreq.wa.assert.NotEmpty(contentType, "content type must be set for marshalling")
	switch contentType[0] {
	case ContentTypeApplicationJSON:
		body, err := json.Marshal(data)
		wreq.wa.assert.Nil(err, "cannot marshal data to JSON")
		wreq.body = body
		wreq.Header().Add(HeaderContentType, ContentTypeApplicationJSON)
		wreq.Header().Add(HeaderAccept, ContentTypeApplicationJSON)
	case ContentTypeApplicationXML:
		body, err := xml.Marshal(data)
		wreq.wa.assert.Nil(err, "cannot marshal data to XML")
		wreq.body = body
		wreq.Header().Add(HeaderContentType, ContentTypeApplicationXML)
		wreq.Header().Add(HeaderAccept, ContentTypeApplicationXML)
	}
}

// AssertRenderTemplate renders the passed data into the template and
// assigns it to the request body. It asserts that everything works fine.
func (wreq *WebRequest) AssertRenderTemplate(templateSource string, data interface{}) {
	restore := wreq.wa.assert.IncrCallstackOffset()
	defer restore()
	// Render template.
	t, err := template.New(wreq.path).Parse(templateSource)
	wreq.wa.assert.Nil(err, "cannot parse template")
	body := &bytes.Buffer{}
	err = t.Execute(body, data)
	wreq.wa.assert.Nil(err, "cannot render template")
	wreq.body = body.Bytes()
}

// Do performes the web request with the passed method.
func (wreq *WebRequest) Do(method string) *WebResponse {
	restore := wreq.wa.assert.IncrCallstackOffset()
	defer restore()
	// First prepare it.
	var bodyReader io.Reader
	if wreq.filename != "" {
		// Upload file content.
		buffer := &bytes.Buffer{}
		writer := multipart.NewWriter(buffer)
		part, err := writer.CreateFormFile(wreq.fieldname, wreq.filename)
		wreq.wa.assert.Nil(err, "cannot create form file")
		_, err = io.WriteString(part, string(wreq.body))
		wreq.wa.assert.Nil(err, "cannot write data")
		wreq.SetContentType(writer.FormDataContentType())
		err = writer.Close()
		wreq.wa.assert.Nil(err, "cannot close multipart writer")
		wreq.method = http.MethodPost
		bodyReader = ioutil.NopCloser(buffer)
	} else if wreq.body != nil {
		// Upload body content.
		bodyReader = ioutil.NopCloser(bytes.NewBuffer(wreq.body))
	}
	req, err := http.NewRequest(wreq.method, wreq.wa.URL()+wreq.path, bodyReader)
	wreq.wa.assert.Nil(err, "cannot prepare request")
	wreq.Header().applyHeader(req)
	wreq.Cookies().applyCookies(req)
	// Create client and perform request.
	c := http.Client{
		Transport: &http.Transport{},
	}
	resp, err := c.Do(req)
	wreq.wa.assert.Nil(err, "cannot perform test request")
	// Create web response.
	wresp := &WebResponse{
		wa:      wreq.wa,
		resp:    resp,
		header:  consumeHeader(wreq.wa, resp),
		cookies: consumeCookies(wreq.wa, resp),
	}
	body, err := ioutil.ReadAll(resp.Body)
	wreq.wa.assert.Nil(err, "cannot read response")
	defer resp.Body.Close()
	wresp.body = body
	return wresp
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
func (wa *WebAsserter) CreateRequest(method, path string) *WebRequest {
	return &WebRequest{
		wa:     wa,
		method: method,
		path:   path,
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
