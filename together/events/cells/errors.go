// Tideland Go Library - Together - Events - Cells
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package cells // import "tideland.dev/go/together/events/cells"

//--------------------
// ERRORS
//--------------------

// Error codes of the cells package.
const (
	ErrRuntimeAdd  = "ERUNTIMEADD"
	ErrCellInit    = "ECELLINIT"
	ErrCellBackend = "ECELLBACKEND"

	msgRuntimeAdd  = "double cell identifier %q cannot be added"
	msgCellInit    = "cell %q cannot initialize"
	msgCellBackend = "cell %q has a backend failure"
)

// EOF
