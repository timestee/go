// Tideland Go Library - Network - Web Toolbox
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package webbox // import "tideland.dev/go/net/webbox"

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"sync"
)

//--------------------
// METHOD MULTIPLEXER
//--------------------

// MethodMux is a handler multiplexing web requests to registered
// handlers based on the HTTP method. This way e.g. different handlers
// can be used for GET or POST. The multiplexer itself can be registered
// at a http.ServeMux. In case of no registered handler for the method
// of a request a http.StatusMethodNotAllowed will be returned.
type MethodMux struct {
	mu       sync.RWMutex
	handlers map[string]http.Handler
}

// NewMethodMux creates an instance of a method multiplexer.
func NewMethodMux() *MethodMux {
	return &MethodMux{
		handlers: make(map[string]http.Handler),
	}
}

// Handle registers the handler for the given method. If a handler already
// exists for method, Handle panics.
func (mux *MethodMux) Handle(method string, handler http.Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if !ValidMethod(method) {
		panic("webbox: invalid method")
	}
	if handler == nil {
		panic("webbox: nil handler")
	}
	if _, exist := mux.handlers[method]; exist {
		panic("webbox: multiple registrations for " + method)
	}

	mux.handlers[method] = handler
}

// HandleFunc registers the handler function for the given method.
func (mux *MethodMux) HandleFunc(method string, handler func(w http.ResponseWriter, r *http.Request)) {
	if handler == nil {
		panic("webbox: nil handler")
	}
	mux.Handle(method, http.HandlerFunc(handler))
}

// ServeHTTP implements the http.Handler interface. It dispatches the
// request to a registered handler depending on the HTTP method.
func (mux *MethodMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.mu.RLock()
	handler, ok := mux.handlers[r.Method]
	mux.mu.RUnlock()
	if !ok {
		http.Error(w, "no matching method handler found", http.StatusMethodNotAllowed)
		return
	}
	handler.ServeHTTP(w, r)
}

// EOF
