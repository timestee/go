// Tideland Go Library - Together - Loop - Unit Tests
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package notifier_test

//--------------------
// IMPORTS
//--------------------

import (
	"math/rand"
	"sync"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/notifier"
)

//--------------------
// CONSTANTS
//--------------------

// timeout is the waitng time..
var timeout time.Duration = 5 * time.Second

//--------------------
// TESTS
//--------------------

// TestCloserOK tests the closing of the Closer when one of the input closer
// channels closes.
func TestCloserOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	ccs := []chan struct{}{
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
	}
	closer := notifier.NewCloser(ccs[0], ccs[1], ccs[2], ccs[3], ccs[4]).Go()
	beenThereDoneThat := false

	// Test.
	go func() {
		time.Sleep(time.Second)
		i := rand.Int() % len(ccs)
		close(ccs[i])
	}()

	select {
	case <-closer.Done():
		beenThereDoneThat = true
	case <-time.After(timeout):
		assert.Fail("timeout")
	}
	assert.True(beenThereDoneThat)
}

// TestCloserTimeout tests what happes if no channel signals the closing.
func TestCloserTimeout(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	closer := notifier.NewCloser(make(chan struct{})).Go()
	beenThereDoneThat := false

	// Test.
	select {
	case <-closer.Done():
		assert.Fail("invalid closing")
	case <-time.After(timeout):
		beenThereDoneThat = true
	}
	assert.True(beenThereDoneThat)
}

// TestNotifiersOK tests the notification of multiple Notifier through one Notifiers.
func TestNotifiersOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	na := notifier.New()
	nb := notifier.New()
	nc := notifier.New()

	var status sync.WaitGroup
	status.Add(12)
	var waiterStart sync.WaitGroup
	waiterStart.Add(3)

	waiter := func(n *notifier.Notifier) {
		waiterStart.Done()

		<-n.Ready()
		status.Done()

		<-n.Working()
		status.Done()

		<-n.Stopping()
		status.Done()

		<-n.Stopped()
		status.Done()
	}
	b := notifier.NewBundle()
	b.Add(na, nb, nc)
	beenThereDoneThat := false
	waitC := make(chan interface{}, 1)

	// Test.
	go waiter(na)
	go waiter(nb)
	go waiter(nc)

	go func() {
		waiterStart.Wait()
		time.Sleep(50 * time.Millisecond)
		b.Notify(notifier.Ready)
		time.Sleep(50 * time.Millisecond)
		b.Notify(notifier.Working)
		time.Sleep(50 * time.Millisecond)
		b.Notify(notifier.Stopping)
		time.Sleep(50 * time.Millisecond)
		b.Notify(notifier.Stopped)

		status.Wait()
		waitC <- "been there, done that"
	}()

	assert.Wait(waitC, "been there, done that", timeout)

	select {
	case <-b.Stopped():
		beenThereDoneThat = true
	case <-time.After(timeout):
		assert.Fail("timeout")
	}

	assert.True(beenThereDoneThat)
}

// TestNotifiersMulti tests the multiple setting of a status and the multiple query
// of a status.
func TestNotifiersMulti(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	n := notifier.New()
	b := notifier.NewBundle()
	b.Add(n)
	beenThereDoneThat := false

	// Test.
	b.Notify(notifier.Ready)
	b.Notify(notifier.Ready)
	b.Notify(notifier.Ready)

	select {
	case <-n.Ready():
		beenThereDoneThat = true
	case <-time.After(timeout):
		assert.Fail("timeout")
	}

	assert.True(beenThereDoneThat)
	assert.Equal(n.Status(), notifier.Ready)

	beenThereDoneThat = false

	select {
	case <-n.Ready():
		beenThereDoneThat = true
	case <-time.After(timeout):
		assert.Fail("timeout")
	}

	assert.True(beenThereDoneThat)
}

// EOF
