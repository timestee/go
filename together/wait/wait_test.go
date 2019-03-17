// Tideland Go Library - Together - Wait - Unit Tests
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package wait_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/wait"
	"tideland.dev/go/trace/errors"
)

//--------------------
// TESTS
//--------------------

// TestPollWithInterval tests the polling of conditions in intervals.
func TestPollWithInterval(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)

	// Tests.
	assert.Logf("end with positive condition")
	count := 0
	err := wait.Poll(
		context.Background(),
		wait.MakeIntervalTicker(50*time.Millisecond),
		func() (bool, error) {
			count++
			if count == 5 {
				return true, nil
			}
			return false, nil
		},
	)
	assert.NoError(err)
	assert.Equal(count, 5)

	assert.Logf("end with cancelled context")
	count = 0
	ctx, cancel := context.WithTimeout(context.Background(), 110*time.Millisecond)
	defer cancel()
	err = wait.Poll(
		ctx,
		wait.MakeIntervalTicker(20*time.Millisecond),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, "context deadline exceeded")
	assert.Equal(count, 5)
}

// TestPollWithChangingInterval tests the polling of conditions in a maximum
// number of intervals.
func TestPollWithChangingInterval(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	changer := func(interval time.Duration) (time.Duration, bool) {
		if interval == 0 {
			interval = 10 * time.Millisecond
		} else {
			interval *= 2
		}
		if interval > time.Second {
			return 0, false
		}
		return interval, true
	}

	// Tests.
	assert.Logf("end with positive condition")
	count := 0
	err := wait.Poll(
		context.Background(),
		wait.MakeChangingIntervalTicker(changer),
		func() (bool, error) {
			count++
			if count == 5 {
				return true, nil
			}
			return false, nil
		},
	)
	assert.NoError(err)
	assert.Equal(count, 5)

	assert.Logf("end with deadline, 7 checks")
	count = 0
	err = wait.Poll(
		context.Background(),
		wait.MakeChangingIntervalTicker(changer),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.True(errors.IsError(err, wait.ErrTickerExceeded))
	assert.Equal(count, 7, "exceeded with a count")

	assert.Logf("end with cancelled context")
	count = 0
	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()
	err = wait.Poll(
		ctx,
		wait.MakeChangingIntervalTicker(changer),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, "context deadline exceeded")
	assert.Equal(count, 5)
}

// TestPollWithMaxInterval tests the polling of conditions in a maximum
// number of intervals.
func TestPollWithMaxInterval(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)

	// Tests.
	assert.Logf("end with positive condition")
	count := 0
	err := wait.Poll(
		context.Background(),
		wait.MakeMaxIntervalTicker(20*time.Millisecond, 10),
		func() (bool, error) {
			count++
			if count == 5 {
				return true, nil
			}
			return false, nil
		},
	)
	assert.NoError(err)
	assert.Equal(count, 5)

	assert.Logf("end with deadline, 10 checks")
	count = 0
	err = wait.Poll(
		context.Background(),
		wait.MakeMaxIntervalTicker(20*time.Millisecond, 10),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.True(errors.IsError(err, wait.ErrTickerExceeded))
	assert.Equal(count, 10, "exceeded with a count")

	assert.Logf("end with deadline, no check")
	count = 0
	err = wait.Poll(
		context.Background(),
		wait.MakeMaxIntervalTicker(20*time.Millisecond, -1),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.True(errors.IsError(err, wait.ErrTickerExceeded))
	assert.Equal(count, 0)

	assert.Logf("end with cancelled context")
	count = 0
	ctx, cancel := context.WithTimeout(context.Background(), 110*time.Millisecond)
	defer cancel()
	err = wait.Poll(
		ctx,
		wait.MakeMaxIntervalTicker(20*time.Millisecond, 10),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, "context deadline exceeded")
	assert.Equal(count, 5)
}

// TestPollWithDeadline tests the polling of conditions with deadlines.
func TestPollWithDeadline(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)

	// Tests.
	assert.Logf("end with positive condition")
	count := 0
	err := wait.Poll(
		context.Background(),
		wait.MakeDeadlinedIntervalTicker(20*time.Millisecond, time.Now().Add(210*time.Millisecond)),
		func() (bool, error) {
			count++
			if count == 5 {
				return true, nil
			}
			return false, nil
		},
	)
	assert.NoError(err)
	assert.Equal(count, 5)

	assert.Logf("end with deadline, 10 checks")
	count = 0
	err = wait.Poll(
		context.Background(),
		wait.MakeDeadlinedIntervalTicker(20*time.Millisecond, time.Now().Add(210*time.Millisecond)),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.True(errors.IsError(err, wait.ErrTickerExceeded))
	assert.Equal(count, 10, "exceeded with a count")

	assert.Logf("end with deadline, no check")
	count = 0
	err = wait.Poll(
		context.Background(),
		wait.MakeDeadlinedIntervalTicker(20*time.Millisecond, time.Now().Add(-time.Second)),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.True(errors.IsError(err, wait.ErrTickerExceeded))
	assert.Equal(count, 0)

	assert.Logf("end with cancelled context")
	count = 0
	ctx, cancel := context.WithTimeout(context.Background(), 110*time.Millisecond)
	defer cancel()
	err = wait.Poll(
		ctx,
		wait.MakeDeadlinedIntervalTicker(20*time.Millisecond, time.Now().Add(time.Second)),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, "context deadline exceeded")
	assert.Equal(count, 5)
}

// TestPollWithTimeout tests the polling of conditions with timeouts.
func TestPollWithTimeout(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)

	// Tests.
	assert.Logf("end with positive condition")
	count := 0
	err := wait.Poll(
		context.Background(),
		wait.MakeExpiringIntervalTicker(20*time.Millisecond, 210*time.Millisecond),
		func() (bool, error) {
			count++
			if count == 5 {
				return true, nil
			}
			return false, nil
		},
	)
	assert.NoError(err)
	assert.Equal(count, 5)

	assert.Logf("end with timeout, 10 checks")
	count = 0
	err = wait.Poll(
		context.Background(),
		wait.MakeExpiringIntervalTicker(20*time.Millisecond, 210*time.Millisecond),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.True(errors.IsError(err, wait.ErrTickerExceeded))
	assert.Equal(count, 10, "exceeded with a count")

	assert.Logf("end with timeout, no check")
	count = 0
	err = wait.Poll(
		context.Background(),
		wait.MakeExpiringIntervalTicker(20*time.Millisecond, -10*time.Millisecond),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.True(errors.IsError(err, wait.ErrTickerExceeded))
	assert.Equal(count, 0)

	assert.Logf("end with cancelled context")
	count = 0
	ctx, cancel := context.WithTimeout(context.Background(), 110*time.Millisecond)
	defer cancel()
	err = wait.Poll(
		ctx,
		wait.MakeExpiringIntervalTicker(20*time.Millisecond, 210*time.Millisecond),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, "context deadline exceeded")
	assert.Equal(count, 5)
}

// TestPoll tests the polling of conditions with individual ticker.
func TestPoll(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	ticker := func(ctx context.Context) <-chan struct{} {
		// Ticker runs 1000 times.
		tickc := make(chan struct{})
		go func() {
			count := 0
			defer close(tickc)
			for {
				select {
				case tickc <- struct{}{}:
					count++
					if count == 1000 {
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
		return tickc
	}

	// Tests.
	assert.Logf("end with positive condition")
	count := 0
	err := wait.Poll(
		context.Background(),
		ticker,
		func() (bool, error) {
			count++
			if count == 500 {
				return true, nil
			}
			return false, nil
		},
	)
	assert.NoError(err)
	assert.Equal(count, 500)

	assert.Logf("end with timeout, 1000 checks")
	count = 0
	err = wait.Poll(
		context.Background(),
		ticker,
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.True(errors.IsError(err, wait.ErrTickerExceeded))
	assert.Equal(count, 1000, "exceeded with a count")

	assert.Logf("end with cancelled context")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	err = wait.Poll(
		ctx,
		ticker,
		func() (bool, error) {
			time.Sleep(2 * time.Millisecond)
			return false, nil
		},
	)
	assert.ErrorMatch(err, "context deadline exceeded")
}

// EOF
