// Tideland Go Library - Together - Cells - Mesh
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package mesh // import "tideland.dev/go/together/cells/mesh"

//--------------------
// ERRORS
//--------------------

// Error codes of the runtime package.
const (
	ErrCellInit     = "ECELLINIT"
	ErrCellBackend  = "ECELLBACKEND"
	ErrCellNotFound = "ECELLNOTFOUND"

	msgCellInit     = "cell %q cannot initialize"
	msgCellBackend  = "cell %q has a backend failure"
	msgCellNotFound = "cell %q not found"
)

// EOF
