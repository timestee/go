// Tideland Go Library - Together - Wait
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package wait

//--------------------
// IMPORTS
//--------------------

import (
	"context"
)

//--------------------
// CONDITION
//--------------------

// Condition has to be implemented for checking the wanted condition. A positive
// condition will return true and nil, a negative false and nil. In case of errors
// during the check false and the error have to be returned. The function will
// be used by the poll functions.
type Condition func() (bool, error)

//--------------------
// WAIT
//--------------------

// Until waits until condition returns, it's only called once. A cancelled context,
// e.g. with timeout or deadline, will stop the waiting.
func Until(ctx context.Context, condition Condition) error {
	donec := make(chan error, 1)
	defer close(donec)
	go func() {
		// Wait for condition in background.
		_, err := condition()
		donec <- err
	}()
	select {
	case <-ctx.Done():
		// Cancelled or timeout.
		return ctx.Err()
	case err := <-donec:
		// Condition done, w/ or w/o error.
		return err
	}
}

//--------------------
// HELPER
//--------------------

// EOF
