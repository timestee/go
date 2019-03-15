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
	"time"
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

// WithTimout waits until condition returns, it's only called once. Also a timeout
// or a cancelled context will stop the waiting.
func WithTimout(ctx context.Context, timeout time.Duration, condition Condition) error {
	waitCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return wait(waitCtx, condition)
}

//--------------------
// HELPER
//--------------------

func wait(ctx context.Context, condition Condition) error {
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

// EOF
