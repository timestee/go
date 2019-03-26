// Tideland Go Library - DB - CouchDB Client - Core
//
// Copyright (C) 2016-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package db

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/json"
	"reflect"
	"strings"

	"tideland.dev/go/dsa/identifier"
	"tideland.dev/go/trace/errors"
)

//--------------------
// DATABASE
//--------------------

// DB provides the access to a database.
type DB struct {
	host       string
	database   string
	debugLog   bool
	parameters []Parameter
}

// Open returns a configured connection to a CouchDB server.
// Permanent parameters, e.g. for authentication, are possible.
func Open(host, database string, debugLog bool, params ...Parameter) (*DB, error) {
	db := &DB{
		host:       host,
		database:   database,
		debugLog:   debugLog,
		parameters: params,
	}
	return db, nil
}

// Path creates a document path starting at root.
func (db *DB) Path(parts ...string) string {
	return strings.Join(append([]string{""}, parts...), "/")
}

// DatabasePath creates a document path for the database.
func (db *DB) DatabasePath(parts ...string) string {
	return db.Path(append([]string{db.database}, parts...)...)
}

// Head performs a HEAD request against the configured database.
func (db *DB) Head(path string, doc interface{}, params ...Parameter) *ResultSet {
	req := newRequest(db, path, doc)
	return req.apply(params...).head()
}

// Get performs a GET request against the configured database.
func (db *DB) Get(path string, doc interface{}, params ...Parameter) *ResultSet {
	req := newRequest(db, path, doc)
	return req.apply(params...).get()
}

// Put performs a GET request against the configured database.
func (db *DB) Put(path string, doc interface{}, params ...Parameter) *ResultSet {
	req := newRequest(db, path, doc)
	return req.apply(params...).put()
}

// Post performs a GET request against the configured database.
func (db *DB) Post(path string, doc interface{}, params ...Parameter) *ResultSet {
	req := newRequest(db, path, doc)
	return req.apply(params...).post()
}

// Delete performs a GET request against the configured database.
func (db *DB) Delete(path string, doc interface{}, params ...Parameter) *ResultSet {
	req := newRequest(db, path, doc)
	return req.apply(params...).delete()
}

// GetOrPost decides based on the document if it will perform
// a GET request or a POST request. The document can be set directly
// or by one of the parameters. Several of the CouchDB commands
// work this way.
func (db *DB) GetOrPost(path string, doc interface{}, params ...Parameter) *ResultSet {
	var rs *ResultSet
	req := newRequest(db, path, doc).apply(params...)
	if req.doc != nil {
		rs = req.post()
	} else {
		rs = req.get()
	}
	return rs
}

// Version returns the version number of the database instance.
func (db *DB) Version() (string, error) {
	rs := db.Get("/", nil)
	if !rs.IsOK() {
		return "", rs.Error()
	}
	welcome := map[string]interface{}{}
	err := rs.Document(&welcome)
	if err != nil {
		return "", err
	}
	version, ok := welcome["version"].(string)
	if !ok {
		return "", errors.New(ErrInvalidVersion, msgInvalidVersion, welcome["version"])
	}
	return version, nil
}

// AllDatabases returns a list of all database IDs
// of the connected server.
func (db *DB) AllDatabases() ([]string, error) {
	rs := db.Get("/_all_dbs", nil)
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	ids := []string{}
	err := rs.Document(&ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

// HasDatabase checks if the configured database exists.
func (db *DB) HasDatabase() (bool, error) {
	rs := db.Head(db.DatabasePath(), nil)
	if rs.IsOK() {
		return true, nil
	}
	if rs.StatusCode() == StatusNotFound {
		return false, nil
	}
	return false, rs.Error()
}

// CreateDatabase creates the configured database.
func (db *DB) CreateDatabase(params ...Parameter) *ResultSet {
	return db.Put(db.DatabasePath(), nil, params...)
}

// DeleteDatabase removes the configured database.
func (db *DB) DeleteDatabase(params ...Parameter) *ResultSet {
	return db.Delete(db.DatabasePath(), nil, params...)
}

// AllDesigns returns the list of all design
// document IDs of the configured database.
func (db *DB) AllDesigns() ([]string, error) {
	jstart, _ := json.Marshal("_design/")
	jend, _ := json.Marshal("_design0/")
	startEndKey := Query(KeyValue{"startkey", string(jstart)}, KeyValue{"endkey", string(jend)})
	rs := db.Get(db.DatabasePath("_all_docs"), nil, startEndKey)
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	designRows := couchdbRows{}
	err := rs.Document(&designRows)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, row := range designRows.Rows {
		ids = append(ids, row.ID)
	}
	return ids, nil
}

// Design returns the design document instance for the given ID.
func (db *DB) Design(id string) (*Design, error) {
	return newDesign(db, id)
}

// AllDocuments returns a list of all document IDs
// of the configured database.
func (db *DB) AllDocuments() ([]string, error) {
	rs := db.Get(db.DatabasePath("_all_docs"), nil)
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	designRows := couchdbRows{}
	err := rs.Document(&designRows)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	for _, row := range designRows.Rows {
		ids = append(ids, row.ID)
	}
	return ids, nil
}

// HasDocument checks if the document with the ID exists.
func (db *DB) HasDocument(id string) (bool, error) {
	rs := db.Head(db.DatabasePath(id), nil)
	if rs.IsOK() {
		return true, nil
	}
	if rs.StatusCode() == StatusNotFound {
		return false, nil
	}
	return false, rs.Error()
}

// CreateDocument creates a new document.
func (db *DB) CreateDocument(doc interface{}, params ...Parameter) *ResultSet {
	id, _, err := db.idAndRevision(doc)
	if err != nil {
		return newResultSet(nil, err)
	}
	if id == "" {
		id = identifier.NewUUID().ShortString()
	}
	return db.Put(db.DatabasePath(id), doc, params...)
}

// ReadDocument reads an existing document.
func (db *DB) ReadDocument(id string, params ...Parameter) *ResultSet {
	return db.Get(db.DatabasePath(id), nil, params...)
}

// UpdateDocument update an existing document.
func (db *DB) UpdateDocument(doc interface{}, params ...Parameter) *ResultSet {
	id, _, err := db.idAndRevision(doc)
	if err != nil {
		return newResultSet(nil, err)
	}
	if id == "" {
		return newResultSet(nil, errors.New(ErrNoIdentifier, msgNoIdentifier))
	}
	hasDoc, err := db.HasDocument(id)
	if err != nil {
		return newResultSet(nil, err)
	}
	if !hasDoc {
		return newResultSet(nil, errors.New(ErrNotFound, msgNotFound, id))
	}
	return db.Put(db.DatabasePath(id), doc, params...)
}

// DeleteDocument deletes an existing document.
func (db *DB) DeleteDocument(doc interface{}, params ...Parameter) *ResultSet {
	id, revision, err := db.idAndRevision(doc)
	if err != nil {
		return newResultSet(nil, err)
	}
	hasDoc, err := db.HasDocument(id)
	if err != nil {
		return newResultSet(nil, err)
	}
	if !hasDoc {
		return newResultSet(nil, errors.New(ErrNotFound, msgNotFound, id))
	}
	params = append(params, Revision(revision))
	return db.Delete(db.DatabasePath(id), nil, params...)
}

// DeleteDocumentByID deletes an existing document simply by
// its identifier and revision.
func (db *DB) DeleteDocumentByID(id, revision string, params ...Parameter) *ResultSet {
	hasDoc, err := db.HasDocument(id)
	if err != nil {
		return newResultSet(nil, err)
	}
	if !hasDoc {
		return newResultSet(nil, errors.New(ErrNotFound, msgNotFound, id))
	}
	params = append(params, Revision(revision))
	return db.Delete(db.DatabasePath(id), nil, params...)
}

// BulkWriteDocuments allows to create or update many
// documents en bloc.
func (db *DB) BulkWriteDocuments(docs []interface{}, params ...Parameter) (Statuses, error) {
	bulk := &couchdbBulkDocuments{
		Docs: docs,
	}
	rs := db.Post(db.DatabasePath("_bulk_docs"), bulk, params...)
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	statuses := Statuses{}
	err := rs.Document(&statuses)
	if err != nil {
		return nil, err
	}
	return statuses, nil
}

// idAndRevision retrieves the ID and the revision of the
// passed document.
func (db *DB) idAndRevision(doc interface{}) (string, string, error) {
	v := reflect.Indirect(reflect.ValueOf(doc))
	t := v.Type()
	k := t.Kind()
	if k != reflect.Struct {
		return "", "", errors.New(ErrInvalidDocument, msgInvalidDocument)
	}
	var id string
	var revision string
	var found int
	for i := 0; i < t.NumField(); i++ {
		vf := v.Field(i)
		tf := t.Field(i)
		if json, ok := tf.Tag.Lookup("json"); ok {
			switch json {
			case "_id", "_id,omitempty":
				id = vf.String()
				found++
			case "_rev", "_rev,omitempty":
				revision = vf.String()
				found++
			}
		}
	}
	if found != 2 {
		return "", "", errors.New(ErrInvalidDocument, msgInvalidDocument)
	}
	return id, revision, nil
}

// EOF
