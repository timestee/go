// Tideland Go Library - Together - Mesh - Nodes
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package nodes // import "tideland.dev/go/together/mesh/nodes"

//--------------------
// ERRORS
//--------------------

// Error codes of the runtime package.
const (
	ErrNodeInit     = "ENODEINIT"
	ErrNodeBackend  = "ENODEBACKEND"
	ErrNodeNotFound = "ENODENOTFOUND"

	msgNodeInit     = "node %q cannot initialize"
	msgNodeBackend  = "node %q has a backend failure"
	msgNodeNotFound = "node %q not found"
)

// EOF
