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

	"tideland.dev/go/trace/errors"
)

//--------------------
// WAITER
//--------------------

// Waiter defines a function sending signals for each condition
// check when waiting. The waiter can be canceled via the given
// context.
type Waiter func(ctx context.Context) <-chan struct{}

// MakeLimitedIntervalWaiter returns a waiter signalling in intervals
// and stopping after timeout.
func MakeLimitedIntervalWaiter(interval, timeout time.Duration) Waiter {
	return func(ctx context.Context) <-chan struct{} {
		waitc := make(chan struct{})
		go func() {
			defer close(waitc)
			// Ticker for the interval.
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			// Context with timeout if not 0.
			waitCtx := ctx
			if timeout != 0 {
				var cancel func()
				waitCtx, cancel = context.WithTimeout(ctx, timeout)
				defer cancel()
			}
			// Loop sending signals.
			for {
				select {
				case <-ticker.C:
					// One interval tick.
					select {
					case waitc <- struct{}{}:
					default:
					}
				case <-waitCtx.Done():
					// Timeout or waiter stopped.
					return
				}
			}
		}()
		return waitc
	}
}

//--------------------
// WAIT HELPER
//--------------------

// wait waits until the condition returns true or an error. The waiter controls
// e.g. interval and timeout, or only the interval. Also the context can stop
// the waiting.
func wait(ctx context.Context, waiter Waiter, condition Condition) error {
	waitCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	waitc := waiter(waitCtx)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case _, open := <-waitc:
			ok, err := condition()
			if err != nil {
				return err
			}
			if ok {
				return nil
			}
			if !open {
				return errors.New(ErrWaitTimeout, msgWaitTimeout)
			}
		}
	}
}

// EOF
