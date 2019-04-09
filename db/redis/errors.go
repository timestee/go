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
	ErrInvalidConfiguration   = "err-configuration"
	ErrConnectionEstablishing = "err-new-connection"
	ErrConnectionBroken       = "err-connection-broken"
	ErrInvalidResponse        = "err-inv-response"
	ErrServerResponse         = "err-server-response"
	ErrTimeout                = "err-timeout"
	ErrAuthenticate           = "err-auth"
	ErrSelectDatabase         = "err-database"
	ErrUseSubscription        = "err-subscription"
	ErrInvalidType            = "err-inv-type"
	ErrInvalidKey             = "err-inv-key"
	ErrIllegalItemIndex       = "err-item-index"
	ErrIllegalItemType        = "err-item-type"
	ErrPoolLimitReached       = "err-pool-limit"
	ErrPoolClosed             = "err-pool-closed"

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
	msgIllegalItemIndex       = "item index %d is illegal for result set size %d"
	msgIllegalItemType        = "item at index %d is no %s"
	msgPoolLimitReached       = "connection pool limit (%d) reached"
	msgPoolClosed             = "connection pool closed"
)

// EOF
