// Tideland Go Library - Network - REST - Core
//
// Copyright (C) 2009-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package core

//--------------------
// CONSTANTS
//--------------------

const (
	// Error codes.
	ErrIllegalRequest           = "E001"
	ErrCannotPrepareRequest     = "E002"
	ErrHTTPRequestFailed        = "E004"
	ErrProcessingRequestContent = "E005"
	ErrReadingResponse          = "E006"
	ErrInitHandler              = "E011"
	ErrDuplicateHandler         = "E012"
	ErrNoHandler                = "E013"
	ErrNoGetHandler             = "E014"
	ErrNoHeadHandler            = "E015"
	ErrNoPutHandler             = "E016"
	ErrNoPostHandler            = "E017"
	ErrNoPatchHandler           = "E018"
	ErrNoDeleteHandler          = "E019"
	ErrNoOptionsHandler         = "E020"
	ErrMethodNotSupported       = "E031"
	ErrUploadingFile            = "E032"
	ErrInvalidContentType       = "E033"
	ErrContentNotKeyValue       = "E034"
	ErrNoCachedTemplate         = "E035"
	ErrQueryValueNotFound       = "E036"
	ErrNoServerDefined          = "E037"

	// Error messages.
	msgIllegalRequest           = "illegal request containing too many parts"
	msgCannotPrepareRequest     = "cannot prepare request"
	msgHTTPRequestFailed        = "HTTP request failed"
	msgProcessingRequestContent = "cannot process request content"
	msgReadingResponse          = "cannot read the HTTP response"
	msgInitHandler              = "error during initialization of handler %q"
	msgDuplicateHandler         = "cannot register handler %q, it is already registered"
	msgNoHandler                = "found no handler with ID %q"
	msgNoGetHandler             = "handler %q is no handler for GET requests"
	msgNoHeadHandler            = "handler %q is no handler for HEAD requests"
	msgNoPutHandler             = "handler %q is no handler for PUT requests"
	msgNoPostHandler            = "handler %q is no handler for POST requests"
	msgNoPatchHandler           = "handler %q is no handler for PATCH requests"
	msgNoDeleteHandler          = "handler %q is no handler for DELETE requests"
	msgNoOptionsHandler         = "handler %q is no handler for OPTIONS requests"
	msgMethodNotSupported       = "method %q is not supported"
	msgUploadingFile            = "uploaded file cannot be handled by %q"
	msgInvalidContentType       = "content type is not %q"
	msgContentNotKeyValue       = "content is not key/value"
	msgNoCachedTemplate         = "template %q is not cached"
	msgQueryValueNotFound       = "query value not found"
	msgNoServerDefined          = "no server for domain '%s' configured"
)

// EOF
