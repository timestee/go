// Tideland Go Library - DB - CouchDB Client
//
// Copyright (C) 2016-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package couchdb

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/base64"
)

//--------------------
// PARAMETERIZABLE
//--------------------

// KeyValue is used for generic query and header parameters.
type KeyValue struct {
	Key   string
	Value string
}

//--------------------
// PARAMETERS
//--------------------

// Parameter is a function changing one (or if needed multile) parameter.
type Parameter func(req *Request)

// Query is generic for setting request query parameters.
func Query(kvs ...KeyValue) Parameter {
	return func(req *Request) {
		for _, kv := range kvs {
			req.AddQuery(kv.Key, kv.Value)
		}
	}
}

// Header is generic for setting request header parameters.
func Header(kvs ...KeyValue) Parameter {
	return func(req *Request) {
		for _, kv := range kvs {
			req.SetHeader(kv.Key, kv.Value)
		}
	}
}

// Revision sets the revision for the access to concrete document revisions.
func Revision(revision string) Parameter {
	return func(req *Request) {
		req.SetQuery("rev", revision)
	}
}

// BasicAuthentication is intended for basic authentication
// against the database.
func BasicAuthentication(name, password string) Parameter {
	return func(req *Request) {
		np := []byte(name + ":" + password)
		auth := "Basic " + base64.StdEncoding.EncodeToString(np)

		req.SetHeader("Authorization", auth)
	}
}

// EOF
