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
)

//--------------------
// RESPONSE TOOLS
//--------------------

// MarshalResponseBody allows to directly marshal a value into the content
// types JSON or XML and write it to the passed response writer. The content
// type header is set too.
func MarshalResponseBody(w http.ResponseWriter, contentType string, v interface{}) error {
	switch contentType {
	case ContentTypeJSON:
		body, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("webbox: cannot marshal response body: %v", err)
		}
		w.Header().Set(HeaderContentType, contentType)
		_, err = w.Write(body)
		if err != nil {
			return fmt.Errorf("webbox: cannot write  body: %v", err)
		}
		return nil
	case ContentTypeXML:
		body, err := xml.Marshal(v)
		if err != nil {
			return fmt.Errorf("webbox: cannot marshal response body: %v", err)
		}
		w.Header().Set(HeaderContentType, contentType)
		_, err = w.Write(body)
		if err != nil {
			return fmt.Errorf("webbox: cannot write  body: %v", err)
		}
		return nil
	default:
		return fmt.Errorf("webbox: invalid content-type")
	}

}

// EOF
