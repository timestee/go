// Tideland Go Library - Together - Loop
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package loop // import "tideland.dev/go/together/loop"

//--------------------
// IMPORTS
//--------------------

import (
	"sync"
	"time"

	"tideland.dev/go/together/notifier"
	"tideland.dev/go/trace/failure"
)

//--------------------
// RECOVERER
//--------------------

// Recoverer allows a goroutine to react on a panic during its
// work. If it returns nil the goroutine shall continue
// work. Otherwise it will return with an error the gouroutine
// may use for its continued processing.
type Recoverer func(reason interface{}) error

// DefaultRecoverer simply re-panics.
func DefaultRecoverer(reason interface{}) error {
	panic(reason)
}

//--------------------
// ERROR HELPERS
//--------------------

// IsErrLoopNotWorking helps loop users to check if an errors
// tells that it currently isn't running (anymore).
func IsErrLoopNotWorking(err error) bool {
	return failure.Contains(err, "loop not working")
}

//--------------------
// LOOP
//--------------------

// Worker is a managed Loop function performing the work.
type Worker func(c *notifier.Closer) error

// Finalizer allows to perform some steps to clean-up when
// the worker terminates. The passed error is the state of
// the loop.
type Finalizer func(err error) error

// Loop manages running for-select-loops in the background as goroutines
// in a controlled way. Users can get information about status and possible
// failure as well as control how to stop, restart, or recover via
// options.
type Loop struct {
	mu       sync.RWMutex
	work     Worker
	finalize Finalizer
	closeC   chan struct{}
	closeCs  []<-chan struct{}
	closer   *notifier.Closer
	bundle   *notifier.Bundle
	recover  Recoverer
	err      error
}

// New creates a new loop running the passed worker with the set options.
func New(w Worker, options ...Option) *Loop {
	// Init with default values.
	l := &Loop{
		work:    w,
		closeC:  make(chan struct{}),
		bundle:  notifier.NewBundle(),
		recover: DefaultRecoverer,
	}
	l.closeCs = append(l.closeCs, (<-chan struct{})(l.closeC))
	// Apply options.
	for _, option := range options {
		if err := option(l); err != nil {
			// One of the options made troubles.
			l.err = err
			if l.bundle != nil {
				l.bundle.Notify(notifier.Stopped)
			}
			return l
		}
	}
	// Ensure default settings, first close channel is
	// to stop directly.
	l.closer = notifier.NewCloser(l.closeCs...).Go()
	l.bundle.Notify(notifier.Ready)
	return l
}

// Go starts the Loop worker as goroutine and immediately returns its
// own instance to allow statements like l := loop.New(worker).Go().
func (l *Loop) Go() *Loop {
	// Check if status is correct.
	l.mu.Lock()
	if l.bundle.Status() != notifier.Ready {
		l.mu.Unlock()
		return l
	}
	l.mu.Unlock()
	// Start work as goroutine.
	go l.Work()
	return l
}

// Work starts the Loop worker synchronously.
func (l *Loop) Work() error {
	// Check if status is correct.
	l.mu.Lock()
	if l.bundle.Status() != notifier.Ready {
		l.mu.Unlock()
		return failure.New("loop not ready")
	}
	// Start working.
	defer l.bundle.Notify(notifier.Stopped)
	l.bundle.Notify(notifier.Working)
	l.mu.Unlock()
	for l.bundle.Status() == notifier.Working {
		l.container()
	}
	if l.finalize != nil {
		l.err = l.finalize(l.err)
	}
	return l.err
}

// Stop terminates the Loop with the passed error. That or
// a potential earlier error will be returned.
func (l *Loop) Stop(err error) error {
	// Check if status is correct.
	l.mu.Lock()
	if l.bundle.Status() != notifier.Working {
		l.mu.Unlock()
		if l.err != nil {
			return l.err
		}
		return failure.New("loop not working")
	}
	defer l.mu.Unlock()
	// Stop and wait.
	close(l.closeC)
	select {
	case <-l.bundle.Stopped():
	case <-time.After(30 * time.Second):
		l.err = failure.New("timeout during stopping")
	}
	if l.err == nil {
		l.err = err
	}
	return l.err
}

// Status returns status information about the Loop.
func (l *Loop) Status() notifier.Status {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.bundle.Status()
}

// Err returns information if the Loop has an error.
func (l *Loop) Err() error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.err
}

// container wraps the worker, handles possible failure, and
// manages panics.
func (l *Loop) container() {
	defer func() {
		if reason := recover(); reason != nil {
			// Panic, try to recover.
			if err := l.recover(reason); err != nil {
				l.err = err
				l.bundle.Notify(notifier.Stopping)
			}
		} else {
			// Regular ending.
			l.bundle.Notify(notifier.Stopping)
		}
	}()
	l.err = l.work(l.closer)
}

// EOF
