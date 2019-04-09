// Tideland Go Library - Network - JSON Web Token
//
// Copyright (C) 2016-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package token

//--------------------
// CONSTANTS
//--------------------

const (
	// Error codes.
	ErrCannotEncode               = "err-encode"
	ErrCannotDecode               = "err-decode"
	ErrCannotVerify               = "err-verify"
	ErrCannotSign                 = "err-sign"
	ErrNoKey                      = "err-no-key"
	ErrInvalidKey                 = "err-inv-key"
	ErrInvalidTokenPart           = "err-token"
	ErrInvalidSignature           = "err-signature"
	ErrJSONMarshalling            = "err-marshalling"
	ErrJSONUnmarshalling          = "err-unmarshalling"
	ErrInvalidAlgorithm           = "err-algorithm"
	ErrInvalidCombination         = "err-combination"
	ErrInvalidKeyType             = "err-key-type"
	ErrCannotReadPEM              = "err-read-pem"
	ErrCannotDecodePEM            = "err-decode-pem"
	ErrCannotParseECDSA           = "err-parse-ecdsa"
	ErrNoECDSAKey                 = "err-no-ecdsa"
	ErrCannotParseRSA             = "err-parse-rsa"
	ErrNoRSAKey                   = "err-no-rsa"
	ErrNoAuthorizationHeader      = "err-no-auth"
	ErrInvalidAuthorizationHeader = "err-inv-auth"

	// Error messages.
	msgNoAuthorizationHeader      = "request contains no authorization header"
	msgInvalidAuthorizationHeader = "invalid authorization header: %q"
)

// EOF
