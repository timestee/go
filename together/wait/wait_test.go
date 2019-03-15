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

// TestPollWithDeadline tests the polling of a condition with durations
// and cancels.
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
}

// TestPollWithTimeout tests the polling of a condition with timeouts
// and cancels.
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
}

// EOF
