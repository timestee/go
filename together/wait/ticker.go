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
// ticker ended. Sending ticks should be able to handle not
// received ones in case the condition check of the poller is
// working.
type Ticker func(ctx context.Context) <-chan struct{}

// TickChanger allows to work with changing intervals. The
// current one is the argument, the new has to be returned. In
// case the bool return value is false the ticker will stop.
type TickChanger func(interval time.Duration) (time.Duration, bool)

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

// MakeChangingIntervalTicker returns a ticker signalling in intervals.
// The interval can be changed with the tick changer, it starts with a
// zero.
func MakeChangingIntervalTicker(changer TickChanger) Ticker {
	return func(ctx context.Context) <-chan struct{} {
		tickc := make(chan struct{})
		interval := 0 * time.Millisecond
		ok := true
		go func() {
			defer close(tickc)
			// Loop sending signals.
			for {
				if interval, ok = changer(interval); !ok {
					return
				}
				select {
				case <-time.After(interval):
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

// MakeMaxIntervalTicker returns a ticker signalling in intervals. It
// stops after a maximum number of signals.
func MakeMaxIntervalTicker(interval time.Duration, max int) Ticker {
	return func(ctx context.Context) <-chan struct{} {
		tickc := make(chan struct{})
		count := 0
		go func() {
			defer close(tickc)
			// Ticker for the interval.
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			// Loop sending signals.
			for {
				count++
				if count > max {
					return
				}
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
