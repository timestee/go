// Tideland Go Library - Network - JSON Web Token - Cache
//
// Copyright (C) 2016-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cache

//--------------------
// CONSTANTS
//--------------------

const (
	// Error codes.
	ErrNoAuthorizationHeader      = "E101"
	ErrInvalidAuthorizationHeader = "E102"

	// Error messages.
	msgNoAuthorizationHeader      = "request contains no authorization header"
	msgInvalidAuthorizationHeader = "invalid authorization header: %q"
)

// EOF
