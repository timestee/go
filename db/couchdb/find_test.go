// Tideland Go Library - DB - CouchDB Client
//
// Copyright (C) 2016-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package couchdb_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/db/couchdb"
)

//--------------------
// TESTS
//--------------------

// TestSimpleFind tests tests calling find with a simple search.
func TestSimpleFind(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	count := 1000
	cdb, cleanup := prepareSizedFilledDatabase(assert, "find-simple", count)
	defer cleanup()

	// Try to find some documents a simple way.
	search := couchdb.NewSearch().
		Selector(`{"$or": [
			{"$and": [
				{"age": {"$lt": 30}},
				{"active": {"$eq": false}}
			]},
			{"$and": [
				{"age": {"$gt": 60}},
				{"active": {"$eq": true}}
			]}
		]}`).
		Fields("name", "age", "active")

	fnds, err := cdb.Find(search)
	assert.NoError(err)

	err = fnds.Process(func(document *couchdb.Unmarshable) error {
		fields := struct {
			Name   string `json:"name"`
			Age    int    `json:"age"`
			Active bool   `json:"active"`
		}{}
		if err := document.Unmarshal(&fields); err != nil {
			return err
		}
		assert.True((fields.Age < 30 && !fields.Active) || (fields.Age > 60 && fields.Active))
		return nil
	})
	assert.Nil(err)
}

// EOF
