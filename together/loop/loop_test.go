// Tideland Go Library - Together - Loop - Unit Tests
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package loop_test

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"errors"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/together/loop"
	"tideland.dev/go/together/notifier"
)

//--------------------
// CONSTANTS
//--------------------

// timeout is the waitng time for events from inside of loops.
var timeout time.Duration = 5 * time.Second

//--------------------
// TESTS
//--------------------

// TestPure tests a loop without any options, stopping without an error.
func TestPure(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	waitC := make(chan struct{})
	beenThereDoneThat := false
	worker := func(c *notifier.Closer) error {
		defer func() {
			beenThereDoneThat = true
		}()
		close(waitC)
		for {
			select {
			case <-c.Done():
				return nil
			}
		}
	}
	l := loop.New(worker).Go()

	// Test.
	<-waitC
	assert.Equal(l.Status(), notifier.Working)
	assert.NoError(l.Stop(nil))
	assert.Equal(l.Status(), notifier.Stopped)
	assert.True(beenThereDoneThat)
}

// TestPureError tests a loop without any options, stopping with an error.
func TestPureError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	waitC := make(chan struct{})
	worker := func(c *notifier.Closer) error {
		close(waitC)
		for {
			select {
			case <-c.Done():
				return nil
			}
		}
	}
	l := loop.New(worker).Go()

	// Test.
	<-waitC
	err := l.Stop(errors.New("ouch"))
	assert.ErrorMatch(err, "ouch")
}

// TestPureInternalError tests a loop without any options, stopping leads
// to an internal error.
func TestPureInternalError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	waitC := make(chan struct{})
	worker := func(c *notifier.Closer) (err error) {
		defer func() {
			err = errors.New("ouch")
		}()
		close(waitC)
		for {
			select {
			case <-c.Done():
				return nil
			}
		}
	}
	l := loop.New(worker).Go()

	// Test.
	<-waitC
	err := l.Stop(nil)
	assert.ErrorMatch(err, "ouch")
}

// TestContextCancelOK tests the stopping after a context cancel w/o error.
func TestContextCancelOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	ctx, cancel := context.WithCancel(context.Background())
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return nil
			}
		}
	}
	notifier := notifier.New()
	l := loop.New(worker,
		loop.WithContext(ctx),
		loop.WithNotifier(notifier))

	// Test.
	<-notifier.Ready()
	l.Go()
	<-notifier.Working()
	cancel()
	<-notifier.Stopped()
	assert.NoError(l.Err())
}

// TestContextCancelError tests the stopping after a context cancel w/ error.
func TestContextCancelError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	ctx, cancel := context.WithCancel(context.Background())
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return errors.New("oh, no")
			}
		}
	}
	notifier := notifier.New()
	l := loop.New(worker,
		loop.WithContext(ctx),
		loop.WithNotifier(notifier))

	// Test.
	<-notifier.Ready()
	l.Go()
	<-notifier.Working()
	cancel()
	<-notifier.Stopped()
	assert.ErrorMatch(l.Err(), "oh, no")
}

// TestMultipleNotifier tests the usage of multiple notifiers.
func TestMultipleNotifier(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return nil
			}
		}
	}
	notifierA := notifier.New()
	notifierB := notifier.New()
	notifierC := notifier.New()
	l := loop.New(worker,
		loop.WithNotifier(notifierA),
		loop.WithNotifier(notifierB),
		loop.WithNotifier(notifierC)).Go()

	// Test.
	<-notifierC.Working()
	l.Stop(nil)

	x := 0

	for x != 7 {
		select {
		case <-notifierA.Stopped():
			x |= 1
		case <-notifierB.Stopped():
			x |= 2
		case <-notifierC.Stopped():
			x |= 4
		case <-time.After(time.Second):
			break
		}
	}

	assert.Equal(x, 7)
	assert.NoError(l.Err())
}

// TestFinalizerOK tests the stopping with an error, but cleared by a finalizer.
func TestFinalizerOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	finalized := false
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return errors.New("don't want to stop")
			}
		}
	}
	finalizer := func(err error) error {
		assert.ErrorMatch(err, "don't want to stop")
		finalized = true
		return nil
	}
	notifier := notifier.New()
	l := loop.New(worker,
		loop.WithFinalizer(finalizer),
		loop.WithNotifier(notifier)).Go()

	// Test.
	<-notifier.Working()
	l.Stop(nil)
	<-notifier.Stopped()
	assert.NoError(l.Err())
	assert.True(finalized)
}

// TestFinalizerError tests the stopping with an error, but changed by a finalizer.
func TestContextCancelFinalizerError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return errors.New("don't want to stop")
			}
		}
	}
	finalizer := func(err error) error {
		assert.ErrorMatch(err, "don't want to stop")
		return errors.New("don't care")
	}
	notifier := notifier.New()
	l := loop.New(worker,
		loop.WithFinalizer(finalizer),
		loop.WithNotifier(notifier)).Go()

	// Test.
	<-notifier.Working()
	l.Stop(nil)
	<-notifier.Stopped()
	assert.ErrorMatch(l.Err(), "don't care")
}

// TestInternalOK tests the stopping w/o an error.
func TestInternalOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return nil
			case <-time.After(50 * time.Millisecond):
				return nil
			}
		}
	}
	notifier := notifier.New()
	l := loop.New(worker,
		loop.WithNotifier(notifier)).Go()

	// Test.
	<-notifier.Stopped()
	assert.NoError(l.Err())
}

// TestInternalError tests the stopping after an internal error.
func TestInternalError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return nil
			case <-time.After(50 * time.Millisecond):
				return errors.New("time over")
			}
		}
	}
	notifier := notifier.New()
	l := loop.New(worker,
		loop.WithNotifier(notifier)).Go()

	// Test.
	<-notifier.Stopped()
	assert.ErrorMatch(l.Stop(nil), "time over")
	assert.ErrorMatch(l.Err(), "time over")
}

// TestRecoveredOK tests the stopping without an error if Loop has a recoverer.
// Recoverer must never been called.
func TestRecoveredOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	beenThereDoneThat := false
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return nil
			}
		}
	}
	recoverer := func(reason interface{}) error {
		beenThereDoneThat = true
		return nil
	}
	notifier := notifier.New()
	l := loop.New(worker,
		loop.WithRecoverer(recoverer),
		loop.WithNotifier(notifier)).Go()

	// Test.
	<-notifier.Working()
	l.Stop(nil)
	<-notifier.Stopped()
	assert.Nil(l.Err())
	assert.False(beenThereDoneThat)
}

// TestRecovererError tests the stopping with an error if Loop has a recoverer.
// Recoverer must never been called.
func TestRecovererError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	beenThereDoneThat := false
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return errors.New("oh, no")
			}
		}
	}
	recoverer := func(reason interface{}) error {
		beenThereDoneThat = true
		return nil
	}
	notifier := notifier.New()
	l := loop.New(worker,
		loop.WithRecoverer(recoverer),
		loop.WithNotifier(notifier)).Go()

	// Test.
	<-notifier.Working()
	l.Stop(nil)
	<-notifier.Stopped()
	assert.ErrorMatch(l.Err(), "oh, no")
	assert.False(beenThereDoneThat)
}

// TestRecoverPanicsOK tests the stopping w/o an error.
func TestRecoverPanicsOK(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	panics := 0
	doneC := make(chan struct{})
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return nil
			case <-doneC:
				return nil
			case <-time.After(10 * time.Millisecond):
				panic("bam")
			}
		}
	}
	recoverer := func(reason interface{}) error {
		panics++
		if panics > 10 {
			close(doneC)
		}
		return nil
	}
	notifier := notifier.New()
	l := loop.New(worker,
		loop.WithRecoverer(recoverer),
		loop.WithNotifier(notifier)).Go()

	// Test.
	<-notifier.Stopped()
	assert.NoError(l.Err())
}

// TestRecoverPanicsError tests the stopping w/o an error.
func TestRecoverPanicsError(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	panics := 0
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return nil
			case <-time.After(10 * time.Millisecond):
				panic("bam")
			}
		}
	}
	recoverer := func(reason interface{}) error {
		panics++
		if panics > 10 {
			return errors.New("superbam")
		}
		return nil
	}
	notifier := notifier.New()
	l := loop.New(worker,
		loop.WithRecoverer(recoverer),
		loop.WithNotifier(notifier)).Go()

	// Test.
	<-notifier.Stopped()
	assert.ErrorMatch(l.Err(), "superbam")
}

// TestReasons tests collecting and analysing loop recovery reasons.
func TestReasons(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	rins := []error{
		errors.New("error a"),
		errors.New("error b"),
		errors.New("error c"),
		errors.New("error d"),
		errors.New("error e"),
	}
	rs := loop.MakeReasons()
	for _, rin := range rins {
		time.Sleep(100 * time.Millisecond)
		rs = rs.Append(rin)
	}

	// Test.
	assert.Length(rs, 5)
	assert.Equal(rs.Last().Reason, rins[4])
	assert.True(rs.Frequency(5, time.Second))
	assert.False(rs.Frequency(5, 10*time.Millisecond))

	rs = rs.Trim(3)

	assert.Length(rs, 3)
	assert.Match(rs.String(), `\[\['error c' @ .*\] / \['error d' @ .*\] / \['error e' @ .*\]\]`)
}

//--------------------
// EXAMPLES
//--------------------

// ExampleWorker shows the usage of Loo with no recoverer. The inner loop
// contains a select listening to the channel returned by Closer.Done().
// Other channels are for the standard communication with the Loop.
func ExampleWorker() {
	printC := make(chan string)
	ctx, cancel := context.WithCancel(context.Background())
	// Sample loop worker.
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				// We shall stop.
				return nil
			case str := <-printC:
				// Standard work of example loop.
				if str == "panic" {
					return errors.New("panic")
				}
				println(str)
			}
		}
	}
	l := loop.New(worker, loop.WithContext(ctx)).Go()

	printC <- "Hello"
	printC <- "World"

	cancel()

	if l.Err() != nil {
		panic(l.Err())
	}
}

// ExampleRecoverer demonstrates the usage of a recoverer.
// Here the frequency of the recovered reasons (more than five
// in 10 milliseconds) or the total number is checked. If the
// total number is not interesting the reasons could be
// trimmed by e.g. rs.Trim(5). The fields Time and Reason per
// recovering allow even more diagnosis.
func ExampleRecoverer() {
	panicC := make(chan string)
	// Sample loop worker.
	worker := func(c *notifier.Closer) error {
		for {
			select {
			case <-c.Done():
				return nil
			case str := <-panicC:
				panic(str)
			}
		}
	}
	// Recovery function checking frequency and total number.
	rs := loop.MakeReasons()
	recoverer := func(reason interface{}) error {
		rs = rs.Append(reason)
		if rs.Frequency(5, 10*time.Millisecond) {
			return errors.New("too high error frequency")
		}
		if rs.Len() >= 10 {
			return errors.New("too many errors")
		}
		return nil
	}
	loop.New(worker, loop.WithRecoverer(recoverer)).Go()
}

// EOF
