// Tideland Go Library - Network - REST - Handlers - Unit Tests
//
// Copyright (C) 2009-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package handlers_test

//--------------------
// IMPORTS
//--------------------

import (
	"bufio"
	"context"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"tideland.one/go/audit/asserts"
	"tideland.one/go/net/jwt/cache"
	"tideland.one/go/net/jwt/token"
	"tideland.one/go/net/rest/audit"
	"tideland.one/go/net/rest/core"
	"tideland.one/go/net/rest/handlers"
	"tideland.one/go/text/etc"
)

//--------------------
// TESTS
//--------------------

// TestWrapperHandler tests the usage of standard handler funcs
// wrapped to be used inside the package context.
func TestWrapperHandler(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	data := "Been there, done that!"
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	handler := func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte(data))
	}
	err := mux.Register("test", "wrapper", handlers.NewWrapperHandler("wrapper", handler))
	assert.Nil(err)
	// Perform test requests.
	req := audit.NewRequest("GET", "/test/wrapper")
	resp := ts.DoRequest(req)
	resp.AssertBodyContains(data)
	punctuation := resp.AssertBodyGrep("[,!]")
	assert.Length(punctuation, 2)
}

// TestFileServeHandler tests the serving of files.
func TestFileServeHandler(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	data := "Been there, done that!"
	// Setup the test file.
	dir, err := ioutil.TempDir("", "gorest")
	assert.Nil(err)
	defer os.RemoveAll(dir)
	filename := filepath.Join(dir, "foo.txt")
	f, err := os.Create(filename)
	assert.Nil(err)
	_, err = f.WriteString(data)
	assert.Nil(err)
	assert.Logf("written %s", f.Name())
	err = f.Close()
	assert.Nil(err)
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	err = mux.Register("test", "files", handlers.NewFileServeHandler("files", dir))
	assert.Nil(err)
	// Perform test requests.
	req := audit.NewRequest("GET", "/test/files/foo.txt")
	resp := ts.DoRequest(req)
	resp.AssertBodyContains(data)
	req = audit.NewRequest("GET", "/test/files/does.not.exist")
	resp = ts.DoRequest(req)
	resp.AssertBodyContains("404 page not found")
}

// TestFileUploadHandler tests the uploading of files.
func TestFileUploadHandler(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	data := "Been there, done that!"
	// Setup the file upload processor.
	processor := func(job *core.Job, header *multipart.FileHeader, file multipart.File) error {
		assert.Equal(header.Filename, "test.txt")
		scanner := bufio.NewScanner(file)
		assert.True(scanner.Scan())
		text := scanner.Text()
		assert.Equal(text, data)
		return nil
	}
	// Setup the test server.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	err := mux.Register("test", "files", handlers.NewFileUploadHandler("files", processor))
	assert.Nil(err)
	// Perform test requests.
	ts.DoUpload("/test/files", "testfile", "test.txt", data)
}

// TestJWTAuthorizationHandler tests the authorization process
// using JSON Web Tokens.
func TestJWTAuthorizationHandler(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	key := []byte("secret")
	tests := []struct {
		id      string
		tokener func() *token.JWT
		config  *handlers.JWTAuthorizationConfig
		runs    int
		status  int
		body    string
		auditf  handlers.AuditHandlerFunc
	}{
		{
			id:     "no-token",
			status: http.StatusUnauthorized,
		}, {
			id: "token-decode-no-gatekeeper",
			tokener: func() *token.JWT {
				claims := token.NewClaims()
				claims.SetSubject("test")
				out, err := token.Encode(claims, key, token.HS512)
				assert.Nil(err)
				return out
			},
			status: http.StatusOK,
			auditf: func(assert *asserts.Asserts, job *core.Job) (bool, error) {
				token, ok := token.FromContext(job.Context())
				assert.True(ok)
				assert.NotNil(token)
				subject, ok := token.Claims().Subject()
				assert.True(ok)
				assert.Equal(subject, "test")
				return true, nil
			},
		}, {
			id: "token-verify-no-gatekeeper",
			tokener: func() *token.JWT {
				claims := token.NewClaims()
				claims.SetSubject("test")
				out, err := token.Encode(claims, key, token.HS512)
				assert.Nil(err)
				return out
			},
			config: &handlers.JWTAuthorizationConfig{
				Key: key,
			},
			status: http.StatusOK,
		}, {
			id: "cached-token-verify-no-gatekeeper",
			tokener: func() *token.JWT {
				claims := token.NewClaims()
				claims.SetSubject("test")
				out, err := token.Encode(claims, key, token.HS512)
				assert.Nil(err)
				return out
			},
			config: &handlers.JWTAuthorizationConfig{
				Cache: cache.New(context.Background(), time.Minute, time.Minute, time.Minute, 10),
				Key:   key,
			},
			runs:   5,
			status: http.StatusOK,
		}, {
			id: "cached-token-verify-positive-gatekeeper",
			tokener: func() *token.JWT {
				claims := token.NewClaims()
				claims.SetSubject("test")
				out, err := token.Encode(claims, key, token.HS512)
				assert.Nil(err)
				return out
			},
			config: &handlers.JWTAuthorizationConfig{
				Cache: cache.New(context.Background(), time.Minute, time.Minute, time.Minute, 10),
				Key:   key,
				Gatekeeper: func(job *core.Job, claims token.Claims) error {
					subject, ok := claims.Subject()
					assert.True(ok)
					assert.Equal(subject, "test")
					return nil
				},
			},
			runs:   5,
			status: http.StatusOK,
		}, {
			id: "cached-token-verify-negative-gatekeeper",
			tokener: func() *token.JWT {
				claims := token.NewClaims()
				claims.SetSubject("test")
				out, err := token.Encode(claims, key, token.HS512)
				assert.Nil(err)
				return out
			},
			config: &handlers.JWTAuthorizationConfig{
				Cache: cache.New(context.Background(), time.Minute, time.Minute, time.Minute, 10),
				Key:   key,
				Gatekeeper: func(job *core.Job, claims token.Claims) error {
					_, ok := claims.Subject()
					assert.True(ok)
					return errors.New("subject is test")
				},
			},
			runs:   1,
			status: http.StatusUnauthorized,
		}, {
			id: "token-expired",
			tokener: func() *token.JWT {
				claims := token.NewClaims()
				claims.SetSubject("test")
				claims.SetExpiration(time.Now().Add(-time.Hour))
				out, err := token.Encode(claims, key, token.HS512)
				assert.Nil(err)
				return out
			},
			status: http.StatusForbidden,
		},
	}
	// Run defined tests.
	mux := newMultiplexer(assert)
	ts := audit.StartServer(mux, assert)
	defer ts.Close()
	for i, test := range tests {
		// Prepare one test.
		assert.Logf("JWT test #%d: %s", i, test.id)
		err := mux.Register("jwt", test.id, handlers.NewJWTAuthorizationHandler(test.id, test.config))
		assert.Nil(err)
		if test.auditf != nil {
			err := mux.Register("jwt", test.id, handlers.NewAuditHandler("audit", assert, test.auditf))
			assert.Nil(err)
		}
		// Create request.
		req := audit.NewRequest("GET", "/jwt/"+test.id+"/1234567890")
		if test.tokener != nil {
			req.SetRequestProcessor(func(req *http.Request) *http.Request {
				return token.RequestAdd(req, test.tokener())
			})
		}
		// Make request(s).
		runs := 1
		if test.runs != 0 {
			runs = test.runs
		}
		for i := 0; i < runs; i++ {
			resp := ts.DoRequest(req)
			resp.AssertStatusEquals(test.status)
		}
	}
}

//--------------------
// HELPERS
//--------------------

// newMultiplexer creates a new multiplexer with a testing context
// and a testing configuration.
func newMultiplexer(assert *asserts.Asserts) *core.Multiplexer {
	cfgStr := "{etc {basepath /}{default-domain default}{default-resource default}}"
	cfg, err := etc.ReadString(cfgStr)
	assert.Nil(err)
	return core.NewMultiplexer(context.Background(), cfg)
}

// EOF
