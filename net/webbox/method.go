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
	http.MethodGet:     struct{}{},
	http.MethodHead:    struct{}{},
	http.MethodPost:    struct{}{},
	http.MethodPut:     struct{}{},
	http.MethodPatch:   struct{}{},
	http.MethodDelete:  struct{}{},
	http.MethodConnect: struct{}{},
	http.MethodOptions: struct{}{},
	http.MethodTrace:   struct{}{},
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
