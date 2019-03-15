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

// TestPoll tests the polling of conditions.
func TestPoll(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)

	// Tests.
	assert.Logf("end with positive condition")
	count := 0
	err := wait.Poll(
		context.Background(),
		50*time.Millisecond,
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
	ctx, cancel := context.WithTimeout(context.Background(), 275*time.Millisecond)
	defer cancel()
	err = wait.Poll(
		ctx,
		50*time.Millisecond,
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
	err := wait.PollWithDeadline(
		context.Background(),
		50*time.Millisecond,
		time.Now().Add(500*time.Millisecond),
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
	err = wait.PollWithDeadline(
		context.Background(),
		50*time.Millisecond,
		time.Now().Add(500*time.Millisecond),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.True(errors.IsError(err, wait.ErrWaitTimeout))
	assert.Equal(count, 10, "have a count, but still a timeout")

	assert.Logf("end with deadline, no check")
	count = 0
	err = wait.PollWithDeadline(
		context.Background(),
		50*time.Millisecond,
		time.Now().Add(-time.Second),
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.True(errors.IsError(err, wait.ErrWaitTimeout))
	assert.Equal(count, 0)

	assert.Logf("end with cancelled context")
	count = 0
	ctx, cancel := context.WithTimeout(context.Background(), 275*time.Millisecond)
	defer cancel()
	err = wait.PollWithDeadline(
		ctx,
		50*time.Millisecond,
		time.Now().Add(500*time.Millisecond),
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
	err := wait.PollWithTimeout(
		context.Background(),
		50*time.Millisecond,
		500*time.Millisecond,
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
	err = wait.PollWithTimeout(
		context.Background(),
		50*time.Millisecond,
		500*time.Millisecond,
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.True(errors.IsError(err, wait.ErrWaitTimeout))
	assert.Equal(count, 10, "have a count, but still a timeout")

	assert.Logf("end with timeout, no check")
	count = 0
	err = wait.PollWithTimeout(
		context.Background(),
		50*time.Millisecond,
		-time.Second,
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.True(errors.IsError(err, wait.ErrWaitTimeout))
	assert.Equal(count, 0)

	assert.Logf("end with cancelled context")
	count = 0
	ctx, cancel := context.WithTimeout(context.Background(), 275*time.Millisecond)
	defer cancel()
	err = wait.PollWithTimeout(
		ctx,
		50*time.Millisecond,
		500*time.Millisecond,
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.ErrorMatch(err, "context deadline exceeded")
	assert.Equal(count, 5)
}

// TestPollWithTicker tests the polling of conditions with individual tickers.
func TestPollWithTicker(t *testing.T) {
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
	err := wait.PollWithTicker(
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
	err = wait.PollWithTicker(
		context.Background(),
		ticker,
		func() (bool, error) {
			count++
			return false, nil
		},
	)
	assert.True(errors.IsError(err, wait.ErrWaitTimeout))
	assert.Equal(count, 1000, "have a count, but still a timeout")

	assert.Logf("end with cancelled context")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	err = wait.PollWithTicker(
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
