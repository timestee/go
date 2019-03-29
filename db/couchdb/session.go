// Tideland Go Library - DB - CouchDB Client
//
// Copyright (C) 2016-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package couchdb

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
		pa.SetHeader("Cookie", s.authSession)
		pa.SetHeader("X-CouchDB-WWW-Authenticate", "Cookie")
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

// EOF
