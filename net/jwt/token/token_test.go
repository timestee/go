// Tideland Go Library - Network - JSON Web Token - Unit Tests
//
// Copyright (C) 2016-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package token_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"
	"time"

	"tideland.one/go/audit/asserts"
	"tideland.one/go/net/jwt/claims"
	"tideland.one/go/net/jwt/crypto"
	"tideland.one/go/net/jwt/token"
)

//--------------------
// TESTS
//--------------------

const (
	subClaim   = "1234567890"
	nameClaim  = "John Doe"
	adminClaim = true
	iatClaim   = 1600000000
	rawToken   = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9." +
		"eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTYwMDAwMDAwMH0." +
		"P50peTbENKIPw0tjuHLgosFmJRYGTh_kNA9IcyWIoJ39uYMa4JfKYhnQw5mkgSLB2WYVT68QaDeWWErn4lU69g"
)

// TestDecode tests the decoding without verifying the signature.
func TestDecode(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	// Decode.
	jwt, err := token.Decode(rawToken)
	assert.Nil(err)
	assert.Equal(jwt.Algorithm(), crypto.HS512)
	key, err := jwt.Key()
	assert.Nil(key)
	assert.ErrorMatch(err, ".*no key available, only after encoding or verifying.*")
	assert.Length(jwt.Claims(), 4)

	sub, ok := jwt.Claims().GetString("sub")
	assert.True(ok)
	assert.Equal(sub, subClaim)
	name, ok := jwt.Claims().GetString("name")
	assert.True(ok)
	assert.Equal(name, nameClaim)
	admin, ok := jwt.Claims().GetBool("admin")
	assert.True(ok)
	assert.Equal(admin, adminClaim)
	iat, ok := jwt.Claims().IssuedAt()
	assert.True(ok)
	assert.Equal(iat, time.Unix(iatClaim, 0))
	exp, ok := jwt.Claims().Expiration()
	assert.False(ok)
	assert.Equal(exp, time.Time{})
}

// TestIsValid checks the time validation of a token.
func TestIsValid(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	assert.Logf("testing time validation")
	now := time.Now()
	leeway := time.Minute
	key := []byte("secret")
	// Create token with no times set, encode, decode, validate ok.
	clms := claims.New()
	jwtEnc, err := token.Encode(clms, key, crypto.HS512)
	assert.Nil(err)
	jwtDec, err := token.Decode(jwtEnc.String())
	assert.Nil(err)
	ok := jwtDec.IsValid(leeway)
	assert.True(ok)
	// Now a token with a long timespan, still valid.
	clms = claims.New()
	clms.SetNotBefore(now.Add(-time.Hour))
	clms.SetExpiration(now.Add(time.Hour))
	jwtEnc, err = token.Encode(clms, key, crypto.HS512)
	assert.Nil(err)
	jwtDec, err = token.Decode(jwtEnc.String())
	assert.Nil(err)
	ok = jwtDec.IsValid(leeway)
	assert.True(ok)
	// Now a token with a long timespan in the past, not valid.
	clms = claims.New()
	clms.SetNotBefore(now.Add(-2 * time.Hour))
	clms.SetExpiration(now.Add(-time.Hour))
	jwtEnc, err = token.Encode(clms, key, crypto.HS512)
	assert.Nil(err)
	jwtDec, err = token.Decode(jwtEnc.String())
	assert.Nil(err)
	ok = jwtDec.IsValid(leeway)
	assert.False(ok)
	// And at last a token with a long timespan in the future, not valid.
	clms = claims.New()
	clms.SetNotBefore(now.Add(time.Hour))
	clms.SetExpiration(now.Add(2 * time.Hour))
	jwtEnc, err = token.Encode(clms, key, crypto.HS512)
	assert.Nil(err)
	jwtDec, err = token.Decode(jwtEnc.String())
	assert.Nil(err)
	ok = jwtDec.IsValid(leeway)
	assert.False(ok)
}

// EOF