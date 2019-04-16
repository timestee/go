// Tideland Go Library - Together - Wait
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package wait // import "tideland.dev/go/together/wait"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/trace/errors"
)

//--------------------
// ERRORS
//--------------------

// Error codes of the wait package.
const (
	ErrTickerExceeded    = "EEXCEEDED"
	ErrContextCancelled  = "ECANCELLED"
	ErrConditionPanicked = "EPANIC"

	msgTickerExceeded    = "ticker exceeded while waiting for the condition"
	msgContextCancelled  = "context has been cancelled"
	msgConditionPanicked = "panic during condition check: %v"
)

// IsExceeded returns true of the given error represends an exceeded ticker.
func IsExceeded(err error) bool {
	return errors.IsError(err, ErrTickerExceeded)
}

// IsCancelled returns true of the given error represends a cancelled context.
func IsCancelled(err error) bool {
	return errors.IsError(err, ErrContextCancelled)
}

// IsPanicked returns true of the given error represends a panicked condition check.
func IsPanicked(err error) bool {
	return errors.IsError(err, ErrConditionPanicked)
}

// EOF
