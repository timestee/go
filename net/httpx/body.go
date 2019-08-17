// Tideland Go Library - Network - HTTP Extensions
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package httpx // import "tideland.dev/go/net/httpx"

//--------------------
// IMPORTS
//--------------------

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"

	"tideland.dev/go/trace/failure"
)

//--------------------
// BODY TOOLS
//--------------------

// ReadBody retrieves the whole body out of a HTTP request or response
// and returns it as byte slice.
func ReadBody(body io.ReadCloser) ([]byte, error) {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, failure.Annotate(err, "webbox: cannot read body")
	}
	body.Close()
	return data, nil
}

// UnmarshalBody parses the body data of a request or response based on the
// content type header stores the result in the value pointed by v. Conten types
// JSON and XML expect the according types as arguments, all text types
// expect a string, and all others too, but the data is encoded in BASE64.
func UnmarshalBody(body io.ReadCloser, h http.Header, v interface{}) error {
	data, err := ReadBody(body)
	if err != nil {
		return err
	}
	switch {
	case ContainsContentType(h, contentTypesText):
		switch tv := v.(type) {
		case *string:
			*tv = string(data)
		case *interface{}:
			*tv = string(data)
		default:
			return failure.New("invalid value argument for text or HTML body;,want string or empty interface")
		}
		return nil
	case ContainsContentType(h, ContentTypeJSON):
		if err = json.Unmarshal(data, &v); err != nil {
			return failure.Annotate(err, "cannot unmarshal JSON body")
		}
		return nil
	case ContainsContentType(h, ContentTypeXML):
		if err = xml.Unmarshal(data, &v); err != nil {
			return failure.Annotate(err, "cannot unmarshal XML body")
		}
		return nil
	default:
		sd := base64.StdEncoding.EncodeToString(data)
		switch tv := v.(type) {
		case *string:
			*tv = sd
		case *interface{}:
			*tv = sd
		default:
			return failure.New("invalid value argument for text or HTML body;,want string or empty interface")
		}
		return nil
	}
}

// EOF
