// Tideland Go Library - Trace - Logger
//
// Copyright (C) 2012-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package logger

//--------------------
// IMPORTS
//--------------------

import (
	"context"
)

//--------------------
// CONTEXT
//--------------------

// contextKey describes the type of the context key.
type contextKey int

// loggerContextKey is the context key for a logger.
const loggerContextKey contextKey = 1

// NewContext creates a context containing a logger.
func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// FromContext retrieves a logger from a context.
func FromContext(ctx context.Context) (Logger, bool) {
	logger, ok := ctx.Value(loggerContextKey).(Logger)
	return logger, ok
}

// EOF
