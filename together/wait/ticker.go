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
// TICKER
//--------------------

// Ticker defines a function sending signals for each condition
// check when polling. The ticker can be canceled via the given
// context. Closing the returned signal channel means that the
// ticker ended.
type Ticker func(ctx context.Context) <-chan struct{}

// MakeIntervalTicker returns a ticker signalling in intervals.
func MakeIntervalTicker(interval time.Duration) Ticker {
	return func(ctx context.Context) <-chan struct{} {
		tickc := make(chan struct{})
		go func() {
			defer close(tickc)
			// Ticker for the interval.
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			// Loop sending signals.
			for {
				select {
				case <-ticker.C:
					// One interval tick. Ignore if needed.
					select {
					case tickc <- struct{}{}:
					default:
					}
				case <-ctx.Done():
					// Given context stopped.
					return
				}
			}
		}()
		return tickc
	}
}

// MakeDeadlinedIntervalTicker returns a ticker signalling in intervals
// and stopping after a deadline.
func MakeDeadlinedIntervalTicker(interval time.Duration, deadline time.Time) Ticker {
	return func(ctx context.Context) <-chan struct{} {
		tickc := make(chan struct{})
		if deadline.Before(time.Now()) {
			// Quick stop if deadline is before now.
			close(tickc)
			return tickc
		}
		// Let it tick.
		go func() {
			defer close(tickc)
			// Ticker for the interval.
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			// Context with deadline.
			waitCtx, cancel := context.WithDeadline(ctx, deadline)
			defer cancel()
			// Loop sending signals.
			for {
				select {
				case <-ticker.C:
					// One interval tick. Ignore if needed.
					select {
					case tickc <- struct{}{}:
					default:
					}
				case <-waitCtx.Done():
					// Deadline reached or given context stopped.
					return
				}
			}
		}()
		return tickc
	}
}

// MakeExpiringIntervalTicker returns a ticker signalling in intervals
// and stopping after a timeout.
func MakeExpiringIntervalTicker(interval, timeout time.Duration) Ticker {
	return MakeDeadlinedIntervalTicker(interval, time.Now().Add(timeout))
}

// EOF
