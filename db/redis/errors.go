// Tideland Go Library - Database - Redis Client
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
	ErrNoConfiguration        = "E001"
	ErrInvalidConfiguration   = "E002"
	ErrPoolLimitReached       = "E011"
	ErrConnectionEstablishing = "E021"
	ErrConnectionBroken       = "E022"
	ErrInvalidResponse        = "E031"
	ErrServerResponse         = "E032"
	ErrAuthenticate           = "E101"
	ErrSelectDatabase         = "E102"
	ErrUseSubscription        = "E103"
	ErrInvalidType            = "E901"
	ErrInvalidKey             = "E902"
	ErrIllegalItemIndex       = "E903"
	ErrIllegalItemType        = "E904"
	ErrTimeout                = "E999"

	msgNoConfiguration        = "no configuration"
	msgInvalidConfiguration   = "invalid configuration value in field %q: %v"
	msgPoolLimitReached       = "connection pool limit (%d) reached"
	msgConnectionEstablishing = "cannot establish connection"
	msgConnectionBroken       = "cannot %s, connection is broken"
	msgInvalidResponse        = "invalid server response: %q"
	msgServerResponse         = "server responded error: %v"
	msgAuthenticate           = "cannot authenticate"
	msgSelectDatabase         = "cannot select database"
	msgUseSubscription        = "use subscription type for subscriptions"
	msgInvalidType            = "invalid type conversion of \"%v\" to %q"
	msgInvalidKey             = "invalid key %q"
	msgIllegalItemIndex       = "item index %d is illegal for result set size %d"
	msgIllegalItemType        = "item at index %d is no %s"
	msgTimeout                = "timeout waiting for response"
)

// EOF
