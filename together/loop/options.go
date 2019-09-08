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
	"context"

	"tideland.dev/go/together/notifier"
	"tideland.dev/go/trace/failure"
)

//--------------------
// OPTIONS
//--------------------

// Option defines the signature of an option setting function.
type Option func(l *Loop) error

// WithContext allows to pass a context for cancellation or timeout.
func WithContext(ctx context.Context) Option {
	return func(l *Loop) error {
		if ctx == nil {
			return failure.New("invalid loop option: context is nil")
		}
		l.closeCs = append(l.closeCs, ctx.Done())
		return nil
	}
}

// WithRecoverer defines the panic handler of a loop.
func WithRecoverer(recoverer Recoverer) Option {
	return func(l *Loop) error {
		if recoverer == nil {
			return failure.New("invalid loop option: recoverer is nil")
		}
		l.recover = recoverer
		return nil
	}
}

// WithNotifier add a notifier to make external monitors aware of
// the Loop internal status.
func WithNotifier(notifier *notifier.Notifier) Option {
	return func(l *Loop) error {
		if notifier == nil {
			return failure.New("invalid loop option: notifier is nil")
		}
		l.bundle.Add(notifier)
		return nil
	}
}

// WithFinalizer sets a function for finalizing the
// work of a Loop.
func WithFinalizer(finalizer Finalizer) Option {
	return func(l *Loop) error {
		if finalizer == nil {
			return failure.New("invalid loop option: finalizer is nil")
		}
		l.finalize = finalizer
		return nil
	}
}

// EOF
