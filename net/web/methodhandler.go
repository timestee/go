// Tideland Go Library - Network - Web
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package web // import "tideland.dev/go/net/web"

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
)

//--------------------
// CONSTANTS
//--------------------

const (
	MethodAll = "*"
)

//--------------------
// METHOD HANDLER
//--------------------

// MethodHandler distributes request depending on the HTTP method
// to subhandlers.
type MethodHandler struct {
	handlers map[string]http.Handler
}

// NewMethodHandler creates an empty HTTP method handler.
func NewMethodHandler() *MethodHandler {
	return &MethodHandler{
		handlers: make(map[string]http.Handler),
	}
}

// Handle adds the handler based on the method.
func (mh *MethodHandler) Handle(method string, handler http.Handler) {
	mh.handlers[method] = handler
}

// HandleFunc adds the handler function based on the method.
func (mh *MethodHandler) HandleFunc(method string, hf func(http.ResponseWriter, *http.Request)) {
	mh.handlers[method] = http.HandlerFunc(hf)
}

// ServeHTTP implements http.Handler.
func (mh *MethodHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := mh.handlers[r.Method]
	if !ok {
		handler, ok = mh.handlers[MethodAll]
		if !ok {
			http.Error(w, "cannot handle request", http.StatusMethodNotAllowed)
			return
		}
	}
	handler.ServeHTTP(w, r)
}

// EOF
