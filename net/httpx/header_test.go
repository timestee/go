// Tideland Go Library - Network - HTTP Extensions - Unit Tests
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package httpx_test // import "tideland.dev/go/net/httpx"

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/net/httpx"
)

//--------------------
// TESTS
//--------------------

// TestAcceptsContentType tests if the checking for accepted content
// types works correctly.
func TestAcceptsContentType(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	h := http.Header{}

	h.Set("Accept", "text/plain; q=0.5, text/html")

	assert.True(httpx.AcceptsContentType(h, httpx.ContentTypePlain))
	assert.True(httpx.AcceptsContentType(h, httpx.ContentTypeHTML))
	assert.False(httpx.AcceptsContentType(h, httpx.ContentTypeJSON))
}

// TestUnmarshalBody tests the retrieval of encoded data out of a body.
func TestUnmarshalBody(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	dIn := data{
		Number: 1234,
		Name:   "Test",
		Tags:   []string{"json", "xml", "testing"},
	}
	h := http.Header{}

	// First run: JSON.
	b, err := json.Marshal(dIn)
	assert.NoError(err)
	h.Set("Content-Type", "application/json; charset=ISO-8859-1")

	var dJSONOut data

	rc := ioutil.NopCloser(bytes.NewBuffer(b))
	err = httpx.UnmarshalBody(rc, h, &dJSONOut)
	assert.NoError(err)
	assert.Equal(dIn, dJSONOut)

	// Second run: XML.
	b, err = xml.Marshal(dIn)
	assert.NoError(err)
	h.Set("Content-Type", "application/xml; charset=ISO-8859-1")

	var dXMLOut data

	rc = ioutil.NopCloser(bytes.NewBuffer(b))
	err = httpx.UnmarshalBody(rc, h, &dXMLOut)
	assert.NoError(err)
	assert.Equal(dIn, dXMLOut)

	// Third run: plain text.
	dText := "This is a test!"
	b = []byte(dText)
	h.Set("Content-Type", "text/plain")

	var dTextOut string

	rc = ioutil.NopCloser(bytes.NewBuffer(b))
	err = httpx.UnmarshalBody(rc, h, &dTextOut)
	assert.NoError(err)
	assert.Equal(dText, dTextOut)

	// Fourth run: HTML.
	dHTML := "<html><head><title>Test</title></head><body><p>Hello, World!</p></body></html>"
	b = []byte(dHTML)
	h.Set("Content-Type", "text/html")

	var dHTMLOut string

	rc = ioutil.NopCloser(bytes.NewBuffer(b))
	err = httpx.UnmarshalBody(rc, h, &dHTMLOut)
	assert.NoError(err)
	assert.Equal(dHTML, dHTMLOut)

	// Final run: anything else in BASE64.

	dBASE64 := "VEVTVFZBTFVF"
	b = []byte{'T', 'E', 'S', 'T', 'V', 'A', 'L', 'U', 'E'}
	h.Set("Content-Type", "application/x-tideland-testing")

	var dBASE64Out string

	rc = ioutil.NopCloser(bytes.NewBuffer(b))
	err = httpx.UnmarshalBody(rc, h, &dBASE64Out)
	assert.NoError(err)
	assert.Equal(dBASE64, dBASE64Out)
}

//--------------------
// HELPERS
//--------------------

// data is used in marshalling tests.
type data struct {
	Number int
	Name   string
	Tags   []string
}

// EOF
