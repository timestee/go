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
	"fmt"
)

//--------------------
// SESSION
//--------------------

// Session contains the information of a CouchDB session.
type Session struct {
	db          *Database
	name        string
	authSession string
}

// Name returns the users name of this session.
func (s *Session) Name() string {
	return s.name
}

// Cookie returns the session cookie as parameter
// to be used in the individual database requests.
func (s *Session) Cookie() Parameter {
	return func(pa Parameterizable) {
		pa.SetHeader("X-CouchDB-WWW-Authenticate", "Cookie")
		pa.SetHeader("Cookie", s.authSession)
	}
}

// Stop ends the session.
func (s *Session) Stop() error {
	rs := s.db.delete("/_session", nil, s.Cookie())
	if !rs.IsOK() {
		return rs.Error()
	}
	return nil
}

// String returns a string representation of the session.
func (s *Session) String() string {
	return fmt.Sprintf("[DB: %q USER: %q SESSION: %q]", s.db.name, s.name, s.authSession)
}

// EOF
