// Tideland Go Library - Network - JSON Web Token
//
// Copyright (C) 2016-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package token

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"strings"

	"tideland.dev/go/trace/errors"
)

//--------------------
// REQUEST HELPERS
//--------------------

// RequestAdd adds a token as header to a request for
// usage by a client.
func RequestAdd(req *http.Request, jwt *JWT) *http.Request {
	req.Header.Add("Authorization", "Bearer "+jwt.String())
	return req
}

// RequestDecode tries to retrieve a token from a request header.
func RequestDecode(req *http.Request) (*JWT, error) {
	return decode(req, nil)
}

// RequestVerify retrieves a possible token from a request.
// The JWT then will be verified.
func RequestVerify(req *http.Request, key Key) (*JWT, error) {
	return decode(req, key)
}

//--------------------
// PRIVATE HELPERS
//--------------------

// decodeFromRequest is the generic decoder with possible
// caching and verification.
func decode(req *http.Request, key Key) (*JWT, error) {
	// Retrieve token from header.
	authorization := req.Header.Get("Authorization")
	if authorization == "" {
		return nil, errors.New(ErrNoAuthorizationHeader, msgNoAuthorizationHeader)
	}
	fields := strings.Fields(authorization)
	if len(fields) != 2 || fields[0] != "Bearer" {
		return nil, errors.New(ErrInvalidAuthorizationHeader, msgInvalidAuthorizationHeader, authorization)
	}
	// Decode or verify.
	var jwt *JWT
	var err error
	if key == nil {
		jwt, err = Decode(fields[1])
	} else {
		jwt, err = Verify(fields[1], key)
	}
	if err != nil {
		return nil, err
	}
	return jwt, nil
}

// EOF
