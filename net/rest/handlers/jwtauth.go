// Tideland Go Library - Network - REST - Handlers
//
// Copyright (C) 2009-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package handlers

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"time"

	"tideland.one/go/net/jwt/cache"
	"tideland.one/go/net/jwt/token"
	"tideland.one/go/net/rest/core"
)

//--------------------
// JWT AUTHORIZATION HANDLER
//--------------------

// JWTAuthorizationConfig allows to control how the JWT authorization
// handler works. All values are optional. In this case tokens are only
// decoded without using a cache, validated for the current time plus/minus
// a minute leeway, and there's no user defined gatekeeper function
// running afterwards. In case of a denial a warning is written with
// the standard logger.
type JWTAuthorizationConfig struct {
	Cache      *cache.Cache
	Key        token.Key
	Leeway     time.Duration
	Gatekeeper func(job *core.Job, claims token.Claims) error
	Logger     func(job *core.Job, msg string)
}

// jwtAuthorizationHandler checks for a valid token and then runs
// a gatekeeper function. If everythinh is fine the token is stored
// in the job context for the following handlers.
type jwtAuthorizationHandler struct {
	id         string
	cache      *cache.Cache
	key        token.Key
	leeway     time.Duration
	gatekeeper func(job *core.Job, claims token.Claims) error
	logger     func(job *core.Job, msg string)
}

// NewJWTAuthorizationHandler creates a handler checking for a valid JSON
// Web Token in each request.
func NewJWTAuthorizationHandler(id string, config *JWTAuthorizationConfig) core.ResourceHandler {
	h := &jwtAuthorizationHandler{
		id:     id,
		leeway: time.Minute,
		logger: func(job *core.Job, msg string) {
			job.Environment().Log().Warningf("access denied for %v: %s", job, msg)
		},
	}
	if config != nil {
		if config.Cache != nil {
			h.cache = config.Cache
		}
		if config.Key != nil {
			h.key = config.Key
		}
		if config.Leeway != 0 {
			h.leeway = config.Leeway
		}
		if config.Gatekeeper != nil {
			h.gatekeeper = config.Gatekeeper
		}
		if config.Logger != nil {
			h.logger = config.Logger
		}
	}
	return h
}

// ID is specified on the ResourceHandler interface.
func (h *jwtAuthorizationHandler) ID() string {
	return h.id
}

// Init is specified on the ResourceHandler interface.
func (h *jwtAuthorizationHandler) Init(env *core.Environment, domain, resource string) error {
	return nil
}

// Get is specified on the GetResourceHandler interface.
func (h *jwtAuthorizationHandler) Get(job *core.Job) (bool, error) {
	return h.check(job)
}

// Head is specified on the HeadResourceHandler interface.
func (h *jwtAuthorizationHandler) Head(job *core.Job) (bool, error) {
	return h.check(job)
}

// Put is specified on the PutResourceHandler interface.
func (h *jwtAuthorizationHandler) Put(job *core.Job) (bool, error) {
	return h.check(job)
}

// Post is specified on the PostResourceHandler interface.
func (h *jwtAuthorizationHandler) Post(job *core.Job) (bool, error) {
	return h.check(job)
}

// Patch is specified on the PatchResourceHandler interface.
func (h *jwtAuthorizationHandler) Patch(job *core.Job) (bool, error) {
	return h.check(job)
}

// Delete is specified on the DeleteResourceHandler interface.
func (h *jwtAuthorizationHandler) Delete(job *core.Job) (bool, error) {
	return h.check(job)
}

// Options is specified on the OptionsResourceHandler interface.
func (h *jwtAuthorizationHandler) Options(job *core.Job) (bool, error) {
	return h.check(job)
}

// check is used by all methods to check the token.
func (h *jwtAuthorizationHandler) check(job *core.Job) (bool, error) {
	var jwt *token.JWT
	var err error
	switch {
	case h.cache != nil && h.key != nil:
		jwt, err = h.cache.RequestVerify(job.Request(), h.key)
	case h.cache != nil && h.key == nil:
		jwt, err = h.cache.RequestDecode(job.Request())
	case h.cache == nil && h.key != nil:
		jwt, err = token.RequestVerify(job.Request(), h.key)
	default:
		jwt, err = token.RequestDecode(job.Request())
	}
	// Now do the checks.
	if err != nil {
		return h.deny(job, core.StatusUnauthorized, err.Error())
	}
	if jwt == nil {
		return h.deny(job, core.StatusUnauthorized, "no JSON Web Token")
	}
	if !jwt.IsValid(h.leeway) {
		return h.deny(job, core.StatusForbidden, "the JSON Web Token claims 'nbf' and/or 'exp' are not valid")
	}
	if h.gatekeeper != nil {
		err := h.gatekeeper(job, jwt.Claims())
		if err != nil {
			return h.deny(job, core.StatusUnauthorized, "access rejected by gatekeeper: "+err.Error())
		}
	}
	// All fine, store token in context.
	job.EnhanceContext(func(ctx context.Context) context.Context {
		return token.NewContext(ctx, jwt)
	})
	return true, nil
}

// deny sends a negative feedback to the caller.
func (h *jwtAuthorizationHandler) deny(job *core.Job, statusCode int, msg string) (bool, error) {
	h.logger(job, msg)
	switch {
	case job.AcceptsContentType(core.ContentTypeJSON):
		return core.NegativeFeedback(job.JSON(false), statusCode, msg)
	case job.AcceptsContentType(core.ContentTypeXML):
		return core.NegativeFeedback(job.XML(), statusCode, msg)
	default:
		job.ResponseWriter().WriteHeader(statusCode)
		job.ResponseWriter().Header().Set("Content-Type", core.ContentTypePlain)
		job.ResponseWriter().Write([]byte(msg))
		return false, nil
	}
}

// EOF
