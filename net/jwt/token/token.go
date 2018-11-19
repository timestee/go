// Tideland Go Library - Network - JSON Web Token
//
// Copyright (C) 2016-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package token

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"tideland.one/go/net/jwt/claims"
	"tideland.one/go/net/jwt/crypto"
	"tideland.one/go/trace/errors"
)

//--------------------
// JSON Web Token
//--------------------

// JWT describes the interface to access the parts of a
// JSON Web Token.
type JWT interface {
	// Stringer provides the String() method.
	fmt.Stringer

	// Claims returns the claims payload of the token.
	Claims() claims.Claims

	// Key return the key of the token only when
	// it is a result of encoding or verification.
	Key() (crypto.Key, error)

	// Algorithm returns the algorithm of the token
	// after encoding, decoding, or verification.
	Algorithm() crypto.Algorithm

	// IsValid is a convenience method checking the
	// registered claims if the token is valid.
	IsValid(leeway time.Duration) bool
}

// jwtHeader contains the JWT header fields.
type jwtHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

// jwt implements JWT.
type jwt struct {
	claims    claims.Claims
	key       crypto.Key
	algorithm crypto.Algorithm
	token     string
}

// Encode creates a JSON Web Token for the given claims
// based on key and algorithm.
func Encode(claims claims.Claims, key crypto.Key, algorithm crypto.Algorithm) (JWT, error) {
	jwt := &jwt{
		claims:    claims,
		key:       key,
		algorithm: algorithm,
	}
	headerPart, err := marshallAndEncode(jwtHeader{string(algorithm), "JWT"})
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotEncode, "cannot encode the header")
	}
	claimsPart, err := marshallAndEncode(claims)
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotEncode, "cannot encode the claims")
	}
	dataParts := headerPart + "." + claimsPart
	signaturePart, err := signAndEncode([]byte(dataParts), key, algorithm)
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotEncode, " cannot encode the signature")
	}
	jwt.token = dataParts + "." + signaturePart
	return jwt, nil
}

// Decode creates a token out of a string without verification.
func Decode(token string) (JWT, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New(ErrCannotDecode, "cannot decode the parts")
	}
	var header jwtHeader
	err := decodeAndUnmarshall(parts[0], &header)
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotDecode, "cannot decode the header")
	}
	var claims claims.Claims
	err = decodeAndUnmarshall(parts[1], &claims)
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotDecode, "cannot decode the claims")
	}
	return &jwt{
		claims:    claims,
		algorithm: crypto.Algorithm(header.Algorithm),
		token:     token,
	}, nil
}

// Verify creates a token out of a string and varifies it against
// the passed key.
func Verify(token string, key crypto.Key) (JWT, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New(ErrCannotVerify, "cannot verify the parts")
	}
	var header jwtHeader
	err := decodeAndUnmarshall(parts[0], &header)
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotVerify, "cannot verify the header")
	}
	err = decodeAndVerify(parts, key, crypto.Algorithm(header.Algorithm))
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotVerify, "cannot verify the signature")
	}
	var claims claims.Claims
	err = decodeAndUnmarshall(parts[1], &claims)
	if err != nil {
		return nil, errors.Annotate(err, ErrCannotVerify, "cannot verify the claims")
	}
	return &jwt{
		claims:    claims,
		key:       key,
		algorithm: crypto.Algorithm(header.Algorithm),
		token:     token,
	}, nil
}

// Claims implements the JWT interface.
func (jwt *jwt) Claims() claims.Claims {
	return jwt.claims
}

// Key implements the JWT interface.
func (jwt *jwt) Key() (crypto.Key, error) {
	if jwt.key == nil {
		return nil, errors.New(ErrNoKey, "no key available, only after encoding or verifying")
	}
	return jwt.key, nil
}

// Algorithm implements the JWT interface.
func (jwt *jwt) Algorithm() crypto.Algorithm {
	return jwt.algorithm
}

// IsValid implements the JWT interface.
func (jwt *jwt) IsValid(leeway time.Duration) bool {
	return jwt.claims.IsValid(leeway)
}

// String implements the Stringer interface.
func (jwt *jwt) String() string {
	return jwt.token
}

//--------------------
// PRIVATE HELPERS
//--------------------

// marshallAndEncode marshals the passed value to JSON and
// creates a BASE64 string out of it.
func marshallAndEncode(value interface{}) (string, error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return "", errors.Annotate(err, ErrJSONMarshalling, "error marshalling to JSON")
	}
	encoded := base64.RawURLEncoding.EncodeToString(jsonValue)
	return encoded, nil
}

// decodeAndUnmarshall decodes a BASE64 encoded JSON string and
// unmarshals it into the passed value.
func decodeAndUnmarshall(part string, value interface{}) error {
	decoded, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		return errors.Annotate(err, ErrInvalidTokenPart, "part of the token contains invalid data")
	}
	err = json.Unmarshal(decoded, value)
	if err != nil {
		return errors.Annotate(err, ErrJSONUnmarshalling, "error unmarshalling from JSON")
	}
	return nil
}

// signAndEncode creates the signature for the data part (header and
// payload) of the token using the passed key and algorithm. The result
// is then encoded to BASE64.
func signAndEncode(data []byte, key crypto.Key, algorithm crypto.Algorithm) (string, error) {
	sig, err := algorithm.Sign(data, key)
	if err != nil {
		return "", err
	}
	encoded := base64.RawURLEncoding.EncodeToString(sig)
	return encoded, nil
}

// decodeAndVerify decodes a BASE64 encoded signature and verifies
// the correct signing of the data part (header and payload) using the
// passed key and algorithm.
func decodeAndVerify(parts []string, key crypto.Key, algorithm crypto.Algorithm) error {
	data := []byte(parts[0] + "." + parts[1])
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return errors.Annotate(err, ErrInvalidTokenPart, "part of the token contains invalid data")
	}
	return algorithm.Verify(data, sig, key)
}

// EOF
