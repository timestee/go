// Tideland Go Library - DB - Redis Client
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package redis

//--------------------
// CONSTANTS
//--------------------

// Error codes.
const (
	ErrInvalidConfiguration   = "E001"
	ErrConnectionEstablishing = "E002"
	ErrConnectionBroken       = "E003"
	ErrInvalidResponse        = "E004"
	ErrServerResponse         = "E005"
	ErrTimeout                = "E006"
	ErrAuthenticate           = "E007"
	ErrSelectDatabase         = "E008"
	ErrUseSubscription        = "E009"
	ErrInvalidType            = "E010"
	ErrInvalidKey             = "E011"
	ErrIllegalItemIndex       = "E012"
	ErrIllegalItemType        = "E013"
	ErrPoolLimitReached       = "E101"
	ErrPoolClosed             = "E199"

	msgInvalidConfiguration   = "invalid configuration value in field %q: %v"
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
	msgPoolLimitReached       = "connection pool limit (%d) reached"
	msgPoolClosed             = "connection pool closed"
)

// EOF
