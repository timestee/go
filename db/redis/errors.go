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
	ErrInvalidConfiguration   = "ECONFIG"
	ErrConnectionEstablishing = "ENEWCONN"
	ErrConnectionBroken       = "ELOSTCONN"
	ErrInvalidResponse        = "EINVRESP"
	ErrServerResponse         = "ESERVER"
	ErrTimeout                = "ETIMEOUT"
	ErrAuthenticate           = "EAUTH"
	ErrSelectDatabase         = "EDBSELECT"
	ErrUseSubscription        = "ESUBSCR"
	ErrInvalidType            = "EINVTYP"
	ErrInvalidKey             = "EINVKEY"
	ErrInvlidItemIndex        = "EINVITIDX"
	ErrInvalidItemType        = "EINVITTYP"
	ErrPoolLimitReached       = "EPOOLLIM"
	ErrPoolClosed             = "EPOOLCLOSE"

	msgInvalidConfiguration   = "invalid configuration value in field %q: %v"
	msgConnectionEstablishing = "cannot establish new connection"
	msgConnectionBroken       = "cannot %s, connection is broken"
	msgInvalidResponse        = "invalid server response: %q"
	msgServerResponse         = "server responded error"
	msgTimeout                = "timeout waiting for response"
	msgAuthenticate           = "cannot authenticate"
	msgSelectDatabase         = "cannot select database"
	msgUseSubscription        = "use subscription type for subscriptions"
	msgInvalidType            = "invalid type conversion of \"%v\" to %q"
	msgInvalidKey             = "invalid key %q"
	msgInvalidItemIndex       = "invalid item index %d for result set size %d"
	msgInvalidItemType        = "item at index %d is no %s"
	msgPoolLimitReached       = "connection pool limit (%d) reached"
	msgPoolClosed             = "connection pool closed"
)

// EOF
