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
}

//--------------------
// FINDS
//--------------------

// FindProcessor is a function processing the content of a found document.
type FindProcessor func(document Unmarshable) error

// Finds allows to find and process documents by a given selector.
type Finds struct {
	db *Database
}

// newFinds returns a new finds instance.
func newFinds(db *Database, selector *Selector, params ...Parameter) *Finds {
	rs := db.Request().SetPath(db.name, "_find").SetDocument(selector).ApplyParameters(params...).Post()
	return &Finds{
		db: db,
	}
}

// EOF
