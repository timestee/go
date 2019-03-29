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
	"tideland.dev/go/dsa/version"
	"tideland.dev/go/trace/errors"
)

//--------------------
// MANAGER
//--------------------

// Manager bundles the methods to manage the database system
// opposite to handle documents.
type Manager struct {
	db *Database
}

// newManager creates the manager instance.
func newManager(db *Database) *Manager {
	return &Manager{
		db: db,
	}
}

// Version returns the version number of the database instance.
func (m *Manager) Version() (version.Version, error) {
	rs := m.db.get("/", nil)
	if !rs.IsOK() {
		return version.New(0, 0, 0), rs.Error()
	}
	welcome := map[string]interface{}{}
	err := rs.Document(&welcome)
	if err != nil {
		return version.New(0, 0, 0), err
	}
	vsn, ok := welcome["version"].(string)
	if !ok {
		return version.New(0, 0, 0), errors.New(ErrInvalidVersion, msgInvalidVersion, welcome["version"])
	}
	return version.Parse(vsn)
}

// AllDatabaseIDs returns a list of all database IDs
// of the connected server.
func (m *Manager) AllDatabaseIDs() ([]string, error) {
	rs := m.db.get("/_all_dbs", nil)
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
func (m *Manager) HasDatabase() (bool, error) {
	rs := m.db.head(m.db.databasePath(), nil)
	if rs.IsOK() {
		return true, nil
	}
	if rs.StatusCode() == StatusNotFound {
		return false, nil
	}
	return false, rs.Error()
}

// CreateDatabase creates the configured database.
func (m *Manager) CreateDatabase(params ...Parameter) *ResultSet {
	return m.db.put(m.db.databasePath(), nil, params...)
}

// DeleteDatabase removes the configured database.
func (m *Manager) DeleteDatabase(params ...Parameter) *ResultSet {
	return m.db.delete(m.db.databasePath(), nil, params...)
}

// HasAdministrator checks if a given administrator account exists.
func (m *Manager) HasAdministrator(nodename, name string, params ...Parameter) (bool, error) {
	path := m.db.path("_node", nodename, "_config", "admins", name)
	rs := m.db.get(path, nil, params...)
	if !rs.IsOK() {
		if rs.StatusCode() == StatusNotFound {
			return false, nil
		}
		return false, rs.Error()
	}
	return true, nil
}

// WriteAdministrator adds or updates an administrator to the given database.
func (m *Manager) WriteAdministrator(nodename, name, password string, params ...Parameter) error {
	path := m.db.path("_node", nodename, "_config", "admins", name)
	rs := m.db.put(path, password, params...)
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

// DeleteAdministrator deletes an administrator from the given database.
func (m *Manager) DeleteAdministrator(nodename, name string, params ...Parameter) error {
	path := m.db.path("_node", nodename, "_config", "admins", name)
	rs := m.db.delete(path, nil, params...)
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

// CreateUser adds a new user to the system.
func (m *Manager) CreateUser(user *User, params ...Parameter) error {
	if err := ensureUsersDatabase(m.db); err != nil {
		return err
	}
	user.DocumentID = userDocumentID(user.Name)
	user.Type = "user"
	path := m.db.path("_users", user.DocumentID)
	rs := m.db.put(path, user, params...)
	return rs.Error()
}

// ReadUser reads an existing user from the system.
func (m *Manager) ReadUser(name string, params ...Parameter) (*User, error) {
	path := m.db.path("_users", userDocumentID(name))
	rs := m.db.get(path, nil, params...)
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	var user User
	err := rs.Document(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates a user in the system.
func (m *Manager) UpdateUser(user *User, params ...Parameter) error {
	if err := ensureUsersDatabase(m.db); err != nil {
		return err
	}
	path := m.db.path("_users", user.DocumentID)
	rs := m.db.put(path, user, params...)
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

// DeleteUser deletes a user from the system.
func (m *Manager) DeleteUser(user *User, params ...Parameter) error {
	params = append(params, Revision(user.DocumentRevision))
	path := m.db.path("_users", user.DocumentID)
	rs := m.db.delete(path, nil, params...)
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

// ReadSecurity returns the security for the given database.
func (m *Manager) ReadSecurity(params ...Parameter) (*Security, error) {
	path := m.db.databasePath("_security")
	rs := m.db.get(path, nil, params...)
	if !rs.IsOK() {
		return nil, rs.Error()
	}
	var security Security
	err := rs.Document(&security)
	if err != nil {
		return nil, err
	}
	return &security, nil
}

// WriteSecurity writes new or changed security data to
// the given database.
func (m *Manager) WriteSecurity(security Security, params ...Parameter) error {
	path := m.db.databasePath("_security")
	rs := m.db.put(path, security, params...)
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

//--------------------
// HELPERS
//--------------------

// ensureUsersDatabase checks if the _users database exists and
// creates it if needed.
func ensureUsersDatabase(db *Database) error {
	rs := db.get("_users", nil)
	if rs.IsOK() {
		return nil
	}
	return db.put("_users", nil).Error()
}

// userDocumentID builds the document ID based
// on the name.
func userDocumentID(name string) string {
	return "org.couchdb.user:" + name
}

// EOF
