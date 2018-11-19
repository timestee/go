// Tideland Go Library - Network - JSON Web Token - Cache - Unit Tests
//
// Copyright (C) 2016-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cache_test

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"testing"
	"time"

	"tideland.one/go/audit/asserts"
	"tideland.one/go/net/jwt/cache"
	"tideland.one/go/net/jwt/claims"
	"tideland.one/go/net/jwt/crypto"
	"tideland.one/go/net/jwt/token"
)

//--------------------
// TESTS
//--------------------

// TestCachePutGet tests the putting and getting of tokens
// to the cache.
func TestCachePutGet(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	assert.Logf("testing cache put and get")
	cache := cache.New(time.Minute, time.Minute, time.Minute, 10)
	key := []byte("secret")
	claims := initClaims()
	jwtIn, err := token.Encode(claims, key, crypto.HS512)
	assert.Nil(err)
	cache.Put(jwtIn)
	token := jwtIn.String()
	jwtOut, ok := cache.Get(token)
	assert.True(ok)
	assert.Equal(jwtIn, jwtOut)
	jwtOut, ok = cache.Get("is.not.there")
	assert.False(ok)
	assert.Nil(jwtOut)
	err = cache.Stop()
	assert.Nil(err)
}

// TestCacheAccessCleanup tests the access based cleanup
// of the JWT cache.
func TestCacheAccessCleanup(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	assert.Logf("testing cache access based cleanup")
	cache := cache.New(time.Second, time.Second, time.Second, 10)
	key := []byte("secret")
	claims := initClaims()
	jwtIn, err := token.Encode(claims, key, crypto.HS512)
	assert.Nil(err)
	cache.Put(jwtIn)
	token := jwtIn.String()
	jwtOut, ok := cache.Get(token)
	assert.True(ok)
	assert.Equal(jwtIn, jwtOut)
	// Now wait a bit an try again.
	time.Sleep(5 * time.Second)
	jwtOut, ok = cache.Get(token)
	assert.False(ok)
	assert.Nil(jwtOut)
}

// TestCacheValidityCleanup tests the validity based cleanup
// of the JWT cache.
func TestCacheValidityCleanup(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	assert.Logf("testing cache validity based cleanup")
	cache := cache.New(time.Minute, time.Second, time.Second, 10)
	key := []byte("secret")
	now := time.Now()
	nbf := now.Add(-2 * time.Second)
	exp := now.Add(2 * time.Second)
	claims := initClaims()
	claims.SetNotBefore(nbf)
	claims.SetExpiration(exp)
	jwtIn, err := token.Encode(claims, key, crypto.HS512)
	assert.Nil(err)
	cache.Put(jwtIn)
	token := jwtIn.String()
	jwtOut, ok := cache.Get(token)
	assert.True(ok)
	assert.Equal(jwtIn, jwtOut)
	// Now access until it is invalid and not
	// available anymore.
	var i int
	for i = 0; i < 5; i++ {
		time.Sleep(time.Second)
		jwtOut, ok = cache.Get(token)
		if !ok {
			break
		}
		assert.Equal(jwtIn, jwtOut)
	}
	assert.True(i > 1 && i < 4)
}

// TestCacheLoad tests the cache load based cleanup.
func TestCacheLoad(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	assert.Logf("testing cache load based cleanup")
	cacheTime := 100 * time.Millisecond
	cache := cache.New(2*cacheTime, cacheTime, cacheTime, 4)
	claims := initClaims()
	// Now fill the cache and check that it doesn't
	// grow too high.
	var i int
	for i = 0; i < 10; i++ {
		time.Sleep(50 * time.Millisecond)
		key := []byte(fmt.Sprintf("secret-%d", i))
		jwtIn, err := token.Encode(claims, key, crypto.HS512)
		assert.Nil(err)
		size := cache.Put(jwtIn)
		assert.True(size < 6)
	}
}

//--------------------
// HELPERS
//--------------------

// initClaims creates test claims.
func initClaims() claims.Claims {
	c := claims.New()
	c.SetSubject("1234567890")
	c.Set("name", "John Doe")
	c.Set("admin", true)
	return c
}

// EOF
