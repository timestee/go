// Tideland Go Library - Network - Web Toolbox
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package webbox // import "tideland.dev/go/net/webbox"

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
)

//--------------------
// CONSTANTS
//--------------------

// httpMethods contains all HTTP methods.
var httpMethods = map[string]struct{}{
	http.MethodGet:     {},
	http.MethodHead:    {},
	http.MethodPost:    {},
	http.MethodPut:     {},
	http.MethodPatch:   {},
	http.MethodDelete:  {},
	http.MethodConnect: {},
	http.MethodOptions: {},
	http.MethodTrace:   {},
}

//--------------------
// METHOD RELATED FUNCTIONS
//--------------------

// ValidMethod returns true if the passed method is valid.
func ValidMethod(method string) bool {
	_, valid := httpMethods[method]
	return valid
}

// EOF
