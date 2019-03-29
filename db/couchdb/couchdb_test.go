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
	"tideland.dev/go/audit/generators"
	"tideland.dev/go/dsa/identifier"
	"tideland.dev/go/trace/errors"
	"tideland.dev/go/trace/logger"

	"tideland.dev/go/db/couchdb"
)

//--------------------
// CONSTANTS
//--------------------

const (
	// testDB is the name of the database used for testing.
	testDB = "tmp-couchdb-testing"
)

//--------------------
// TESTS
//--------------------

// TestInvalidConfiguration tests opening the database with an invalid
//  configuration.
func TestInvalidConfiguration(t *testing.T) {
	assert := asserts.NewTesting(t, true)

	// Open with illegal configuration is okay, only
	// usage of this will fail.
	cdb, err := couchdb.Open(couchdb.Host("some-non-existing-host", 12345))
	assert.Nil(err)

	// Deleting the database has to fail.
	resp := cdb.Manager().DeleteDatabase()
	assert.Equal(resp.StatusCode(), couchdb.StatusBadRequest)
}

// TestCreateDesignDocument tests creating new design documents.
func TestCreateDesignDocument(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	cdb, cleanup := prepareFilledDatabase(assert, "tmp-create-design")
	defer cleanup()

	// Create design document and check if it has been created.
	designIDsA, err := cdb.Designs().IDs()
	assert.Nil(err)

	design, err := cdb.Designs().Design("testing-a")
	assert.Nil(err)
	assert.Equal(design.ID(), "testing-a")
	design.SetView("index-a", "function(doc){ if (doc._id.indexOf('a') !== -1) { emit(doc._id, doc._rev);  } }", "")
	resp := design.Write()
	assert.True(resp.IsOK())

	design, err = cdb.Designs().Design("testing-b")
	assert.Nil(err)
	assert.Equal(design.ID(), "testing-b")
	design.SetView("index-b", "function(doc){ if (doc._id.indexOf('b') !== -1) { emit(doc._id, doc._rev);  } }", "")
	resp = design.Write()
	assert.True(resp.IsOK())

	designIDsB, err := cdb.Designs().IDs()
	assert.Nil(err)
	assert.Equal(len(designIDsB), len(designIDsA)+2)
}

// TestReadDesignDocument tests reading design documents.
func TestReadDesignDocument(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	cdb, cleanup := prepareFilledDatabase(assert, "tmp-read-design")
	defer cleanup()

	// Create design document and read it again.
	designA, err := cdb.Designs().Design("testing-a")
	assert.Nil(err)
	assert.Equal(designA.ID(), "testing-a")
	designA.SetView("index-a", "function(doc){ if (doc._id.indexOf('a') !== -1) { emit(doc._id, doc._rev);  } }", "")
	resp := designA.Write()
	assert.True(resp.IsOK())

	designB, err := cdb.Designs().Design("testing-a")
	assert.Nil(err)
	assert.Equal(designB.ID(), "testing-a")
}

// TestUpdateDesignDocument tests updating design documents.
func TestUpdateDesignDocument(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	cdb, cleanup := prepareFilledDatabase(assert, "tmp-update-design")
	defer cleanup()

	// Create design document and read it again.
	designA, err := cdb.Designs().Design("testing-a")
	assert.Nil(err)
	assert.Equal(designA.ID(), "testing-a")
	designA.SetView("index-a", "function(doc){ if (doc._id.indexOf('a') !== -1) { emit(doc._id, doc._rev);  } }", "")
	resp := designA.Write()
	assert.True(resp.IsOK())

	designB, err := cdb.Designs().Design("testing-a")
	assert.Nil(err)
	assert.Equal(designB.ID(), "testing-a")

	// Now update it and read it again.
	designB.SetView("index-b", "function(doc){ if (doc._id.indexOf('b') !== -1) { emit(doc._id, doc._rev);  } }", "")
	resp = designB.Write()
	assert.True(resp.IsOK())

	designC, err := cdb.Designs().Design("testing-a")
	assert.Nil(err)
	assert.Equal(designC.ID(), "testing-a")
	_, _, ok := designC.View("index-a")
	assert.True(ok)
	_, _, ok = designC.View("index-b")
	assert.True(ok)
}

// TestDeleteDesignDocument tests deleting design documents.
func TestDeleteDesignDocument(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	cdb, cleanup := prepareFilledDatabase(assert, "tmp-delete-design")
	defer cleanup()

	// Create design document and check if it has been created.
	designIDsA, err := cdb.Designs().IDs()
	assert.Nil(err)

	designA, err := cdb.Designs().Design("testing")
	assert.Nil(err)
	designA.SetView("index-a", "function(doc){ if (doc._id.indexOf('a') !== -1) { emit(doc._id, doc._rev);  } }", "")
	resp := designA.Write()
	assert.True(resp.IsOK())

	designIDsB, err := cdb.Designs().IDs()
	assert.Nil(err)
	assert.Equal(len(designIDsB), len(designIDsA)+1)

	// Read it and delete it.
	designB, err := cdb.Designs().Design("testing")
	assert.Nil(err)
	resp = designB.Delete()
	assert.True(resp.IsOK())

	designIDsC, err := cdb.Designs().IDs()
	assert.Nil(err)
	assert.Equal(len(designIDsC), len(designIDsA))
}

// TestCreateDocument tests creating new documents.
func TestCreateDocument(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	cdb, cleanup := prepareDatabase(assert, "tmp-create-document")
	defer cleanup()

	// Create document without ID.
	docA := MyDocument{
		Name: "foo",
		Age:  50,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	assert.Match(id, "[0-9a-f]{32}")

	// Create document with ID.
	docB := MyDocument{
		DocumentID: "bar-12345",
		Name:       "bar",
		Age:        25,
		Active:     true,
	}
	resp = cdb.CreateDocument(docB)
	assert.True(resp.IsOK())
	id = resp.ID()
	assert.Equal(id, "bar-12345")
}

// TestReadDocument tests reading a document.
func TestReadDocument(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	cdb, cleanup := prepareDatabase(assert, "tmp-read-document")
	defer cleanup()

	// Create test document.
	docA := MyDocument{
		DocumentID: "foo-12345",
		Name:       "foo",
		Age:        18,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	assert.Equal(id, "foo-12345")

	// Read test document.
	resp = cdb.ReadDocument(id)
	assert.True(resp.IsOK())
	docB := MyDocument{}
	err := resp.Document(&docB)
	assert.Nil(err)
	assert.Equal(docB.DocumentID, docA.DocumentID)
	assert.Equal(docB.Name, docA.Name)
	assert.Equal(docB.Age, docA.Age)

	// Try to read non-existent document.
	resp = cdb.ReadDocument("i-do-not-exist")
	assert.False(resp.IsOK())
	assert.ErrorMatch(resp.Error(), ".* 404,.*")
}

// TestUpdateDocument tests updating documents.
func TestUpdateDocument(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	cdb, cleanup := prepareDatabase(assert, "tmp-update-document")
	defer cleanup()

	// Create first revision.
	docA := MyDocument{
		DocumentID: "foo-12345",
		Name:       "foo",
		Age:        22,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	revision := resp.Revision()
	assert.Equal(id, "foo-12345")

	resp = cdb.ReadDocument(id)
	assert.True(resp.IsOK())
	docB := MyDocument{}
	err := resp.Document(&docB)
	assert.Nil(err)

	// Update the document.
	docB.Age = 23

	resp = cdb.UpdateDocument(docB)
	assert.True(resp.IsOK())

	// Read the updated revision.
	resp = cdb.ReadDocument(id)
	assert.True(resp.IsOK())
	docC := MyDocument{}
	err = resp.Document(&docC)
	assert.Nil(err)
	assert.Equal(docC.DocumentID, docB.DocumentID)
	assert.Substring("2-", docC.DocumentRevision)
	assert.Equal(docC.Name, docB.Name)
	assert.Equal(docC.Age, docB.Age)

	// Read the first revision.
	resp = cdb.ReadDocument(id, couchdb.Revision(revision))
	assert.True(resp.IsOK())
	assert.Equal(resp.Revision(), revision)

	// Try to update a non-existent document.
	docD := MyDocument{
		DocumentID: "i-do-not-exist",
		Name:       "none",
		Age:        999,
	}
	resp = cdb.UpdateDocument(docD)
	assert.False(resp.IsOK())
	assert.Equal(resp.StatusCode(), couchdb.StatusNotFound)
	assert.True(errors.IsError(resp.Error(), couchdb.ErrNotFound))
}

// TestDeleteDocument tests deleting a document.
func TestDeleteDocument(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	cdb, cleanup := prepareDatabase(assert, "tmp-delete-document")
	defer cleanup()

	// Create test document.
	docA := MyDocument{
		DocumentID: "foo-12345",
		Name:       "foo",
		Age:        33,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	assert.Equal(id, "foo-12345")

	// Read test document, we need it including the revision.
	resp = cdb.ReadDocument(id)
	assert.True(resp.IsOK())
	docB := MyDocument{}
	err := resp.Document(&docB)
	assert.Nil(err)

	// Delete the test document.
	resp = cdb.DeleteDocument(docB)
	assert.True(resp.IsOK())

	// Try to read deleted document.
	resp = cdb.ReadDocument(id)
	assert.False(resp.IsOK())
	assert.Equal(resp.StatusCode(), couchdb.StatusNotFound)

	// Try to delete it a second time.
	resp = cdb.DeleteDocument(docB)
	assert.False(resp.IsOK())
	assert.Equal(resp.StatusCode(), couchdb.StatusNotFound)
	assert.True(errors.IsError(resp.Error(), couchdb.ErrNotFound))
}

// TestDeleteDocumentByID tests deleting a document by identifier.
func TestDeleteDocumentByID(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	cdb, cleanup := prepareDatabase(assert, "tmp-delete-document-by-id")
	defer cleanup()

	// Create test document.
	docA := MyDocument{
		DocumentID: "foo-12345",
		Name:       "foo",
		Age:        33,
	}
	resp := cdb.CreateDocument(docA)
	assert.True(resp.IsOK())
	id := resp.ID()
	revision := resp.Revision()
	assert.Equal(id, "foo-12345")

	// Delete the test document by ID.
	resp = cdb.DeleteDocumentByID(id, revision)
	assert.True(resp.IsOK())

	// Try to read deleted document.
	resp = cdb.ReadDocument(id)
	assert.False(resp.IsOK())
	assert.Equal(resp.StatusCode(), couchdb.StatusNotFound)

	// Try to delete it a second time.
	resp = cdb.DeleteDocumentByID(id, revision)
	assert.False(resp.IsOK())
	assert.Equal(resp.StatusCode(), couchdb.StatusNotFound)
	assert.True(errors.IsError(resp.Error(), couchdb.ErrNotFound))
}

//--------------------
// HELPERS
//--------------------

// MyDocument is used for the tests.
type MyDocument struct {
	DocumentID       string `json:"_id,omitempty"`
	DocumentRevision string `json:"_rev,omitempty"`

	Name        string `json:"name"`
	Age         int    `json:"age"`
	Active      bool   `json:"active"`
	Description string `json:"description"`
}

// prepareDatabase opens the database, deletes a possible test
// database, and creates it newly.
func prepareDatabase(assert *asserts.Asserts, name string) (*couchdb.Database, func()) {
	logger.SetLevel(logger.LevelDebug)
	cdb, err := couchdb.Open(couchdb.Name(name))
	assert.Nil(err)
	rs := cdb.Manager().DeleteDatabase()
	rs = cdb.Manager().CreateDatabase()
	assert.Nil(rs.Error())
	assert.True(rs.IsOK())
	return cdb, func() { cdb.Manager().DeleteDatabase() }
}

// prepareDeletedDatabase opens the database, checks result,
// and deletes it. Ensures a good environment including
// cleanup func.
func prepareDeletedDatabase(assert *asserts.Asserts, name string) (*couchdb.Database, func()) {
	logger.SetLevel(logger.LevelDebug)
	cdb, err := couchdb.Open(couchdb.Name(name))
	assert.Nil(err)
	cdb.Manager().DeleteDatabase()
	cdb.Manager().DeleteNamedDatabase("_users")
	return cdb, func() {
		cdb.Manager().DeleteDatabase()
		cdb.Manager().DeleteNamedDatabase("_users")
	}
}

// prepareFilledDatabase opens the database, deletes a possible test
// database, creates it newly and adds some data.
func prepareFilledDatabase(assert *asserts.Asserts, name string) (*couchdb.Database, func()) {
	logger.SetLevel(logger.LevelDebug)
	cdb, err := couchdb.Open(couchdb.Name(name))
	assert.Nil(err)
	rs := cdb.Manager().DeleteDatabase()
	rs = cdb.Manager().CreateDatabase()
	assert.Nil(rs.Error())
	assert.True(rs.IsOK())

	gen := generators.New(generators.FixedRand())
	docs := []interface{}{}
	for i := 0; i < 1000; i++ {
		first, middle, last := gen.Name()
		doc := MyDocument{
			DocumentID:  identifier.Identifier(last, first, i),
			Name:        first + " " + middle + " " + last,
			Age:         gen.Int(18, 65),
			Active:      gen.FlipCoin(75),
			Description: gen.Sentence(),
		}
		docs = append(docs, doc)
	}
	results, err := cdb.BulkWriteDocuments(docs)
	assert.Nil(err)
	for _, result := range results {
		assert.True(result.OK)
	}

	return cdb, func() { cdb.Manager().DeleteDatabase() }
}

// EOF
