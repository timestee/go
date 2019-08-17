// Tideland Go Library - Network - Web Toolbox - Unit Tests
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package webbox_test // import "tideland.dev/go/net/webbox"

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"testing"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/net/webbox"
)

//--------------------
// TESTS
//--------------------

// TestAcceptsContentType tests if the checking for accepted content
// types works correctly.
func TestAcceptsContentType(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	r, err := http.NewRequest(http.MethodGet, "http://localhost/", nil)
	assert.NoError(err)
	r.Header.Set("Accept", "text/plain; q=0.5, text/html")

	assert.True(webbox.AcceptsContentType(r, webbox.ContentTypePlain))
	assert.True(webbox.AcceptsContentType(r, webbox.ContentTypeHTML))
	assert.False(webbox.AcceptsContentType(r, webbox.ContentTypeJSON))
}

// TestHasContentType tests if the checking for contained content
// types works correctly.
func TestHasContentType(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	r, err := http.NewRequest(http.MethodGet, "http://localhost/", nil)
	assert.NoError(err)
	r.Header.Set("Content-Type", "application/json; charset=ISO-8859-1")

	assert.True(webbox.HasContentType(r, webbox.ContentTypeJSON))
	assert.False(webbox.HasContentType(r, webbox.ContentTypeURLEncoded))
}

// TestUnmarshalRequestBody tests the retrieval of encoded data
// out of a request body.
func TestUnmarshalRequestBody(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	dIn := data{
		Number: 1234,
		Name:   "Test",
		Tags:   []string{"json", "xml", "testing"},
	}

	// First run: JSON.
	b, err := json.Marshal(dIn)
	assert.NoError(err)
	r, err := http.NewRequest(http.MethodPost, "http://localhost/", bytes.NewBuffer(b))
	assert.NoError(err)
	r.Header.Set("Content-Type", "application/json; charset=ISO-8859-1")

	var dJSONOut data

	err = webbox.UnmarshalRequestBody(r, &dJSONOut)
	assert.NoError(err)
	assert.Equal(dIn, dJSONOut)

	// Second run: XML.
	b, err = xml.Marshal(dIn)
	assert.NoError(err)
	r, err = http.NewRequest(http.MethodPost, "http://localhost/", bytes.NewBuffer(b))
	assert.NoError(err)
	r.Header.Set("Content-Type", "application/xml; charset=ISO-8859-1")

	var dXMLOut data

	err = webbox.UnmarshalRequestBody(r, &dXMLOut)
	assert.NoError(err)
	assert.Equal(dIn, dXMLOut)

	// Third run: plain text with error.
	b = []byte("Boom!")
	r, err = http.NewRequest(http.MethodPost, "http://localhost/", bytes.NewBuffer(b))
	assert.NoError(err)
	r.Header.Set("Content-Type", "text/plain")

	var dTextOut string

	err = webbox.UnmarshalRequestBody(r, &dTextOut)
	assert.ErrorMatch(err, "webbox: invalid content-type")
}

// EOF
