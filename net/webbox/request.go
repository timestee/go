// Tideland Go Library - Network - Web Toolbox
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package webbox

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"strings"
)

//--------------------
// REQUEST TOOLS
//--------------------

// PathFields splits the request path into its fields.
func PathFields(r *http.Request) []string {
	rawFields := strings.Split(r.URL.Path, "/")
	fields := []string{}
	for _, field := range rawFields {
		if field != "" {
			fields = append(fields, field)
		}
	}
	return fields
}

// PathField returns the nth field of the request path and true if it exists.
// Otherwise an empty string and false.
func PathField(r *http.Request, n int) (string, bool) {
	if n < 0 {
		panic("webbox: illegal path index")
	}
	fields := PathFields(r)
	if len(fields) < n+1 {
		return "", false
	}
	return fields[n], true
}

// EOF
