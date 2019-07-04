// Tideland Go Library - Together - Events - Runtime
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package runtime // import "tideland.dev/go/together/events/runtime"

//--------------------
// ERRORS
//--------------------

// Error codes of the runtime package.
const (
	ErrEngineInit        = "EENGINEINIT"
	ErrEngineBackend     = "EENGINEBACKEND"
	ErrProcessorNotFound = "EPROCNOTFOUND"

	msgEngineInit        = "process engine %q cannot initialize"
	msgEngineBackend     = "process engine %q has a backend failure"
	msgProcessorNotFound = "processor %q not found"
)

// EOF
