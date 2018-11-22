// Tideland Go Library - Network - JSON Web Token
//
// Copyright (C) 2016-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package request

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"strings"

	"tideland.one/go/net/jwt/cache"
	"tideland.one/go/net/jwt/crypto"
	"tideland.one/go/net/jwt/token"
	"tideland.one/go/trace/errors"
)

//--------------------
// CONSTANTS
//--------------------

const (
	// Error codes.
	ErrNoAuthorizationHeader      = "E001"
	ErrInvalidAuthorizationHeader = "E002"

	// Error messages.
	msgNoAuthorizationHeader      = "request contains no authorization header"
	msgInvalidAuthorizationHeader = "invalid authorization header: %q"
)

//--------------------
// REQUEST HELPERS
//--------------------

// AddToRequest adds a token as header to a request for
// usage by a client.
func AddToRequest(req *http.Request, jwt *token.JWT) *http.Request {
	req.Header.Add("Authorization", "Bearer "+jwt.String())
	return req
}

// DecodeFromRequest tries to retrieve a token from a request
// header.
func DecodeFromRequest(req *http.Request) (*token.JWT, error) {
	return decodeFromRequest(req, nil, nil)
}

// VerifyFromRequest retrieves a possible token from
// the request. The JWT then will be verified.
func VerifyFromRequest(req *http.Request, key crypto.Key) (*token.JWT, error) {
	return decodeFromRequest(req, nil, key)
}

// VerifyCachedFromRequest retrieves a possible token from the request
// and checks if it already is cached. The JWT otherwise will be
// verified and added to the cache.
func VerifyCachedFromRequest(req *http.Request, cache *cache.Cache, key crypto.Key) (*token.JWT, error) {
	return decodeFromRequest(req, cache, key)
}

//--------------------
// PRIVATE HELPERS
//--------------------

// decodeFromRequest is the generic decoder with possible
// caching and verification.
func decodeFromRequest(req *http.Request, cache *cache.Cache, key crypto.Key) (*token.JWT, error) {
	// Retrieve token from header.
	authorization := req.Header.Get("Authorization")
	if authorization == "" {
		return nil, errors.New(ErrNoAuthorizationHeader, msgNoAuthorizationHeader)
	}
	fields := strings.Fields(authorization)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, errors.New(ErrInvalidAuthorizationHeader, msgInvalidAuthorizationHeader, authorization)
	}
	// Check cache.
	if cache != nil {
		jwt, ok := cache.Get(fields[1])
		if ok {
			return jwt, nil
		}
	}
	// Decode or verify.
	var jwt *token.JWT
	var err error
	if key == nil {
		jwt, err = token.Decode(fields[1])
	} else {
		jwt, err = token.Verify(fields[1], key)
	}
	if err != nil {
		return nil, err
	}
	// Add to cache and return.
	if cache != nil {
		cache.Put(jwt)
	}
	return jwt, nil
}

// EOF
