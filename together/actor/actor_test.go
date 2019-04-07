// Tideland Go Library - Together - Actor - Unit Tests
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package actor_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"errors"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/actor"
	"tideland.dev/go/together/notifier"
)

//--------------------
// TESTS
//--------------------

// TestPureGoOK is simply starting and stopping an Actor
// with Go().
func TestPureGoOK(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	notifier := notifier.New()
	act := actor.New(
		actor.WithNotifier(notifier)).Go()
	assert.NotNil(act)

	<-notifier.Working()
	assert.NoError(act.Stop(nil))
	assert.NoError(act.Err())
}

// TestPureGoError is simply starting and stopping an Actor
// with Go(). Returning the stop error.
func TestPureGo(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	notifier := notifier.New()
	act := actor.New(
		actor.WithNotifier(notifier)).Go()
	assert.NotNil(act)

	<-notifier.Working()
	assert.ErrorMatch(act.Stop(errors.New("damn")), "damn")
	assert.ErrorMatch(act.Err(), "damn")
}

// TestWithContextGo is simply starting and stopping an Actor
// with a context and Go().
func TestNewActorGo(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	act := actor.New(
		actor.WithContext(ctx)).Go()
	assert.NotNil(act)

	cancel()
	assert.NoError(act.Err())
}

// TestWithContextWork is simply starting an Actor with Work()
// and terminate externally.
func TestNewActorWork(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	act := actor.New(
		actor.WithContext(ctx))
	assert.NotNil(act)

	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()

	assert.NoError(act.Work())
}

// TestSync tests synchronous calls.
func TestSync(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	act := actor.New().Go()
	defer act.Stop(nil)

	counter := 0

	for i := 0; i < 5; i++ {
		err := act.DoSync(func() error {
			counter++
			return nil
		})
		assert.Nil(err)
	}

	assert.Equal(counter, 5)
}

// TestTimeout tests timout error of a synchronous Action.
func TestTimeout(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	act := actor.New().Go()
	defer act.Stop(nil)

	// Scenario: Timeout is shorter than needed time, so error
	// is returned.
	err := act.DoSyncTimeout(func() error {
		time.Sleep(time.Second)
		return nil
	}, 500*time.Millisecond)

	assert.True(actor.IsTimedOut(err))
}

// TestAsyncWithQueueLen tests running multiple calls asynchronously.
func TestAsyncWithQueueLen(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	act := actor.New(
		actor.WithQueueLen(100)).Go()
	defer act.Stop(nil)

	sigC := make(chan bool, 1)
	doneC := make(chan bool, 1)

	// Start background func waiting for the signals of
	// the asynchrounous calls.
	go func() {
		count := 0
		for range sigC {
			count++
			if count == 100 {
				break
			}
		}
		doneC <- true
	}()

	// Now start asynchrounous calls.
	start := time.Now()
	for i := 0; i < 100; i++ {
		act.DoAsync(func() error {
			time.Sleep(5 * time.Millisecond)
			sigC <- true
			return nil
		})
	}
	enqueued := time.Since(start)

	// Expect signal done to be sent about one second later.
	<-doneC
	done := time.Since(start)

	assert.True((done - 500*time.Millisecond) > enqueued)
}

// TestRecovery tests handling panics.
func TestRecovery(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	counter := 0
	recovered := false
	doneC := make(chan struct{})
	recoverer := func(reason interface{}) error {
		recovered = true
		close(doneC)
		return nil
	}
	act := actor.New(
		actor.WithRecoverer(recoverer)).Go()
	defer act.Stop(nil)

	err := act.DoSyncTimeout(func() error {
		counter++
		print(counter / (counter - counter))
		return nil
	}, time.Second)
	assert.True(actor.IsTimedOut(err))
	<-doneC
	assert.True(recovered)
	err = act.DoSync(func() error {
		counter++
		return nil
	})
	assert.NoError(err)
	assert.Equal(counter, 2)
}

// EOF
