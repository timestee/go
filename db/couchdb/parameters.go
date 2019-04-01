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
	"encoding/json"
	"strconv"
)

//--------------------
// CONSTANTS
//--------------------

// Fixed values for some of the view parameters.
const (
	SinceNow = "now"

	StyleMainOnly = "main_only"
	StyleAllDocs  = "all_docs"
)

//--------------------
// PARAMETERS
//--------------------

// KeyValue is used for generic query and header parameters.
type KeyValue struct {
	Key   string
	Value string
}

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

// Limit sets the maximum number of result rows.
func Limit(limit int) Parameter {
	return func(req *Request) {
		req.SetQuery("limit", strconv.Itoa(limit))
	}
}

// Since sets the start of the changes gathering, can also be "now".
func Since(sequence string) Parameter {
	return func(req *Request) {
		req.SetQuery("since", sequence)
	}
}

// Descending sets the flag for a descending order of changes gathering.
func Descending() Parameter {
	return func(req *Request) {
		req.SetQuery("descending", "true")
	}
}

// Style sets how many revisions are returned. Default is
// StyleMainOnly only returning the winning document revision.
// StyleAllDocs will return all revision including possible
// conflicts.
func Style(style string) Parameter {
	return func(req *Request) {
		req.SetQuery("style", style)
	}
}

// FilterDocumentIDs sets a filtering of the changes to the
// given document identifiers.
func FilterDocumentIDs(documentIDs ...string) Parameter {
	update := func(doc interface{}) interface{} {
		if doc == nil {
			doc = &couchdbDocumentIDs{}
		}
		idsdoc, ok := doc.(*couchdbDocumentIDs)
		if ok {
			idsdoc.DocumentIDs = append(idsdoc.DocumentIDs, documentIDs...)
			return idsdoc
		}
		return doc
	}
	return func(req *Request) {
		req.SetQuery("filter", "_doc_ids")
		req.UpdateDocument(update)
	}
}

// FilterSelector sets the filter to the passed selector expression.
func FilterSelector(selector json.RawMessage) Parameter {
	update := func(doc interface{}) interface{} {
		// TODO 2019-03-31 Mue Set selector expression.
		return doc
	}
	return func(req *Request) {
		req.SetQuery("filter", "_selector")
		req.UpdateDocument(update)
	}
}

// FilterView sets the name of a view which map function acts as
// filter in case it emits at least one record.
func FilterView(view string) Parameter {
	return func(req *Request) {
		req.SetQuery("filter", "_view")
		req.SetQuery("view", view)
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
