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
	"context"

	"tideland.one/go/together/loop"
	"tideland.one/go/together/notifier"
)

//--------------------
// OPTIONS
//--------------------

// Option defines the signature of an option setting function.
type Option func(act *Actor) error

// WithContext allows to pass a context for cancellation or timeout.
func WithContext(ctx context.Context) Option {
	return func(act *Actor) error {
		act.options = append(act.options, loop.WithContext(ctx))
		return nil
	}
}

// WithQueueLen defines the channel size for actions to
// send to an Actor.
func WithQueueLen(size int) Option {
	return func(act *Actor) error {
		if size < 1 {
			size = 1
		}
		act.actionC = make(chan Action, size)
		return nil
	}
}

// WithRecoverer defines the panic handler of an actor.
func WithRecoverer(recoverer Recoverer) Option {
	return func(act *Actor) error {
		act.options = append(act.options, loop.WithRecoverer(loop.Recoverer(recoverer)))
		return nil
	}
}

// WithNotifier add a notifier to make external monitors aware of
// the Actors internal status.
func WithNotifier(notifier *notifier.Notifier) Option {
	return func(act *Actor) error {
		act.options = append(act.options, loop.WithNotifier(notifier))
		return nil
	}
}

// WithFinalizer sets a function for finalizing the
// work of a Loop.
func WithFinalizer(finalizer loop.Finalizer) Option {
	return func(act *Actor) error {
		act.options = append(act.options, loop.WithFinalizer(finalizer))
		return nil
	}
}

// EOF
