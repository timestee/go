// Tideland Go Library - DB - Redis Client
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package redis

//--------------------
// IMPORTS
//--------------------

//--------------------
// CONSTANTS
//--------------------

// Error codes.
const (
	ErrInvalidConfiguration   = "E001"
	ErrPoolLimitReached       = "E002"
	ErrConnectionEstablishing = "E003"
	ErrConnectionBroken       = "E004"
	ErrInvalidResponse        = "E005"
	ErrServerResponse         = "E006"
	ErrTimeout                = "E007"
	ErrAuthenticate           = "E008"
	ErrSelectDatabase         = "E009"
	ErrUseSubscription        = "E010"
	ErrInvalidType            = "E011"
	ErrInvalidKey             = "E012"
	ErrIllegalItemIndex       = "E013"
	ErrIllegalItemType        = "E014"

	msgInvalidConfiguration   = "invalid configuration value in field %q: %v"
	msgPoolLimitReached       = "connection pool limit (%d) reached"
	msgConnectionEstablishing = "cannot establish connection"
	msgConnectionBroken       = "cannot %s, connection is broken"
	msgInvalidResponse        = "invalid server response: %q"
	msgServerResponse         = "server responded error"
	msgTimeout                = "timeout waiting for response"
	msgAuthenticate           = "cannot authenticate"
	msgSelectDatabase         = "cannot select database"
	msgUseSubscription        = "use subscription type for subscriptions"
	msgInvalidType            = "invalid type conversion of \"%v\" to %q"
	msgInvalidKey             = "invalid key %q"
	msgIllegalItemIndex       = "item index %d is illegal for result set size %d"
	msgIllegalItemType        = "item at index %d is no %s"
)

// EOF
