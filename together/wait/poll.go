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
// TICKER
//--------------------

// Ticker defines a function sending signals for each condition
// check when polling. The ticker can be canceled via the given
// context. Closing the returned signal channel means that the
// ticker ended.
type Ticker func(ctx context.Context) <-chan struct{}

//--------------------
// POLL
//--------------------

// Poll checks endlessly in intervals if the condition returns true. In case of a
// cancelled context the polling will stop.
func Poll(ctx context.Context, interval time.Duration, condition Condition) error {
	return PollWithTicker(ctx, makeIntervalTicker(interval), condition)
}

// PollWithDeadline checks in intervals if the condition returns true or the deadline is
// reached. Also in case of a cancelled context the polling will stop.
func PollWithDeadline(ctx context.Context, interval time.Duration, deadline time.Time, condition Condition) error {
	return PollWithTicker(ctx, makeIntervalDeadlineTicker(interval, deadline), condition)
}

// PollWithTimeout checks in intervals if the condition returns true or the timeout is
// reached. Also in case of a cancelled context the polling will stop.
func PollWithTimeout(ctx context.Context, interval, timeout time.Duration, condition Condition) error {
	return PollWithDeadline(ctx, interval, time.Now().Add(timeout), condition)
}

// PollWithTicker checks the condition until it returns true or an error. The ticker
// sends signals whenever the condition shall be checked. It closes the returned
// channel when the polling shall stop with a timeout.
func PollWithTicker(ctx context.Context, ticker Ticker, condition Condition) error {
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
				return errors.New(ErrWaitTimeout, msgWaitTimeout)
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

//--------------------
// HELPER
//--------------------

// makeIntervalTicker returns a ticker signalling in intervals.
func makeIntervalTicker(interval time.Duration) Ticker {
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

// makeIntervalDeadlineTicker returns a ticker signalling in intervals
// and stopping after a deadline.
func makeIntervalDeadlineTicker(interval time.Duration, deadline time.Time) Ticker {
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

// EOF
