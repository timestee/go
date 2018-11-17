// Tideland Go Library - Together - Actor
//
// Copyright (C) 2017-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package actor

//--------------------
// ERRORS
//--------------------

// Errors of the actor package.
var (
	errTimeout = "synchronous action execution timed out"
	errStopped = "actor has been stopped"
)

// IsTimedOut checks if the error signals an action timeout.
func IsTimedOut(err error) bool {
	return err.Error() == errTimeout
}

// HasBeenStopped checks if the error signals a stopped actor.
func HasBeenStopped(err error) bool {
	return err.Error() == errStopped
}

// EOF
