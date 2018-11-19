// Tideland Go Library - Network - JSON Web Token - Crypto - Unit Tests
//
// Copyright (C) 2016-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package crypto_test

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"testing"

	"tideland.one/go/audit/asserts"
	"tideland.one/go/net/jwt/crypto"
)

//--------------------
// TESTS
//--------------------

var (
	esTests = []crypto.Algorithm{crypto.ES256, crypto.ES384, crypto.ES512}
	hsTests = []crypto.Algorithm{crypto.HS256, crypto.HS384, crypto.HS512}
	psTests = []crypto.Algorithm{crypto.PS256, crypto.PS384, crypto.PS512}
	rsTests = []crypto.Algorithm{crypto.RS256, crypto.RS384, crypto.RS512}
	data    = []byte("the quick brown fox jumps over the lazy dog")
)

// TestESAlgorithms tests the ECDSA algorithms.
func TestESAlgorithms(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(err)
	for _, algo := range esTests {
		assert.Logf("testing algorithm %q", algo)
		// Sign.
		signature, err := algo.Sign(data, privateKey)
		assert.Nil(err)
		assert.NotEmpty(signature)
		// Verify.
		err = algo.Verify(data, signature, privateKey.Public())
		assert.Nil(err)
	}
}

// TestHSAlgorithms tests the HMAC algorithms.
func TestHSAlgorithms(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	key := []byte("secret")
	for _, algo := range hsTests {
		assert.Logf("testing algorithm %q", algo)
		// Sign.
		signature, err := algo.Sign(data, key)
		assert.Nil(err)
		assert.NotEmpty(signature)
		// Verify.
		err = algo.Verify(data, signature, key)
		assert.Nil(err)
	}
}

// TestPSAlgorithms tests the RSAPSS algorithms.
func TestPSAlgorithms(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)
	for _, algo := range psTests {
		assert.Logf("testing algorithm %q", algo)
		// Sign.
		signature, err := algo.Sign(data, privateKey)
		assert.Nil(err)
		assert.NotEmpty(signature)
		// Verify.
		err = algo.Verify(data, signature, privateKey.Public())
		assert.Nil(err)
	}
}

// TestRSAlgorithms tests the RSA algorithms.
func TestRSAlgorithms(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)
	for _, algo := range rsTests {
		assert.Logf("testing algorithm %q", algo)
		// Sign.
		signature, err := algo.Sign(data, privateKey)
		assert.Nil(err)
		assert.NotEmpty(signature)
		// Verify.
		err = algo.Verify(data, signature, privateKey.Public())
		assert.Nil(err)
	}
}

// TestNoneAlgorithm tests the none algorithm.
func TestNoneAlgorithm(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	assert.Logf("testing algorithm \"none\"")
	// Sign.
	signature, err := crypto.NONE.Sign(data, "")
	assert.Nil(err)
	assert.Empty(signature)
	// Verify.
	err = crypto.NONE.Verify(data, signature, "")
	assert.Nil(err)
}

// TestNotMatchingAlgorithm checks when algorithms of
// signing and verifying don't match.'
func TestNotMatchingAlgorithm(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	esPrivateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	esPublicKey := esPrivateKey.Public()
	assert.Nil(err)
	hsKey := []byte("secret")
	rsPrivateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	rsPublicKey := rsPrivateKey.Public()
	assert.Nil(err)
	noneKey := ""
	errorMatch := ".* combination of algorithm .* and key type .*"
	tests := []struct {
		description string
		algorithm   crypto.Algorithm
		key         crypto.Key
		signKeys    []crypto.Key
		verifyKeys  []crypto.Key
	}{
		{"ECDSA", crypto.ES512, esPrivateKey,
			[]crypto.Key{hsKey, rsPrivateKey, noneKey}, []crypto.Key{hsKey, rsPublicKey, noneKey}},
		{"HMAC", crypto.HS512, hsKey,
			[]crypto.Key{esPrivateKey, rsPrivateKey, noneKey}, []crypto.Key{esPublicKey, rsPublicKey, noneKey}},
		{"RSA", crypto.RS512, rsPrivateKey,
			[]crypto.Key{esPrivateKey, hsKey, noneKey}, []crypto.Key{esPublicKey, hsKey, noneKey}},
		{"RSAPSS", crypto.PS512, rsPrivateKey,
			[]crypto.Key{esPrivateKey, hsKey, noneKey}, []crypto.Key{esPublicKey, hsKey, noneKey}},
		{"none", crypto.NONE, noneKey,
			[]crypto.Key{esPrivateKey, hsKey, rsPrivateKey}, []crypto.Key{esPublicKey, hsKey, rsPublicKey}},
	}
	// Run the tests.
	for _, test := range tests {
		assert.Logf("testing %q algorithm key type mismatch", test.description)
		for _, key := range test.signKeys {
			_, err := test.algorithm.Sign(data, key)
			assert.ErrorMatch(err, errorMatch)
		}
		signature, err := test.algorithm.Sign(data, test.key)
		assert.Nil(err)
		for _, key := range test.verifyKeys {
			err = test.algorithm.Verify(data, signature, key)
			assert.ErrorMatch(err, errorMatch)
		}
	}
}

// TestESTools tests the tools for the reading of PEM encoded
// ECDSA keys.
func TestESTools(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	assert.Logf("testing \"ECDSA\" tools")
	// Generate keys and PEMs.
	privateKeyIn, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	assert.Nil(err)
	privateBytes, err := x509.MarshalECPrivateKey(privateKeyIn)
	assert.Nil(err)
	privateBlock := pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateBytes,
	}
	privatePEM := pem.EncodeToMemory(&privateBlock)
	publicBytes, err := x509.MarshalPKIXPublicKey(privateKeyIn.Public())
	assert.Nil(err)
	publicBlock := pem.Block{
		Type:  "EC PUBLIC KEY",
		Bytes: publicBytes,
	}
	publicPEM := pem.EncodeToMemory(&publicBlock)
	assert.NotNil(publicPEM)
	// Now read them.
	buf := bytes.NewBuffer(privatePEM)
	privateKeyOut, err := crypto.ReadECPrivateKey(buf)
	assert.Nil(err)
	buf = bytes.NewBuffer(publicPEM)
	publicKeyOut, err := crypto.ReadECPublicKey(buf)
	assert.Nil(err)
	// And as a last step check if they are correctly usable.
	signature, err := crypto.ES512.Sign(data, privateKeyOut)
	assert.Nil(err)
	err = crypto.ES512.Verify(data, signature, publicKeyOut)
	assert.Nil(err)
}

// TestRSTools tests the tools for the reading of PEM encoded
// RSA keys.
func TestRSTools(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	assert.Logf("testing \"RSA\" tools")
	// Generate keys and PEMs.
	privateKeyIn, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)
	privateBytes := x509.MarshalPKCS1PrivateKey(privateKeyIn)
	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateBytes,
	}
	privatePEM := pem.EncodeToMemory(&privateBlock)
	publicBytes, err := x509.MarshalPKIXPublicKey(privateKeyIn.Public())
	assert.Nil(err)
	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicBytes,
	}
	publicPEM := pem.EncodeToMemory(&publicBlock)
	assert.NotNil(publicPEM)
	// Now read them.
	buf := bytes.NewBuffer(privatePEM)
	privateKeyOut, err := crypto.ReadRSAPrivateKey(buf)
	assert.Nil(err)
	buf = bytes.NewBuffer(publicPEM)
	publicKeyOut, err := crypto.ReadRSAPublicKey(buf)
	assert.Nil(err)
	// And as a last step check if they are correctly usable.
	signature, err := crypto.RS512.Sign(data, privateKeyOut)
	assert.Nil(err)
	err = crypto.RS512.Verify(data, signature, publicKeyOut)
	assert.Nil(err)
}

// EOF
