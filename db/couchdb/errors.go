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
	ErrStartupActionFailed = "err-startup-action"
	ErrInvalidVersion      = "err-version"
	ErrInvalidDocument     = "err-document"
	ErrNoIdentifier        = "err-identifier"
	ErrNotFound            = "err-not-found"
	ErrMarshallingDoc      = "err-marshalling"
	ErrUnmarshallingDoc    = "err-unmarshalling"
	ErrPreparingRequest    = "err-preparing"
	ErrPerformingRequest   = "err-performing"
	ErrClientRequest       = "err-client"
	ErrReadingResponseBody = "err-response"
	ErrUserNotFound        = "err-no-user"
	ErrUserExists          = "err-user-exists"

	msgStartupActionFailed = "startup action failed for version '%v'"
	msgInvalidVersion      = "CouchDB returns no or invalid version"
	msgInvalidDocument     = "document needs _id and _rev"
	msgNoIdentifier        = "document contains no identifier"
	msgNotFound            = "document with identifier '%s' not found"
	msgMarshallingDoc      = "cannot marshal into database document"
	msgUnmarshallingDoc    = "cannot unmarshal database document"
	msgPreparingRequest    = "cannot prepare request"
	msgPerformingRequest   = "cannot perform request"
	msgClientRequest       = "client request failed: status code %d, error '%s', reason '%s'"
	msgReadingResponseBody = "cannot read response body"
	msgUserNotFound        = "user not found"
	msgUserExists          = "user already exists"
)

// EOF
