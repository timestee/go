// Tideland Go Library - Network - JSON Web Token
//
// Copyright (C) 2016-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package token

//--------------------
// CONSTANTS
//--------------------

const (
	// Error codes.
	ErrCannotEncode       = "E101"
	ErrCannotDecode       = "E102"
	ErrCannotVerify       = "E103"
	ErrCannotSign         = "E104"
	ErrNoKey              = "E111"
	ErrInvalidKey         = "E112"
	ErrInvalidTokenPart   = "E113"
	ErrInvalidSignature   = "E114"
	ErrJSONMarshalling    = "E191"
	ErrJSONUnmarshalling  = "E192"
	ErrInvalidAlgorithm   = "E201"
	ErrInvalidCombination = "E202"
	ErrInvalidKeyType     = "E203"
	ErrCannotReadPEM      = "E221"
	ErrCannotDecodePEM    = "E222"
	ErrCannotParseECDSA   = "E231"
	ErrNoECDSAKey         = "E232"
	ErrCannotParseRSA     = "E241"
	ErrNoRSAKey           = "E242"
)

// EOF
