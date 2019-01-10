// Tideland Go Library - Together - Actor
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package actor

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"sync"
	"time"

	"tideland.one/go/together/loop"
	"tideland.one/go/together/notifier"
)

//--------------------
// CONSTANTS
//--------------------

const (
	// DefaultTimeout is used in a DoSync() call.
	DefaultTimeout = 5 * time.Second
)

//--------------------
// RECOVERER
//--------------------

// Recoverer allows a goroutine to react on a panic during its
// work. If it returns nil the goroutine shall continue
// work. Otherwise it will return with an error the gouroutine
// may use for its continued processing.
type Recoverer func(reason interface{}) error

//--------------------
// ACTOR
//--------------------

// Action defines the signature of an actor action.
type Action func() error

// Actor allows to simply use and control a goroutine.
type Actor struct {
	mu      sync.Mutex
	actionC chan Action
	options []loop.Option
	loop    *loop.Loop
	err     error
}

// New creates an Actor with the passed options.
func New(options ...Option) *Actor {
	// Init with options.
	act := &Actor{}
	for _, option := range options {
		if err := option(act); err != nil {
			act.err = err
			return act
		}
	}
	// Ensure default settings.
	if act.actionC == nil {
		act.actionC = make(chan Action, 1)
	}
	// Create loop with its options.
	act.loop = loop.New(act.worker, act.options...)
	return act
}

// Go starts the Actor goroutine and immediately returns its
// own instance to allow statements like act := actor.New().Go().
func (act *Actor) Go() *Actor {
	if act.loop != nil {
		act.loop.Go()
	}
	return act
}

// Work starts the Actors backend.
func (act *Actor) Work() error {
	if act.loop != nil {
		return act.loop.Work()
	}
	return act.err
}

// DoSync executes the actor function and returns when it's done
// or it has the default timeout.
func (act *Actor) DoSync(action Action) error {
	return act.DoSyncTimeout(action, DefaultTimeout)
}

// DoSyncTimeout executes the action and returns when it's done
// or it has a timeout.
func (act *Actor) DoSyncTimeout(action Action, timeout time.Duration) error {
	waitC := make(chan struct{})
	if err := act.DoAsync(func() error {
		err := action()
		close(waitC)
		return err
	}); err != nil {
		return err
	}
	select {
	case <-waitC:
	case <-time.After(timeout):
		return errors.New(errTimeout)
	}
	return nil
}

// DoAsync executes the actor function and returns immediately
func (act *Actor) DoAsync(action Action) error {
	act.mu.Lock()
	defer act.mu.Unlock()
	if act.err != nil {
		return act.err
	}
	act.actionC <- action
	return nil
}

// Stop terminates the Actor with the passed error. That or
// a potential earlier error will be returned.
func (act *Actor) Stop(err error) error {
	if act.loop != nil {
		return act.loop.Stop(err)
	}
	return act.err
}

// Err returns information if the Actor has an error.
func (act *Actor) Err() error {
	act.mu.Lock()
	defer act.mu.Unlock()
	if act.err != nil {
		return act.err
	}
	return act.loop.Err()
}

// worker is the Loop worker of the Actor.
func (act *Actor) worker(c *notifier.Closer) error {
	for {
		select {
		case <-c.Done():
			return nil
		case action := <-act.actionC:
			if err := action(); err != nil {
				return err
			}
		}
	}
}

// EOF
