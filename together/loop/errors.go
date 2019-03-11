// Tideland Go Library - Together - Loop
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package loop

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/trace/errors"
)

//--------------------
// ERRORS
//--------------------

// Error codes of the loop package.
const (
	ErrInvalidLoopOption = "E001"
	ErrLoopNotReady      = "E002"
	ErrLoopNotWorking    = "E003"
	ErrTimeout           = "E999"

	msgInvalidLoopOption = "invalid loop option: %v"
	msgLoopNotReady      = "loop not ready"
	msgLoopNotWorking    = "loop not working"
	msgTimeout           = "timeout during stopping"
)

//--------------------
// TEST FUNCTIONS
//--------------------

// IsErrLoopNotReady returns true if the error marks a not ready loop.
func IsErrLoopNotReady(err error) bool {
	return errors.IsError(err, ErrLoopNotReady)
}

// IsErrLoopNotWorking returns true if the error marks a not working loop.
func IsErrLoopNotWorking(err error) bool {
	return errors.IsError(err, ErrLoopNotWorking)
}

// EOF
