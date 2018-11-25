// Tideland Go Library - Network - REST - Core
//
// Copyright (C) 2009-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package core

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"tideland.one/go/text/etc"
	"tideland.one/go/trace/errors"
)

//--------------------
// REGISTRATIONS
//--------------------

// Registration encapsulates one handler registration.
type Registration struct {
	Domain   string
	Resource string
	Handler  ResourceHandler
}

// Registrations is a number handler registratons.
type Registrations []Registration

//--------------------
// MULTIPLEXER
//--------------------

// Multiplexer implements the http.Handler interface and adds registration
// an deregistration of resource handlers.
type Multiplexer struct {
	mutex       sync.RWMutex
	environment *Environment
	mapping     *mapping
}

// NewMultiplexer creates a new HTTP multiplexer. The passed context
// will be  used if a handler requests a context from a job, the
// configuration allows to configure the multiplexer. The allowed
// parameters are
//
//     {etc
//         {basepath /}
//         {default-domain default}
//         {default-resource default}
//         {ignore-favicon true}
//     }
//
// The values shown here are the default values if the configuration
// is nil or missing these settings.
func NewMultiplexer(ctx context.Context, cfg *etc.Etc) *Multiplexer {
	return &Multiplexer{
		environment: newEnvironment(ctx, cfg),
		mapping:     newMapping(cfg),
	}
}

// Register adds a resource handler for a given domain and resource.
func (mux *Multiplexer) Register(domain, resource string, handler ResourceHandler) error {
	mux.mutex.Lock()
	defer mux.mutex.Unlock()
	err := handler.Init(mux.environment, domain, resource)
	if err != nil {
		return err
	}
	return mux.mapping.register(domain, resource, handler)
}

// RegisterAll allows to register multiple handler in one run.
func (mux *Multiplexer) RegisterAll(registrations Registrations) error {
	for _, registration := range registrations {
		err := mux.Register(registration.Domain, registration.Resource, registration.Handler)
		if err != nil {
			return err
		}
	}
	return nil
}

// RegisteredHandlers returns the ID stack of registered handlers
// for a domain and resource.
func (mux *Multiplexer) RegisteredHandlers(domain, resource string) []string {
	mux.mutex.Lock()
	defer mux.mutex.Unlock()
	return mux.mapping.registeredHandlers(domain, resource)
}

// Deregister removes one, more, or all resource handler for a
// given domain and resource.
func (mux *Multiplexer) Deregister(domain, resource string, ids ...string) {
	mux.mutex.Lock()
	defer mux.mutex.Unlock()
	mux.mapping.deregister(domain, resource, ids...)
}

// ServeHTTP implements the http.Handler interface.
func (mux *Multiplexer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.mutex.RLock()
	defer mux.mutex.RUnlock()
	job := newJob(mux.environment, r, w)
	if err := mux.mapping.handle(job); err != nil {
		mux.handleError("error handling request", job, err)
	}
}

// handleError logs an error and returns it to the user.
func (mux *Multiplexer) handleError(format string, job *Job, err error) {
	code := http.StatusInternalServerError
	msg := fmt.Sprintf(format+" %q: %v", job, err)
	job.Environment().Log().Errorf(msg)
	if errors.IsError(err, ErrMethodNotSupported) {
		code = http.StatusMethodNotAllowed
	}
	http.Error(job.ResponseWriter(), msg, code)
}

// EOF
