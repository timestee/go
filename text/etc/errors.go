// Tideland Go Library - Text - Etc
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package etc

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/trace/errors"
)

//--------------------
// CONSTANTS
//--------------------

// Error codes of the etc package.
const (
	ErrInvalidSourceFormat = "err-inv-source"
	ErrIllegalConfigSource = "err-illegal-config"
	ErrCannotReadFile      = "err-reading-file"
	ErrCannotPostProcess   = "err-post-process"
	ErrInvalidPath         = "err-inv-path"
	ErrCannotSplit         = "err-cannot-split"
	ErrCannotApply         = "err-cannot-apply"
)

//--------------------
// ERROR CHECKING
//--------------------

// IsInvalidPathError checks if a path cannot be found.
func IsInvalidPathError(err error) bool {
	return errors.IsError(err, ErrInvalidPath)
}

// EOF
