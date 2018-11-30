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
	"net/url"
	"strings"

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

// Values wraps the query and cookie values for easy assertion
// and access.
type Values struct {
	assert *asserts.Asserts
	values url.Values
}

// NewValues creates a new values instance.
func NewValues(assert *asserts.Asserts, values url.Values) Values {
	return Values{
		assert: assert,
		values: values,
	}
}

// Values returns the inner values for direct usage.
func (v Values) Values() url.Values {
	return v.values
}

// AssertContainsKey tests if the values contain the passed key.
func (v Values) AssertContainsKey(key string, msgs ...string) {
	restore := v.assert.IncrCallstackOffset()
	defer restore()
	_, ok := v.values[key]
	v.assert.True(ok, msgs...)
}

// AssertKeyContainsValue tests if the values contain the passed key
// and that the passed value.
func (v Values) AssertContainsKey(key, value string, msgs ...string) {
	restore := v.assert.IncrCallstackOffset()
	defer restore()
	vs, ok := v.values[key]
	v.assert.True(ok, msgs...)
	v.assert.Contents(value, vs, msgs...)
}

//--------------------
// KEY/VAUES
//--------------------

// KeyValues handles keys and values for request headers and cookies.
type KeyValues map[string]string

// AssertContainsKey tests if the key/values contains the passed key.
func (kvs KeyValues) AssertContainsKey(assert *asserts.Asserts, key string) {
	restore := assert.IncrCallstackOffset()
	defer restore()
	_, ok := kvc[key]
	assert.True(ok, "does not contain key")
}

// ResponseHeader returns the response header as KeyValues.
func ResponseHeader(resp *http.Response) KeyValues {
	kvs := KeyValues{}
	for key, values := range hr.Header {
		kvs[key] = strings.Join(values, ", ")
	}
	return kvs
}

// ResponseCookies returns the response cookies as KeyValues.
func ResponseHeader(resp *http.Response) KeyValues {
	kvs := KeyValues{}
	for _, cookie := range hr.Cookies() {
		kvs[cookie.Name] = cookie.Value
	}
	return kvs
}

// EOF
