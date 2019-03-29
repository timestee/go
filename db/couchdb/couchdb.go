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
	"reflect"
	"strings"
	"sync"

	"tideland.dev/go/dsa/identifier"
	"tideland.dev/go/trace/errors"
)

//--------------------
// DATABASE
//--------------------

// Database provides the access to a database.
type Database struct {
	mu      sync.Mutex
	host    string
	name    string
	logging bool
}

// Open returns a configured connection to a CouchDB server.
// Permanent parameters, e.g. for authentication, are possible.
func Open(options ...Option) (*Database, error) {
	db := &Database{
		host:    defaultHost,
		name:    defaultName,
		logging: defaultLogging,
	}
	for _, option := range options {
		if err := option(db); err != nil {
			return nil, err
		}
	}
	return db, nil
}

// Manager returns the database system manager.
func (db *Database) Manager() *Manager {
	db.mu.Lock()
	defer db.mu.Unlock()
	return newManager(db)
}

// Designs returns the design document manager.
func (db *Database) Designs() *Designs {
	db.mu.Lock()
	defer db.mu.Unlock()
	return newDesigns(db)
}

// StartSession starts a cookie based session for the given user.
func (db *Database) StartSession(name, password string) (*Session, error) {
	user := User{
		Name:     name,
		Password: password,
	}
	rs := db.post("_session", user)
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	roles := couchdbRoles{}
	err := rs.Document(&roles)
	if err != nil {
		return nil, err
	}
	setCookie := rs.Header("Set-Cookie")
	authSession := ""
	for _, part := range strings.Split(setCookie, ";") {
		if strings.HasPrefix(part, "AuthSession=") {
			authSession = part
			break
		}
	}
	s := &Session{
		db:          db,
		name:        roles.Name,
		authSession: authSession,
	}
	return s, nil
}

// AllDocumentIDs returns a list of all document IDs
// of the configured database.
func (db *Database) AllDocumentIDs(params ...Parameter) ([]string, error) {
	rs := db.get(db.databasePath("_all_docs"), nil, params...)
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
func (db *Database) HasDocument(id string, params ...Parameter) (bool, error) {
	rs := db.head(db.databasePath(id), nil, params...)
	if rs.IsOK() {
		return true, nil
	}
	if rs.StatusCode() == StatusNotFound {
		return false, nil
	}
	return false, rs.Error()
}

// CreateDocument creates a new document.
func (db *Database) CreateDocument(doc interface{}, params ...Parameter) *ResultSet {
	id, _, err := db.idAndRevision(doc)
	if err != nil {
		return newResultSet(nil, err)
	}
	if id == "" {
		id = identifier.NewUUID().ShortString()
	}
	return db.put(db.databasePath(id), doc, params...)
}

// ReadDocument reads the a document by ID.
func (db *Database) ReadDocument(id string, params ...Parameter) *ResultSet {
	return db.get(db.databasePath(id), nil, params...)
}

// UpdateDocument update a document if exists.
func (db *Database) UpdateDocument(doc interface{}, params ...Parameter) *ResultSet {
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
	return db.put(db.databasePath(id), doc, params...)
}

// DeleteDocument deletes a existing document.
func (db *Database) DeleteDocument(doc interface{}, params ...Parameter) *ResultSet {
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
	return db.delete(db.databasePath(id), nil, params...)
}

// DeleteDocumentByID deletes an existing document simply by
// its identifier and revision.
func (db *Database) DeleteDocumentByID(id, revision string, params ...Parameter) *ResultSet {
	hasDoc, err := db.HasDocument(id)
	if err != nil {
		return newResultSet(nil, err)
	}
	if !hasDoc {
		return newResultSet(nil, errors.New(ErrNotFound, msgNotFound, id))
	}
	params = append(params, Revision(revision))
	return db.delete(db.databasePath(id), nil, params...)
}

// BulkWriteDocuments allows to create or update many
// documents en bloc.
func (db *Database) BulkWriteDocuments(docs []interface{}, params ...Parameter) (Statuses, error) {
	bulk := &couchdbBulkDocuments{
		Docs: docs,
	}
	rs := db.post(db.databasePath("_bulk_docs"), bulk, params...)
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

// path creates a document path starting at root.
func (db *Database) path(parts ...string) string {
	return strings.Join(append([]string{""}, parts...), "/")
}

// databasePath creates a document path for the database.
func (db *Database) databasePath(parts ...string) string {
	return db.path(append([]string{db.name}, parts...)...)
}

// head performs a HEAD request against the database.
func (db *Database) head(path string, doc interface{}, params ...Parameter) *ResultSet {
	req := newRequest(db, path, doc)
	return req.apply(params...).head()
}

// get performs a GET request against the database.
func (db *Database) get(path string, doc interface{}, params ...Parameter) *ResultSet {
	req := newRequest(db, path, doc)
	return req.apply(params...).get()
}

// put performs a PUT request against the database.
func (db *Database) put(path string, doc interface{}, params ...Parameter) *ResultSet {
	req := newRequest(db, path, doc)
	return req.apply(params...).put()
}

// Post performs a POST request against the database.
func (db *Database) post(path string, doc interface{}, params ...Parameter) *ResultSet {
	req := newRequest(db, path, doc)
	return req.apply(params...).post()
}

// delete performs a DELETE request against the database.
func (db *Database) delete(path string, doc interface{}, params ...Parameter) *ResultSet {
	req := newRequest(db, path, doc)
	return req.apply(params...).delete()
}

// getOrPost decides based on the document if it will perform
// a GET request or a POST request. The document can be set directly
// or by one of the parameters. Several of the CouchDB commands
// work this way.
func (db *Database) getOrPost(path string, doc interface{}, params ...Parameter) *ResultSet {
	var rs *ResultSet
	req := newRequest(db, path, doc).apply(params...)
	if req.doc != nil {
		rs = req.post()
	} else {
		rs = req.get()
	}
	return rs
}

// idAndRevision retrieves the ID and the revision of the
// passed document.
func (db *Database) idAndRevision(doc interface{}) (string, string, error) {
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
