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
func NewValues(assert *asserts.Asserts, values url.Values) *Values {
	vs := &Values{
		assert: assert,
		values: values,
	}
	if vs.values == nil {
		vs.values = url.Values{}
	}
	return vs
}

// Add adds or appends a value to a named field.
func (vs *Values) Add(key, value string) {
	vsvs, ok := v.values[key]
	if ok {
		vs.values[key] = append(vsvs, value)
	} else {
		vs.values[key] = []string{value}
	}
}

// Values returns the inner values for direct usage.
func (vs *Values) Values() url.Values {
	return vs.values
}

// AssertKeyExists tests if the values contain the passed key.
func (vs *Values) AssertKeyExists(key string, msgs ...string) {
	restore := v.assert.IncrCallstackOffset()
	defer restore()
	_, ok := vs.values[key]
	v.assert.True(ok, msgs...)
}

// AssertKeyContainsValue tests if the values contain the passed key
// and that the passed value.
func (vs *Values) AssertKeyContainsValue(key, expected string, msgs ...string) {
	restore := v.assert.IncrCallstackOffset()
	defer restore()
	vsvs, ok := vs.values[key]
	v.assert.True(ok, msgs...)
	v.assert.Contents(expected, vsvs, msgs...)
}

// AssertKeyValueEquals tests if the key value equals the expected value.
func (vs *Values) AssertKeyValueEquals(key, expected string, msgs ...string) {
	restore := v.assert.IncrCallstackOffset()
	defer restore()
	vsvs, ok := vs.values[key]
	v.assert.True(ok, msgs...)
	v.assert.Equal(vs.values.Get(key), expected, msgs...)
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
