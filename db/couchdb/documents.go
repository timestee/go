// Tideland Go Library - DB - CouchDB Client
//
// Copyright (C) 2016-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package couchdb

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

// User contains name and password for
// user management and authentication.
type User struct {
	DocumentID       string `json:"_id,omitempty"`
	DocumentRevision string `json:"_rev,omitempty"`

	Name     string   `json:"name"`
	Password string   `json:"password"`
	Type     string   `json:"type,omitempty"`
	Roles    []string `json:"roles,omitempty"`
}

// NamesRoles contains names and roles for
// administrators and users.
type NamesRoles struct {
	Names []string `json:"names,omitempty"`
	Roles []string `json:"roles,omitempty"`
}

// Security contains administrators and
// members for one database.
type Security struct {
	Admins  NamesRoles `json:"admins,omitempty"`
	Members NamesRoles `json:"members,omitempty"`
}

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

// couchdRoles contains the roles of a user if the
// authentication succeeded.
type couchdbRoles struct {
	OK       bool     `json:"ok"`
	Name     string   `json:"name"`
	Password string   `json:"password_sha,omitempty"`
	Salt     string   `json:"salt,omitempty"`
	Type     string   `json:"type"`
	Roles    []string `json:"roles"`
}

// EOF
