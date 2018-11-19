// Tideland Go Library - Network - JSON Web Token - Crypto
//
// Copyright (C) 2016-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package crypto

//--------------------
// CONSTANTS
//--------------------

// Error codes of the crypto package.
const (
	ErrInvalidAlgorithm   = "E201"
	ErrInvalidCombination = "E202"
	ErrInvalidKeyType     = "E203"
	ErrCannotSign         = "E204"
	ErrCannotVerify       = "E205"
	ErrInvalidKey         = "E206"
	ErrInvalidSignature   = "E207"
	ErrCannotReadPEM      = "E221"
	ErrCannotDecodePEM    = "E222"
	ErrCannotParseECDSA   = "E231"
	ErrNoECDSAKey         = "E232"
	ErrCannotParseRSA     = "E241"
	ErrNoRSAKey           = "E242"
)

// EOF
