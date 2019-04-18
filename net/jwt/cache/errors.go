// Tideland Go Library - Network - JSON Web Token - Cache
//
// Copyright (C) 2016-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cache // import "tideland.dev/go/net/jwt/cache"

//--------------------
// CONSTANTS
//--------------------

// Error codes and messages.
const (
	ErrNoAuthorizationHeader      = "err-no-auth"
	ErrInvalidAuthorizationHeader = "err-inv-auth"

	msgNoAuthorizationHeader      = "request contains no authorization header"
	msgInvalidAuthorizationHeader = "invalid authorization header: %q"
)

// EOF
