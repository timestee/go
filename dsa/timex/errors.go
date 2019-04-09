// Tideland Go Library - DSA - Time Extensions
//
// Copyright (C) 2015-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package timex

//--------------------
// CONSTANTS
//--------------------

// Error codes of the timex package.
const (
	ErrRetriedTooLong  = "err-too-long"
	ErrRetriedTooOften = "err-too-often"

	msgRetriedTooLong  = "retried longer than %v"
	msgRetriedTooOften = "retried more than %d times"
)

// EOF
