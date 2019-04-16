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
	"fmt"
	"io/ioutil"
	"net/http"
)

//--------------------
// CONSTANTS
//--------------------

const (
	HeaderAccept      = "Accept"
	HeaderContentType = "Content-Type"

	ContentTypePlain      = "text/plain"
	ContentTypeHTML       = "text/html"
	ContentTypeXML        = "application/xml"
	ContentTypeJSON       = "application/json"
	ContentTypeURLEncoded = "application/x-www-form-urlencoded"
)

//--------------------
// PRIVATE HELPERS
//--------------------

// readBody retrieves the whole body out of a HTTP request and
// returns it as byte slice.
func readBody(r *http.Request) ([]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("webbox: cannot read body: %v", err)
	}
	r.Body.Close()
	return body, nil
}

// EOF
