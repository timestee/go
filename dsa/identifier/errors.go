// Tideland Go Library - DSA - Identifier - Errors
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package identifier

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/trace/errors"
)

//--------------------
// CONSTANTS
//--------------------

// Error codes of the identifier package.
const (
	ErrInvalidHexLength = "EINVHEXLEN"
	ErrInvalidHexValue  = "EINVHEXVAL"
)

//--------------------
// TESTING
//--------------------

// IsInvalidHexLengthError returns true, if the error signals that
// the passed hex string for a UUID hasn't the correct size of 32.
func IsInvalidHexLengthError(err error) bool {
	return errors.IsError(err, ErrInvalidHexLength)
}

// IsInvalidHexValueError returns true, if the error signals an
// invalid hex string as input.
func IsInvalidHexValueError(err error) bool {
	return errors.IsError(err, ErrInvalidHexValue)
}

// EOF
