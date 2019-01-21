// Tideland Go Library - Network - Web Toolbox
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package webbox

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
)

//--------------------
// METHOD WRAPPER
//--------------------

// GetHandler has to be implemented by a handler for GET requests
// dispatched through the MethodWrapper.
type GetHandler interface {
	ServeGet(w http.ResponseWriter, r *http.Request)
}

// HeadHandler has to be implemented by a handler for HEAD requests
// dispatched through the MethodWrapper.
type HeadHandler interface {
	ServeHead(w http.ResponseWriter, r *http.Request)
}

// PostHandler has to be implemented by a handler for POST requests
// dispatched through the MethodWrapper.
type PostHandler interface {
	ServePost(w http.ResponseWriter, r *http.Request)
}

// PutHandler has to be implemented by a handler for PUT requests
// dispatched through the MethodWrapper.
type PutHandler interface {
	ServePut(w http.ResponseWriter, r *http.Request)
}

// PatchHandler has to be implemented by a handler for PATCH requests
// dispatched through the MethodWrapper.
type PatchHandler interface {
	ServePatch(w http.ResponseWriter, r *http.Request)
}

// DeleteHandler has to be implemented by a handler for DELETE requests
// dispatched through the MethodWrapper.
type DeleteHandler interface {
	ServeDelete(w http.ResponseWriter, r *http.Request)
}

// ConnectHandler has to be implemented by a handler for CONNECT requests
// dispatched through the MethodWrapper.
type ConnectHandler interface {
	ServeConnect(w http.ResponseWriter, r *http.Request)
}

// OptionsHandler has to be implemented by a handler for OPTIONS requests
// dispatched through the MethodWrapper.
type OptionsHandler interface {
	ServeOptions(w http.ResponseWriter, r *http.Request)
}

// TraceHandler has to be implemented by a handler for TRACE requests
// dispatched through the MethodWrapper.
type TraceHandler interface {
	ServeTrace(w http.ResponseWriter, r *http.Request)
}

// MethodWrapper takes a handler and dispatches requests based on the
// HTTP method to methods like ServeGet() or ServePost(). In case the
// handler provides no matching method the standard ServeHTTP() is
// called.
type MethodWrapper struct {
	handler http.Handler
}

// NewMethodWrapper creates a MethodWrapper instance.
func NewMethodWrapper(handler http.Handler) http.Handler {
	if handler == nil {
		panic("webbox: nil handler")
	}
	return MethodWrapper{
		handler: handler,
	}
}

// ServeHTTP implements the http.Handler interface. It checks the HTTP method
// and dispatches the request to the according handler method if possible.
func (mw MethodWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if h, ok := mw.handler.(GetHandler); ok {
			h.ServeGet(w, r)
			return
		}
	case http.MethodHead:
		if h, ok := mw.handler.(HeadHandler); ok {
			h.ServeHead(w, r)
			return
		}
	case http.MethodPost:
		if h, ok := mw.handler.(PostHandler); ok {
			h.ServePost(w, r)
			return
		}
	case http.MethodPut:
		if h, ok := mw.handler.(PutHandler); ok {
			h.ServePut(w, r)
			return
		}
	case http.MethodPatch:
		if h, ok := mw.handler.(PatchHandler); ok {
			h.ServePatch(w, r)
			return
		}
	case http.MethodDelete:
		if h, ok := mw.handler.(DeleteHandler); ok {
			h.ServeDelete(w, r)
			return
		}
	case http.MethodConnect:
		if h, ok := mw.handler.(ConnectHandler); ok {
			h.ServeConnect(w, r)
			return
		}
	case http.MethodOptions:
		if h, ok := mw.handler.(OptionsHandler); ok {
			h.ServeOptions(w, r)
			return
		}
	case http.MethodTrace:
		if h, ok := mw.handler.(TraceHandler); ok {
			h.ServeTrace(w, r)
			return
		}
	}
	mw.handler.ServeHTTP(w, r)
}

//--------------------
// NESTED MULTIPLEXER
//--------------------

// partHandler assigns parts to handlers.
type partHandler struct {
	part    string
	handler http.Handler
}

// NestedWrapper allows to put a number of handlers in a row. Every two
// parts of a path are assigned to one handler. Additionally a path prefix
// can be stored, e.g. for paths like /api/orders/{oid}/items/{iid}.
type NestedWrapper struct {
	prefix   string
	handlers []partHandler
}

// NewNestedWrapper creates a wrapper for nested handlers.
func NewNestedWrapper(prefix string) *NestedWrapper {
	// Only in case the prefix is not empty ensure a leading slash.
	if prefix != "" && prefix[0] != '/' {
		prefix = "+" + prefix
	}
	return &NestedWrapper{
		prefix:   prefix,
		handlers: []partHandler{},
	}
}

// Append adds a combination of part and handler.
func (nw *NestedWrapper) Append(part string, handler http.Handler) {
	if part == "" {
		panic("webbox: empty part")
	}
	if handler == nil {
		panic("webbox: nil handler")
	}
	nw.handlers = append(nw.handlers, partHandler{
		part:    part,
		handler: handler,
	})
}

// ServeHTTP implements the http.Handler interface. It analyzes the path
// and dispatches the request to the first or any later handler.
func (nw *NestedWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Strip prefix.
	// Check for leading slash.
	// Dispatch according to length.
}

// EOF
