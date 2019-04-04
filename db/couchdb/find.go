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

//--------------------
// SELECTOR
//--------------------

// Selector allows to formulate what documents shall be selected.
type Selector struct {
	parameters map[string]interface{}
}

// NewSelector creates a selector for the search of documents.
func NewSelector() *Selector {
	s := &Selector{
		parameters: make(map[string]interface{}),
	}
	return s
}

// MarshalJSON implements json.Marshaler.
func (s *Selector) MarshalJSON() ([]byte, error) {
	return nil, nil
}

//--------------------
// FINDS
//--------------------

// FindProcessor is a function processing the content of a found document.
type FindProcessor func(document *Unmarshable) error

// Find allows to find and process documents by a given selector.
type Find struct {
	db   *Database
	find *couchdbFind
}

// newFind returns a new finds instance.
func newFind(db *Database, selector *Selector, params ...Parameter) (*Find, error) {
	rs := db.Request().SetPath(db.name, "_find").SetDocument(selector).ApplyParameters(params...).Post()
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	find := couchdbFind{}
	err := rs.Document(&find)
	if err != nil {
		return nil, err
	}
	return &Find{
		db:   db,
		find: &find,
	}, nil
}

// Len returns the number of found documents.
func (f *Find) Len() int {
	return len(f.find.Documents)
}

// Process iterates over the found documents and processes them.
func (f *Find) Process(process FindProcessor) error {
	for _, doc := range f.find.Documents {
		unmarshableDoc := NewUnmarshableJSON(doc)
		if err := process(unmarshableDoc); err != nil {
			return err
		}
	}
	return nil
}

// EOF
