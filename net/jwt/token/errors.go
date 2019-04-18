// Tideland Go Library - Network - JSON Web Token
//
// Copyright (C) 2016-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package token // import "tideland.dev/go/net/jwt/token"

//--------------------
// CONSTANTS
//--------------------

// Error codes and messages.
const (
	ErrCannotEncode               = "EENCODE"
	ErrCannotDecode               = "EDECODE"
	ErrCannotVerify               = "EVERIFY"
	ErrCannotSign                 = "ESIGN"
	ErrNoKey                      = "ENOKEY"
	ErrInvalidKey                 = "EINVKEY"
	ErrInvalidTokenPart           = "ETOKEN"
	ErrInvalidSignature           = "ESIGNATURE"
	ErrJSONMarshalling            = "EJSONENC"
	ErrJSONUnmarshalling          = "EJSONDEC"
	ErrInvalidAlgorithm           = "EINVALGO"
	ErrInvalidCombination         = "EINVCOMBI"
	ErrInvalidKeyType             = "EINVKEYTYP"
	ErrCannotReadPEM              = "EREADPEM"
	ErrCannotDecodePEM            = "EDECODEPEM"
	ErrCannotParseECDSA           = "EPARSEDCSA"
	ErrNoECDSAKey                 = "ENOEDCSA"
	ErrCannotParseRSA             = "EPARSERSA"
	ErrNoRSAKey                   = "ENORSA"
	ErrNoAuthorizationHeader      = "ENOAUTH"
	ErrInvalidAuthorizationHeader = "EINVAUTH"

	msgNoAuthorizationHeader      = "request contains no authorization header"
	msgInvalidAuthorizationHeader = "invalid authorization header: %q"
)

// EOF
