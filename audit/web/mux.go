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

	"tideland.one/go/audit/asserts"
)

//--------------------
// MULTIPLEXER
//--------------------

// MultiplexMapper functions shall analyse requests and return the ID of
// the handler where to map the request to.
type MultiplexMapper func(req *http.Request) (string, error)

// Multiplexer may be passed to the test server to multiplex web requests
// to individual test handlers, e.g. stubbing the functionality of integrated
// external systems.
type Multiplexer struct {
	assert     *asserts.Asserts
	mapRequest MultiplexMapper
	registry   map[string]http.HandlerFunc
}

// NewMultiplexer creates a new multiplexer with the passed mapper.
func NewMultiplexer(assert *asserts.Asserts, mapper MultiplexMapper) *Multiplexer {
	return &Multiplexer{
		assert:     assert,
		mapRequest: mapper,
		registry:   make(map[string]http.HandlerFunc),
	}
}

// ServerHTTP implements http.Handler.
func (mux *Multiplexer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
}

// EOF
