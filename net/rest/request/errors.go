// Tideland Go Library - Network - REST - Request
//
// Copyright (C) 2009-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package request

//--------------------
// CONSTANTS
//--------------------

const (
	// Error codes.
	ErrNoServerDefined          = "E001"
	ErrCannotPrepareRequest     = "E002"
	ErrHTTPRequestFailed        = "E003"
	ErrProcessingRequestContent = "E004"
	ErrInvalidContent           = "E005"
	ErrAnalyzingResponse        = "E006"
	ErrDecodingResponse         = "E007"
	ErrInvalidContentType       = "E008"

	// Error messages.
	msgNoServerDefined          = "no server for domain '%s' configured"
	msgCannotPrepareRequest     = "cannot prepare request"
	msgHTTPRequestFailed        = "HTTP request failed"
	msgProcessingRequestContent = "cannot process request content"
	msgInvalidContent           = "content invalid for URL encoding"
	msgAnalyzingResponse        = "cannot analyze the HTTP response"
	msgDecodingResponse         = "cannot decode the HTTP response"
	msgInvalidContentType       = "invalid content type '%s'"
)

// EOF
