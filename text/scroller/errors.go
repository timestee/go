// Tideland Go Library - Text - Scroller
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package scroller

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.one/gotrace/errors"
)

//--------------------
// CONSTANTS
//--------------------

// Error codes of the scroller package.
const (
	ErrNoSource      = "E001"
	ErrNoTarget      = "E002"
	ErrNegativeLines = "E003"
)

//--------------------
// TESTING
//--------------------

// IsNoSourceError returns true, if the error signals that
// no source has been passed.
func IsNoSourceError(err error) bool {
	return errors.IsError(err, ErrNoSource)
}

// IsNoTargetError returns true, if the error signals that
// no target has been passed.
func IsNoTargetError(err error) bool {
	return errors.IsError(err, ErrNoTarget)
}

// IsNegativeLinesError returns true, if the error shows the
// setting of a negative number of lines to start with.
func IsNegativeLinesError(err error) bool {
	return errors.IsError(err, ErrNegativeLines)
}

// EOF
