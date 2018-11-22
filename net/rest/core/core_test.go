// Tideland Go Library - Network - REST - Core - Unit Tests
//
// Copyright (C) 2009-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package core_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"tideland.one/go/audit/asserts"
	"tideland.one/go/dsa/version"
	"tideland.one/go/net/rest/audit"
	"tideland.one/go/net/rest/core"
	"tideland.one/go/text/etc"
)

//--------------------
// TESTS
//--------------------

// TestGetJSON tests the GET command with a JSON result.
func TestGetJSON(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "json", NewTestHandler("json", assert))
	assert.Nil(err)
	// Perform test requests.
	req := audit.NewRequest("GET", "/base/test/json/4711?foo=0815")
	req.AddHeader(audit.HeaderAccept, audit.ApplicationJSON)
	resp := ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	resp.AssertBodyContains(`"ResourceID":"4711"`)
	resp.AssertBodyContains(`"Query":"0815"`)
	resp.AssertBodyContains(`"Context":"foo"`)
}

// TestPutJSON tests the PUT command with a JSON payload and result.
func TestPutJSON(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "json", NewTestHandler("json", assert))
	assert.Nil(err)
	// Perform test requests.
	req := audit.NewRequest("PUT", "/base/test/json/4711")
	reqData := TestRequestData{"foo", "bar", "4711", "0815", ""}
	req.MarshalBody(assert, audit.ApplicationJSON, reqData)
	resp := ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	respData := TestRequestData{}
	resp.AssertUnmarshalledBody(&respData)
	assert.Equal(respData, reqData)
}

// TestGetXML tests the GET command with an XML result.
func TestGetXML(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "xml", NewTestHandler("xml", assert))
	assert.Nil(err)
	// Perform test requests.
	req := audit.NewRequest("GET", "/base/test/xml/4711")
	req.AddHeader(audit.HeaderAccept, audit.ApplicationXML)
	resp := ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	resp.AssertBodyContains(`<ResourceID>4711</ResourceID>`)
}

// TestPutXML tests the PUT command with a XML payload and result.
func TestPutXML(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "xml", NewTestHandler("xml", assert))
	assert.Nil(err)
	// Perform test requests.
	req := audit.NewRequest("PUT", "/base/test/xml/4711")
	reqData := TestRequestData{"foo", "bar", "4711", "0815", ""}
	req.MarshalBody(assert, audit.ApplicationXML, reqData)
	resp := ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	respData := TestRequestData{}
	resp.AssertUnmarshalledBody(&respData)
	assert.Equal(respData, reqData)
}

// TestPutGOB tests the PUT command with a GOB payload and result.
func TestPutGOB(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "gob", NewTestHandler("putgob", assert))
	assert.Nil(err)
	// Perform test requests.
	req := audit.NewRequest("POST", "/base/test/gob")
	reqData := TestCounterData{"test", 4711}
	req.MarshalBody(assert, audit.ApplicationGOB, reqData)
	resp := ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	respData := TestCounterData{}
	resp.AssertUnmarshalledBody(&respData)
	assert.Equal(respData, reqData)
}

// TestLongPath tests the setting of long path tail as resource ID.
func TestLongPath(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("content", "blog", NewTestHandler("default", assert))
	assert.Nil(err)
	// Perform test requests.
	req := audit.NewRequest("GET", "/base/content/blog/2014/09/30/just-a-test")
	resp := ts.DoRequest(req)
	resp.AssertBodyContains(`Resource ID: 2014/09/30/just-a-test`)
	// Now with path elements.
	req = audit.NewRequest("GET", "/base/content/blog/2014/09/30/just-another-test")
	req.AddHeader(audit.HeaderAccept, core.ContentTypePlain)
	resp = ts.DoRequest(req)
	resp.AssertBodyContains(`0: "content" 1: "blog" 2: "2014/09/30/just-another-test" 3: "09" 4: "30" 5: "just-another-test" 6: ""`)
}

// TestFallbackDefault tests the fallback to default.
func TestFallbackDefault(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("testing", "index", NewTestHandler("default", assert))
	assert.Nil(err)
	// Perform test requests.
	req := audit.NewRequest("GET", "/base/x/y")
	resp := ts.DoRequest(req)
	resp.AssertBodyContains(`Resource: y`)
}

// TestHandlerStack tests a complete handler stack.
func TestHandlerStack(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.RegisterAll(core.Registrations{
		{"authentication", "token", NewTestHandler("auth:token", assert)},
		{"test", "stack", NewAuthHandler("stack:auth", assert)},
		{"test", "stack", NewTestHandler("stack:test", assert)},
	})
	assert.Nil(err)
	// Perform test requests.
	req := audit.NewRequest("GET", "/base/test/stack")
	resp := ts.DoRequest(req)
	resp.AssertBodyContains("Resource: token")
	resp.AssertHeaderEquals("Token", "foo")
	req = audit.NewRequest("GET", "/base/test/stack")
	req.AddHeader("token", "foo")
	resp = ts.DoRequest(req)
	resp.AssertBodyContains("Resource: stack")
	resp = ts.DoRequest(req)
	resp.AssertBodyContains("Resource: stack")
}

// TestVersion tests request and response version.
func TestVersion(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "json", NewTestHandler("json", assert))
	assert.Nil(err)
	// Perform test requests.
	req := audit.NewRequest("GET", "/base/test/json/4711?foo=0815")
	req.AddHeader(audit.HeaderAccept, audit.ApplicationJSON)
	resp := ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	resp.AssertHeaderEquals("Version", "1.0.0")

	req = audit.NewRequest("GET", "/base/test/json/4711?foo=0815")
	req.AddHeader(audit.HeaderAccept, audit.ApplicationJSON)
	req.AddHeader("Version", "2")
	resp = ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	resp.AssertHeaderEquals("Version", "2.0.0")

	req = audit.NewRequest("GET", "/base/test/json/4711?foo=0815")
	req.AddHeader(audit.HeaderAccept, audit.ApplicationJSON)
	req.AddHeader("Version", "3.0")
	resp = ts.DoRequest(req)
	resp.AssertStatusEquals(200)
	resp.AssertHeaderEquals("Version", "4.0.0-alpha")
}

//  TestDeregister tests the different possibilities to stop handlers.
func TestDeregister(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.RegisterAll(core.Registrations{
		{"deregister", "single", NewTestHandler("s1", assert)},
		{"deregister", "pair", NewTestHandler("p1", assert)},
		{"deregister", "pair", NewTestHandler("p2", assert)},
		{"deregister", "group", NewTestHandler("g1", assert)},
		{"deregister", "group", NewTestHandler("g2", assert)},
		{"deregister", "group", NewTestHandler("g3", assert)},
		{"deregister", "group", NewTestHandler("g4", assert)},
		{"deregister", "group", NewTestHandler("g5", assert)},
		{"deregister", "group", NewTestHandler("g6", assert)},
	})
	assert.Nil(err)
	// Perform tests.
	assert.Equal(mux.RegisteredHandlers("deregister", "single"), []string{"s1"})
	assert.Equal(mux.RegisteredHandlers("deregister", "pair"), []string{"p1", "p2"})
	assert.Equal(mux.RegisteredHandlers("deregister", "group"), []string{"g1", "g2", "g3", "g4", "g5", "g6"})

	mux.Deregister("deregister", "single", "s1")
	assert.Nil(mux.RegisteredHandlers("deregister", "single"))
	mux.Deregister("deregister", "single")
	assert.Nil(mux.RegisteredHandlers("deregister", "single"))

	mux.Deregister("deregister", "pair")
	assert.Nil(mux.RegisteredHandlers("deregister", "pair"))

	mux.Deregister("deregister", "group", "x99")
	assert.Equal(mux.RegisteredHandlers("deregister", "group"), []string{"g1", "g2", "g3", "g4", "g5", "g6"})
	mux.Deregister("deregister", "group", "g5")
	assert.Equal(mux.RegisteredHandlers("deregister", "group"), []string{"g1", "g2", "g3", "g4", "g6"})
	mux.Deregister("deregister", "group", "g4", "g2")
	assert.Equal(mux.RegisteredHandlers("deregister", "group"), []string{"g1", "g3", "g6"})
	mux.Deregister("deregister", "group")
	assert.Nil(mux.RegisteredHandlers("deregister", "group"))
}

// TestMethodNotSupported tests the handling of a not support HTTP method.
func TestMethodNotSupported(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "no-options", NewTestHandler("no-options", assert))
	assert.Nil(err)
	// Perform test requests.
	req := audit.NewRequest("OPTIONS", "/base/test/no-options")
	resp := ts.DoRequest(req)
	resp.AssertStatusEquals(http.StatusMethodNotAllowed)
	resp.AssertBodyContains("no-options")
}

// TestRESTHandler tests the mapping of requests to the REST methods
// of a handler.
func TestRESTHandler(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "rest", NewRESTHandler("rest", assert))
	assert.Nil(err)
	err = mux.Register("test", "double", NewDoubleHandler("double", assert))
	assert.Nil(err)
	// Perform test requests on rest handler.
	req := audit.NewRequest("POST", "/base/test/rest")
	resp := ts.DoRequest(req)
	resp.AssertBodyContains("CREATE test/rest")
	req = audit.NewRequest("GET", "/base/test/rest/12345")
	resp = ts.DoRequest(req)
	resp.AssertBodyContains("READ test/rest/12345")
	req = audit.NewRequest("PUT", "/base/test/rest/12345")
	resp = ts.DoRequest(req)
	resp.AssertBodyContains("UPDATE test/rest/12345")
	req = audit.NewRequest("PATCH", "/base/test/rest/12345")
	resp = ts.DoRequest(req)
	resp.AssertBodyContains("MODIFY test/rest/12345")
	req = audit.NewRequest("DELETE", "/base/test/rest/12345")
	resp = ts.DoRequest(req)
	resp.AssertBodyContains("DELETE test/rest/12345")
	req = audit.NewRequest("OPTIONS", "/base/test/rest/12345")
	resp = ts.DoRequest(req)
	resp.AssertBodyContains("INFO test/rest/12345")
	// Perform test requests on double handler.
	req = audit.NewRequest("GET", "/base/test/double/12345")
	resp = ts.DoRequest(req)
	resp.AssertBodyContains("GET test/double/12345")
}

//--------------------
// AUTHENTICATION HANDLER
//--------------------

type AuthHandler struct {
	id     string
	assert asserts.Asserts
}

func NewAuthHandler(id string, assert asserts.Asserts) core.ResourceHandler {
	return &AuthHandler{id, assert}
}

func (ah *AuthHandler) ID() string {
	return ah.id
}

func (ah *AuthHandler) Init(env *core.Environment, domain, resource string) error {
	return nil
}

func (ah *AuthHandler) Get(job *core.Job) (bool, error) {
	token := job.Request().Header.Get("Token")
	if token != "foo" {
		job.Redirect("authentication", "token", "")
		return false, nil
	}
	job.EnhanceContext(func(ctx context.Context) context.Context {
		return context.WithValue(ctx, "Token", "foo")
	})
	return true, nil
}

//--------------------
// TEST HANDLER
//--------------------

type TestRequestData struct {
	Domain     string
	Resource   string
	ResourceID string
	Query      string
	Context    string
}

type TestCounterData struct {
	ID    string
	Count int64
}

type TestErrorData struct {
	Error string
}

const testTemplateHTML = `
<?DOCTYPE html?>
<html>
<head><title>Test</title></head>
<body>
<ul>
<li>Domain: {{.Domain}}</li>
<li>Resource: {{.Resource}}</li>
<li>Resource ID: {{.ResourceID}}</li>
<li>Query {{.Query}}</li>
<li>Context: {{.Context}}</li>
</ul>
</body>
</html>
`

type testHandler struct {
	id     string
	assert asserts.Asserts
}

func NewTestHandler(id string, assert asserts.Asserts) core.ResourceHandler {
	return &testHandler{id, assert}
}

func (th *testHandler) ID() string {
	return th.id
}

func (th *testHandler) Init(env *core.Environment, domain, resource string) error {
	env.TemplatesCache().Parse("test:context:html", testTemplateHTML, "text/html")
	return nil
}

func (th *testHandler) Get(job *core.Job) (bool, error) {
	if th.id == "auth:token" {
		job.ResponseWriter().Header().Add("Token", "foo")
	}
	if th.id == "stack:test" {
		ctxToken := job.Context().Value("Token")
		th.assert.Equal(ctxToken, "foo")
	}
	ctxTest := job.Context().Value("test")
	query := job.Query().ValueAsString("foo", "bar")
	precedence, _ := job.Version().Compare(version.New(3, 0, 0))
	// Create response.
	data := TestRequestData{job.Path().Domain(), job.Path().Resource(), job.Path().ResourceID(), query, ctxTest.(string)}
	if precedence == version.Equal {
		job.SetVersion(version.New(4, 0, 0, "alpha"))
	}
	switch {
	case job.AcceptsContentType(core.ContentTypeXML):
		th.assert.Logf("GET XML")
		job.XML().Write(core.StatusOK, data)
	case job.AcceptsContentType(core.ContentTypeJSON):
		th.assert.Logf("GET JSON")
		job.JSON(true).Write(core.StatusOK, data)
	case job.AcceptsContentType(core.ContentTypePlain):
		th.assert.True(job.Path().ContainsSubResourceIDs())
		p0 := job.Path().Domain()
		p1 := job.Path().Resource()
		p2 := job.Path().ResourceID()
		p3 := job.Path().Part(3)
		p4 := job.Path().Part(4)
		p5 := job.Path().Part(5)
		p6 := job.Path().Part(6)
		s := fmt.Sprintf("0: %q 1: %q 2: %q 3: %q 4: %q 5: %q 6: %q", p0, p1, p2, p3, p4, p5, p6)
		job.ResponseWriter().Write([]byte(s))
	default:
		th.assert.Logf("GET HTML")
		job.Renderer().Render("test:context:html", data)
	}
	return true, nil
}

func (th *testHandler) Head(job *core.Job) (bool, error) {
	return false, nil
}

func (th *testHandler) Put(job *core.Job) (bool, error) {
	var data TestRequestData
	switch {
	case job.HasContentType(core.ContentTypeJSON):
		err := job.JSON(true).Read(&data)
		if err != nil {
			job.JSON(true).Write(core.StatusBadRequest, TestErrorData{err.Error()})
		} else {
			job.JSON(true).Write(core.StatusOK, data)
		}
	case job.HasContentType(core.ContentTypeXML):
		err := job.XML().Read(&data)
		if err != nil {
			job.XML().Write(core.StatusBadRequest, TestErrorData{err.Error()})
		} else {
			job.XML().Write(core.StatusOK, data)
		}
	}
	return true, nil
}

func (th *testHandler) Post(job *core.Job) (bool, error) {
	var data TestCounterData
	err := job.GOB().Read(&data)
	if err != nil {
		job.GOB().Write(core.StatusBadRequest, err)
	} else {
		job.GOB().Write(core.StatusOK, data)
	}
	return true, nil
}

func (th *testHandler) Delete(job *core.Job) (bool, error) {
	return false, nil
}

//--------------------
// REST HANDLER
//--------------------

type restHandler struct {
	id     string
	assert asserts.Asserts
}

func NewRESTHandler(id string, assert asserts.Asserts) core.ResourceHandler {
	return &restHandler{id, assert}
}

func (rh *restHandler) ID() string {
	return rh.id
}

func (rh *restHandler) Init(env *core.Environment, domain, resource string) error {
	return nil
}

func (rh *restHandler) Create(job *core.Job) (bool, error) {
	s := fmt.Sprintf("CREATE %v/%v", job.Path().Domain(), job.Path().Resource())
	job.ResponseWriter().Write([]byte(s))
	return true, nil
}

func (rh *restHandler) Read(job *core.Job) (bool, error) {
	s := fmt.Sprintf("READ %v/%v/%v", job.Path().Domain(), job.Path().Resource(), job.Path().ResourceID())
	job.ResponseWriter().Write([]byte(s))
	return true, nil
}

func (rh *restHandler) Update(job *core.Job) (bool, error) {
	s := fmt.Sprintf("UPDATE %v/%v/%v", job.Path().Domain(), job.Path().Resource(), job.Path().ResourceID())
	job.ResponseWriter().Write([]byte(s))
	return true, nil
}

func (rh *restHandler) Modify(job *core.Job) (bool, error) {
	s := fmt.Sprintf("MODIFY %v/%v/%v", job.Path().Domain(), job.Path().Resource(), job.Path().ResourceID())
	job.ResponseWriter().Write([]byte(s))
	return true, nil
}

func (rh *restHandler) Delete(job *core.Job) (bool, error) {
	s := fmt.Sprintf("DELETE %v/%v/%v", job.Path().Domain(), job.Path().Resource(), job.Path().ResourceID())
	job.ResponseWriter().Write([]byte(s))
	return true, nil
}

func (rh *restHandler) Info(job *core.Job) (bool, error) {
	s := fmt.Sprintf("INFO %v/%v/%v", job.Path().Domain(), job.Path().Resource(), job.Path().ResourceID())
	job.ResponseWriter().Write([]byte(s))
	return true, nil
}

//--------------------
// DOUBLE HANDLER
//--------------------

// doubleHandler implements Get and Read. So Get should be chosen.
type doubleHandler struct {
	id     string
	assert asserts.Asserts
}

func NewDoubleHandler(id string, assert asserts.Asserts) core.ResourceHandler {
	return &doubleHandler{id, assert}
}

func (dh *doubleHandler) ID() string {
	return dh.id
}

func (dh *doubleHandler) Init(env *core.Environment, domain, resource string) error {
	return nil
}

func (dh *doubleHandler) Get(job *core.Job) (bool, error) {
	s := fmt.Sprintf("GET %v/%v/%v", job.Path().Domain(), job.Path().Resource(), job.Path().ResourceID())
	job.ResponseWriter().Write([]byte(s))
	return true, nil
}

func (dh *doubleHandler) Read(job *core.Job) (bool, error) {
	s := fmt.Sprintf("READ %v/%v/%v", job.Path().Domain(), job.Path().Resource(), job.Path().ResourceID())
	job.ResponseWriter().Write([]byte(s))
	return true, nil
}

//--------------------
// HELPERS
//--------------------

// newMultiplexer creates a new multiplexer with a testing context
// and a testing configuration.
func newMultiplexer(assert asserts.Asserts) *core.Multiplexer {
	ctx := context.WithValue(context.Background(), "test", "foo")
	cfgStr := "{etc {basepath /base/}{default-domain testing}{default-resource index}}"
	cfg, err := etc.ReadString(cfgStr)
	assert.Nil(err)
	return core.NewMultiplexer(ctx, cfg)
}

// EOF
