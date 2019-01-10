// Tideland Go Library - Network - REST - Request
//
// Copyright (C) 2009-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package request

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"net/url"
	"strings"

	"tideland.one/go/net/rest/core"
	"tideland.one/go/trace/errors"
)

//--------------------
// RESPONSE
//--------------------

// KeyValues handles keys and values for request headers and cookies.
type KeyValues map[string]string

// Response wraps all infos of a test response.
type Response struct {
	httpResp    *http.Response
	contentType string
	content     []byte
}

// StatusCode returns the HTTP status code of the response.
func (r *Response) StatusCode() int {
	return r.httpResp.StatusCode
}

// Header returns the HTTP header of the response.
func (r *Response) Header() http.Header {
	return r.httpResp.Header
}

// HasContentType checks the content type regardless of charsets.
func (r *Response) HasContentType(contentType string) bool {
	return strings.Contains(r.contentType, contentType)
}

// Read decodes the content into the passed data depending
// on the content type.
func (r *Response) Read(data interface{}) error {
	switch {
	case r.HasContentType(core.ContentTypeGOB):
		dec := gob.NewDecoder(bytes.NewBuffer(r.content))
		if err := dec.Decode(data); err != nil {
			return errors.Annotate(err, ErrDecodingResponse, msgDecodingResponse)
		}
		return nil
	case r.HasContentType(core.ContentTypeJSON):
		if err := json.Unmarshal(r.content, &data); err != nil {
			return errors.Annotate(err, ErrDecodingResponse, msgDecodingResponse)
		}
		return nil
	case r.HasContentType(core.ContentTypeXML):
		if err := xml.Unmarshal(r.content, &data); err != nil {
			return errors.Annotate(err, ErrDecodingResponse, msgDecodingResponse)
		}
		return nil
	case r.HasContentType(core.ContentTypeURLEncoded):
		values, err := url.ParseQuery(string(r.content))
		if err != nil {
			return errors.Annotate(err, ErrDecodingResponse, msgDecodingResponse)
		}
		// Check for data type url.Values.
		duv, ok := data.(url.Values)
		if ok {
			for key, value := range values {
				duv[key] = value
			}
			return nil
		}
		// Check for data type KeyValues.
		kvv, ok := data.(KeyValues)
		if !ok {
			return errors.New(ErrDecodingResponse, msgDecodingResponse)
		}
		for key, value := range values {
			kvv[key] = strings.Join(value, " / ")
		}
		return nil
	}
	return errors.New(ErrInvalidContentType, msgInvalidContentType, r.contentType)
}

// ReadFeedback tries to unmarshal the content of the
// response into a rest package feedback.
func (r *Response) ReadFeedback() (core.Feedback, bool) {
	fb := core.Feedback{}
	err := r.Read(&fb)
	if err != nil {
		return core.Feedback{
			StatusCode: -1,
			Status:     "fail",
			Message:    err.Error(),
			Payload:    r.content,
		}, false
	}
	return fb, true
}

// EOF
