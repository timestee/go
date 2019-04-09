// Tideland Go Library - Text - Simple Markup Language
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package sml

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/trace/errors"
)

//--------------------
// CONSTANTS
//--------------------

// Error codes of the SML package.
const (
	ErrBuilder          = "err-builder"
	ErrReader           = "err-reader"
	ErrNoRootProcessor  = "err-root-processor"
	ErrRegisteredPlugin = "err-plugin"
)

//--------------------
// ERROR
//--------------------

// IsBuilderError checks for an error during node building.
func IsBuilderError(err error) bool {
	return errors.IsError(err, ErrBuilder)
}

// IsReaderError checks for an error during SML text reading.
func IsReaderError(err error) bool {
	return errors.IsError(err, ErrBuilder)
}

// IsNoRootProcessorError checks for an unregistered root
// processor.
func IsNoRootProcessorError(err error) bool {
	return errors.IsError(err, ErrNoRootProcessor)
}

// IsRegisteredPluginError checks for the error of an already
// registered plugin.
func IsRegisteredPluginError(err error) bool {
	return errors.IsError(err, ErrRegisteredPlugin)
}

// EOF
