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

	"tideland.dev/go/trace/errors"
)

//--------------------
// POLL
//--------------------

// Poll checks the condition until it returns true or an error. The ticker
// sends signals whenever the condition shall be checked. It closes the returned
// channel when the polling shall stop with a timeout.
func Poll(ctx context.Context, ticker Ticker, condition Condition) error {
	tickCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tickc := ticker(tickCtx)
	for {
		select {
		case <-ctx.Done():
			// Context is cancelled.
			return ctx.Err()
		case _, open := <-tickc:
			// Ticker sent a signal to check for condition.
			if !open {
				// Oh, ticker tells to end.
				return errors.New(ErrTickerExceeded, msgTickerExceeded)
			}
			ok, err := condition()
			if err != nil {
				// Condition has an error.
				return err
			}
			if ok {
				// Condition is happy.
				return nil
			}
		}
	}
}

// EOF
