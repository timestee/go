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
	ErrInvalidSourceFormat = "E001"
	ErrIllegalConfigSource = "E002"
	ErrCannotReadFile      = "E003"
	ErrCannotPostProcess   = "E004"
	ErrInvalidPath         = "E005"
	ErrCannotSplit         = "E006"
	ErrCannotApply         = "E007"
)

//--------------------
// ERROR CHECKING
//--------------------

// IsInvalidPathError checks if a path cannot be found.
func IsInvalidPathError(err error) bool {
	return errors.IsError(err, ErrInvalidPath)
}

// EOF
