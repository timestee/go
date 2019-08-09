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
	"encoding/json"
	"encoding/xml"
	"fmt"
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

// PathAt returns the nth part of the given request path and true if
// it exists. Otherwise an empty string and false. This way users of
// the nested handlers can retrieve an entity ID out of the path.
func PathAt(p string, n int) (string, bool) {
	if n < 0 {
		panic("webbox: illegal path index")
	}
	parts := strings.Split(p[1:], "/")
	if len(parts) < n+1 || parts[n] == "" {
		return "", false
	}
	return parts[n], true
}

// AcceptsContentType checks if the requestor accepts a given content type.
func AcceptsContentType(r *http.Request, contentType string) bool {
	return strings.Contains(r.Header.Get(HeaderAccept), contentType)
}

// HasContentType checks if the requestor has a given content type.
func HasContentType(r *http.Request, contentType string) bool {
	return strings.Contains(r.Header.Get(HeaderContentType), contentType)
}

// UnmarshalRequestBody parses the body data of a request based on the
// content type header stores the result in the value pointed by v.
// Currently JSON and XML are supported.
func UnmarshalRequestBody(r *http.Request, v interface{}) error {
	switch {
	case HasContentType(r, ContentTypeJSON):
		body, err := readBody(r)
		if err != nil {
			return err
		}
		err = json.Unmarshal(body, &v)
		if err != nil {
			return fmt.Errorf("webbox: cannot unmarshal request body: %v", err)
		}
		return nil
	case HasContentType(r, ContentTypeXML):
		body, err := readBody(r)
		if err != nil {
			return err
		}
		err = xml.Unmarshal(body, &v)
		if err != nil {
			return fmt.Errorf("webbox: cannot unmarshal request body: %v", err)
		}
		return nil
	default:
		return fmt.Errorf("webbox: invalid content-type")
	}
}

// EOF
