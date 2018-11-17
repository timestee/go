// Tideland Go Library - Together - Limiter - Unit Tests
//
// Copyright (C) 2017-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package limiter_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"tideland.one/go/audit/asserts"
	"tideland.one/go/together/limiter"
)

//--------------------
// TESTS
//--------------------

// TestLimitOK tests the limiting of a number of function calls
// in multiple goroutines without an error.
func TestLimitOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	max := 0
	act := 0
	job := func() error {
		act++
		if act > max {
			max = act
		}
		time.Sleep(25 * time.Millisecond)
		act--
		return nil
	}
	l := limiter.New(10)
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(25)

	// Test.
	for i := 0; i < 25; i++ {
		go func() {
			err := l.Do(ctx, job)
			assert.NoError(err)
			wg.Done()
		}()
	}

	wg.Wait()
	assert.Equal(max, 10)
}

// TestLimitError tests the returning of en error by an
// executed function.
func TestLimitError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	job := func() error {
		time.Sleep(25 * time.Millisecond)
		return errors.New("ouch")
	}
	l := limiter.New(5)
	ctx := context.Background()

	var wg sync.WaitGroup
	wg.Add(25)

	// Test.
	for i := 0; i < 25; i++ {
		go func() {
			err := l.Do(ctx, job)
			assert.ErrorMatch(err, "ouch")
			wg.Done()
		}()
	}

	wg.Wait()
}

// EOF
