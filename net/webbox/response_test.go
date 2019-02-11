// Tideland Go Library - Network - Web Toolbox - Unit Tests
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package webbox_test

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"testing"

	"tideland.one/go/audit/asserts"
	"tideland.one/go/net/webbox"
)

//--------------------
// TESTS
//--------------------

// TestMarshalResponseBody tests the encoding of data into a response body.
func TestMarshalResponseBody(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	wa := startWebAsserter(assert)
	defer wa.Close()

	wa.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var d data
		err := webbox.UnmarshalRequestBody(r, &d)
		assert.NoError(err)
		d.Number = 4321
		switch {
		case webbox.HasContentType(r, webbox.ContentTypeJSON):
			d.Name = "JSON"
			err = webbox.MarshalResponseBody(w, webbox.ContentTypeJSON, d)
			assert.NoError(err)
		case webbox.HasContentType(r, webbox.ContentTypeXML):
			d.Name = "XML"
			err = webbox.MarshalResponseBody(w, webbox.ContentTypeXML, d)
			assert.NoError(err)
		}
		w.WriteHeader(http.StatusOK)
	})

	dSend := data{
		Number: 1234,
		Name:   "Test",
		Tags:   []string{"json", "xml", "testing"},
	}

	// First run: JSON.
	req := wa.CreateRequest(http.MethodPost, "/")
	req.SetContentType(webbox.ContentTypeJSON)
	req.AssertMarshalBody(dSend)
	resp := req.Do()

	var dRecv data

	resp.AssertStatusCodeEquals(http.StatusOK)
	resp.AssertUnmarshalledBody(&dRecv)

	assert.Equal(dRecv.Number, 4321)
	assert.Equal(dRecv.Name, "JSON")
	assert.Equal(dRecv.Tags, dSend.Tags)

	// Second run: XML.
	req = wa.CreateRequest(http.MethodPost, "/")
	req.SetContentType(webbox.ContentTypeXML)
	req.AssertMarshalBody(dSend)
	resp = req.Do()

	dRecv = data{}

	resp.AssertStatusCodeEquals(http.StatusOK)
	resp.AssertUnmarshalledBody(&dRecv)

	assert.Equal(dRecv.Number, 4321)
	assert.Equal(dRecv.Name, "XML")
	assert.Equal(dRecv.Tags, dSend.Tags)
}


// EOF