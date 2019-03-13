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

// TestPoll
func TestPoll(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	count := 0

	tests := []struct {
		description string
		interval    time.Duration
		timeout     time.Duration
		condition   wait.Condition
		runs        int
		errCode     string
	}{
		{
			description: "end with timeout",
			interval:    10 * time.Millisecond,
			timeout:     100 * time.Millisecond,
			condition: func() (bool, error) {
				count++
				return false, nil
			},
			runs:    11,
			errCode: wait.ErrWaitTimeout,
		}, {
			description: "end with positive condition",
			interval:    10 * time.Millisecond,
			timeout:     100 * time.Millisecond,
			condition: func() (bool, error) {
				count++
				if count == 5 {
					return true, nil
				}
				return false, nil
			},
			runs:    5,
			errCode: "",
		},
	}

	// Test.
	for i, test := range tests {
		assert.Logf("#%d: %s", i, test.description)
		count = 0
		ctx := context.Background()
		err := wait.Poll(ctx, test.interval, test.timeout, test.condition)
		assert.Equal(count, test.runs)
		if err != nil {
			assert.True(errors.IsError(err, test.errCode))
		}
	}
}

// EOF
