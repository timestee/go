// Tideland Go Library - DB - CouchDB Client
//
// Copyright (C) 2016-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package couchdb

//--------------------
// CONSTANTS
//--------------------

// Error codes.
const (
	ErrStartupActionFailed = "E001"
	ErrInvalidVersion      = "E002"
	ErrInvalidDocument     = "E003"
	ErrNoIdentifier        = "E004"
	ErrNotFound            = "E005"
	ErrMarshallingDoc      = "E006"
	ErrPreparingRequest    = "E006"
	ErrPerformingRequest   = "E008"
	ErrClientRequest       = "E009"
	ErrUnmarshallingDoc    = "E010"
	ErrReadingResponseBody = "E012"
	ErrNoSession           = "E101"
	ErrUserNotFound        = "E102"
	ErrUserExists          = "E103"

	msgStartupActionFailed = "startup action failed for version '%v'"
	msgInvalidVersion      = "CouchDB returns no or invalid version"
	msgInvalidDocument     = "document needs _id and _rev"
	msgNoIdentifier        = "document contains no identifier"
	msgNotFound            = "document with identifier '%s' not found"
	msgMarshallingDoc      = "cannot marshal into database document"
	msgPreparingRequest    = "cannot prepare request"
	msgPerformingRequest   = "cannot perform request"
	msgClientRequest       = "client request failed: status code %d, error '%s', reason '%s'"
	msgUnmarshallingDoc    = "cannot unmarshal database document"
	msgReadingResponseBody = "cannot read response body"
	msgNoSession           = "command needs authenticated session"
	msgUserNotFound        = "user not found"
	msgUserExists          = "user already exists"
)

// EOF
