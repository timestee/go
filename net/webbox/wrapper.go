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
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strconv"
	"sync"
	"time"

	"tideland.dev/go/net/jwt/cache"
	"tideland.dev/go/net/jwt/token"
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

// NestedWrapper allows to put a number of handlers in a row. Every two
// parts of a path are assigned to one handler.
type NestedWrapper struct {
	mu       sync.RWMutex
	handlers []http.Handler
}

// NewNestedWrapper creates a wrapper for nested handlers.
func NewNestedWrapper() *NestedWrapper {
	return &NestedWrapper{
		handlers: []http.Handler{},
	}
}

// Append adds a handler.
func (nw *NestedWrapper) Append(handler http.Handler) {
	nw.mu.Lock()
	defer nw.mu.Unlock()

	if handler == nil {
		panic("webbox: nil handler")
	}
	nw.handlers = append(nw.handlers, handler)
}

// AppendFunc adds a handler function.
func (nw *NestedWrapper) AppendFunc(handler func(w http.ResponseWriter, r *http.Request)) {
	nw.Append(http.HandlerFunc(handler))
}

// ServeHTTP implements the http.Handler interface. It analyzes the path
// and dispatches the request to the first or any later handler.
func (nw *NestedWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fields := PathFields(r)
	fieldsLen := len(fields)
	n := 0
	if fieldsLen > 0 {
		n = (fieldsLen - 1) / 2
	}
	handler, ok := nw.handler(n)
	if !ok {
		http.Error(w, "handler not found", http.StatusNotFound)
		return
	}
	handler.ServeHTTP(w, r)
}

// handler returns the nth handler and true or nil and false.
func (nw *NestedWrapper) handler(n int) (http.Handler, bool) {
	nw.mu.RLock()
	defer nw.mu.RUnlock()
	if n < 0 {
		n = 0
	}
	if len(nw.handlers) < n+1 {
		return nil, false
	}
	return nw.handlers[n], true
}

//--------------------
// JWT WRAPPER
//--------------------

// JWTWrapperConfig allows to control how the JWT wrapper
// works. All values are optional. In this case tokens are only
// decoded without using a cache, validated for the current time
// plus/minus a minute leeway, and there's no user defined gatekeeper
// function running afterwards.
type JWTWrapperConfig struct {
	Cache      *cache.Cache
	Key        token.Key
	Leeway     time.Duration
	Gatekeeper func(w http.ResponseWriter, r *http.Request, claims token.Claims) error
}

// JWTWrapper checks for a valid token and then runs
// a gatekeeper function.
type JWTWrapper struct {
	handler    http.Handler
	cache      *cache.Cache
	key        token.Key
	leeway     time.Duration
	gatekeeper func(w http.ResponseWriter, r *http.Request, claims token.Claims) error
}

// NewJWTWrapper creates a handler checking for a valid JSON
// Web Token in each request.
func NewJWTWrapper(handler http.Handler, config *JWTWrapperConfig) *JWTWrapper {
	jw := &JWTWrapper{
		handler: handler,
		leeway:  time.Minute,
	}
	if config != nil {
		if config.Cache != nil {
			jw.cache = config.Cache
		}
		if config.Key != nil {
			jw.key = config.Key
		}
		if config.Leeway != 0 {
			jw.leeway = config.Leeway
		}
		if config.Gatekeeper != nil {
			jw.gatekeeper = config.Gatekeeper
		}
	}
	return jw
}

// ServeHTTP implements the http.Handler interface. It checks for an existing
// and valid token before calling the wrapped handler.
func (jw *JWTWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if jw.isAuthorized(w, r) {
		jw.handler.ServeHTTP(w, r)
	}
}

// isAuthorized checks the request for a valid token and if configured
// asks the gatekeepr if the request may pass.
func (jw *JWTWrapper) isAuthorized(w http.ResponseWriter, r *http.Request) bool {
	var jwt *token.JWT
	var err error
	switch {
	case jw.cache != nil && jw.key != nil:
		jwt, err = jw.cache.RequestVerify(r, jw.key)
	case jw.cache != nil && jw.key == nil:
		jwt, err = jw.cache.RequestDecode(r)
	case jw.cache == nil && jw.key != nil:
		jwt, err = token.RequestVerify(r, jw.key)
	default:
		jwt, err = token.RequestDecode(r)
	}
	// Now do the checks.
	if err != nil {
		jw.deny(w, r, err.Error(), http.StatusUnauthorized)
		return false
	}
	if jwt == nil {
		jw.deny(w, r, "no JSON Web Token", http.StatusUnauthorized)
		return false
	}
	if !jwt.IsValid(jw.leeway) {
		jw.deny(w, r, "the JSON Web Token claims 'nbf' and/or 'exp' are not valid", http.StatusForbidden)
		return false
	}
	if jw.gatekeeper != nil {
		err := jw.gatekeeper(w, r, jwt.Claims())
		if err != nil {
			jw.deny(w, r, "access rejected by gatekeeper: "+err.Error(), http.StatusUnauthorized)
			return false
		}
	}
	// All fine.
	return true
}

// deny sends a negative feedback to the caller.
func (jw *JWTWrapper) deny(w http.ResponseWriter, r *http.Request, msg string, statusCode int) {
	feedback := map[string]string{
		"statusCode": strconv.Itoa(statusCode),
		"message":    msg,
	}
	switch {
	case AcceptsContentType(r, ContentTypeJSON):
		b, _ := json.Marshal(feedback)
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", ContentTypeJSON)
		w.Write(b)
	case AcceptsContentType(r, ContentTypeXML):
		b, _ := xml.Marshal(feedback)
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", ContentTypeXML)
		w.Write(b)
	default:
		w.WriteHeader(statusCode)
		w.Header().Set("Content-Type", ContentTypePlain)
		w.Write([]byte(msg))
	}
}

// EOF
