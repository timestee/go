// Tideland Go Library - DB - CouchDB Client - Core
//
// Copyright (C) 2016-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package db

//--------------------
// EXTERNAL DOCUMENT TYPES
//--------------------

// Status contains internal status information CouchDB returns.
type Status struct {
	OK       bool   `json:"ok"`
	ID       string `json:"id"`
	Revision string `json:"rev"`
	Error    string `json:"error"`
	Reason   string `json:"reason"`
}

// Statuses is the list of status information after a bulk writing.
type Statuses []Status

//--------------------
// INTERNAL DOCUMENT TYPES
//--------------------

// couchdbBulkDocuments contains a number of documents added at once.
type couchdbBulkDocuments struct {
	Docs     []interface{} `json:"docs"`
	NewEdits bool          `json:"new_edits,omitempty"`
}

// couchdbDocument is used to simply retrieve ID and revision of
// a document.
type couchdbDocument struct {
	ID       string `json:"_id"`
	Revision string `json:"_rev"`
	Deleted  bool   `json:"_deleted"`
}

// couchdbRows returns rows containing IDs of documents. It's
// part of a view document.
type couchdbRows struct {
	Rows []struct {
		ID string `json:"id"`
	}
}

// EOF
