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

	"tideland.one/go/text/etc"
	"tideland.one/go/text/stringex"
	"tideland.one/go/trace/logger"
)

//--------------------
// ENVIRONMENT
//--------------------

// Environment describes the environment of a RESTful application.
type Environment struct {
	ctx             context.Context
	log             *logger.Logger
	basepath        string
	baseparts       []string
	basepartsLen    int
	defaultDomain   string
	defaultResource string
	templatesCache  *TemplatesCache
}

// newEnvironment crerates an environment using the
// passed context and configuration.
func newEnvironment(ctx context.Context, cfg etc.Etc) *Environment {
	env := &Environment{
		basepath:        "/",
		baseparts:       []string{},
		defaultDomain:   "default",
		defaultResource: "default",
		templatesCache:  newTemplatesCache(),
	}
	// Check configuration.
	if cfg != nil {
		env.basepath = cfg.ValueAsString("basepath", env.basepath)
		env.defaultDomain = cfg.ValueAsString("default-domain", env.defaultDomain)
		env.defaultResource = cfg.ValueAsString("default-resource", env.defaultResource)
	}
	// Check basepath and remove empty parts.
	env.baseparts = stringex.SplitMap(env.basepath, "/", func(p string) (string, bool) {
		if p == "" {
			return "", false
		}
		return p, true
	})
	env.basepartsLen = len(env.baseparts)
	// Set log.
	// TODO: According configuration options have to be added.
	ctx.log = logger.NewStandard(logger.NewStandardOutWriter())
	// Set context.
	if ctx == nil {
		ctx = context.Background()
	}
	env.ctx = newEnvironmentContext(ctx, env)
	return env
}

// Context returns the context of the environment.
func (env *Environment) Context() context.Context {
	return env.ctx
}

// Log returns the configured logger of the environment.
func (env *Environment) Log() *context.Logger {
	return env.log
}

// Basepath returns the configured basepath.
func (env *Environment) Basepath() string {
	return env.basepath
}

// DefaultDomain returns the configured default domain.
func (env *Environment) DefaultDomain() string {
	return env.defaultDomain
}

// DefaultResource returns the configured default resource.
func (env *Environment) DefaultResource() string {
	return env.defaultResource
}

// TemplatesCache returns the template cache.
func (env *Environment) TemplatesCache() *TemplatesCache {
	return env.templatesCache
}

// EOF
