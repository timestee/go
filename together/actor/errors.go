// Tideland Go Library - Together - Actor
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package actor

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/trace/errors"
)

//--------------------
// ERRORS
//--------------------

// Error codes of the actor package.
const (
	ErrTimeout = "err-timeout"

	msgTimeout = "synchronous action execution timed out"
)

// IsTimedOut checks if the error signals an action timeout.
func IsTimedOut(err error) bool {
	return errors.IsError(err, ErrTimeout)
}

// EOF
