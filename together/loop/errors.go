// Tideland Go Library - Together - Loop
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package loop

//--------------------
// ERRORS
//--------------------

// Error codes of the loop package.
const (
	errInvalidLoopOption = "invalid loop option: %v"
	errLoopNotReady      = "loop not ready"
	errLoopNotWorking    = "loop not working"
	errLoopPanicked      = "loop panicked: %v"
	errHandlingFailed    = "error handling for %q failed"
	errRestartNonStopped = "cannot restart unstopped loop"
	errTimeout           = "timeout during stopping"
)

//--------------------
// TEST FUNCTIONS
//--------------------

// IsErrLoopNotReady returns true if the error marks a not ready loop.
func IsErrLoopNotReady(err error) bool {
	return err != nil && err.Error() == errLoopNotReady
}

// IsErrLoopNotWorking returns true if the error marks a not working loop.
func IsErrLoopNotWorking(err error) bool {
	return err != nil && err.Error() == errLoopNotWorking
}

// EOF
