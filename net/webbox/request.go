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

// PathField returns the nth field of the request path and true if it exists.
// Otherwise an empty string and false.
func PathField(r *http.Request, n int) (string, bool) {
	if n < 1 {
		panic("webbox: illegal path index")
	}
	fields := strings.Split(r.URL.Path, "/")
	// Empty string before slash is field zero.
	if len(fields)-1 < n {
		// Does not exist.
		return "", false
	}
	return fields[n], true
}

// EOF
