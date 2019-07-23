// Tideland Go Library - Together - Cells - Event
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package event // import "tideland.dev/go/together/cells/event"

//--------------------
// ERRORS
//--------------------

// Error codes of the runtime package.
const (
	ErrNoValue    = "ENOVAL"
	ErrConverting = "ECONV"

	msgNoValue    = "key %q has no value"
	msgConverting = "value of key %q cannot be converted to %v"
)

// EOF
