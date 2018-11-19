// Tideland Go Library - DSA - Time Extensions
//
// Copyright (C) 2015-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package timex

//--------------------
// CONSTANTS
//--------------------

// Error codes of the timex package.
const (
	ErrRetriedTooLong  = "E001"
	ErrRetriedTooOften = "E002"

	msgRetriedTooLong  = "retried longer than %v"
	msgRetriedTooOften = "retried more than %d times"
)

// EOF
