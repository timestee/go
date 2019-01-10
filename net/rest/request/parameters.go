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
	"io"
	"net/url"

	"tideland.one/go/dsa/version"
	"tideland.one/go/net/jwt/token"
	"tideland.one/go/net/rest/core"
	"tideland.one/go/trace/errors"
)

//--------------------
// CALL PARAMETERS
//--------------------

// Parameters allows to pass parameters to a call.
type Parameters struct {
	Version     *version.Version
	Token       *token.JWT
	ContentType string
	Content     interface{}
	Accept      string
}

// body returns the content as body data depending on
// the content type.
func (p *Parameters) body() (io.Reader, error) {
	buffer := bytes.NewBuffer(nil)
	if p.Content == nil {
		return buffer, nil
	}
	// Process content based on content type.
	switch p.ContentType {
	case core.ContentTypeXML:
		tmp, err := xml.Marshal(p.Content)
		if err != nil {
			return nil, errors.Annotate(err, ErrProcessingRequestContent, msgProcessingRequestContent)
		}
		buffer.Write(tmp)
	case core.ContentTypeJSON:
		tmp, err := json.Marshal(p.Content)
		if err != nil {
			return nil, errors.Annotate(err, ErrProcessingRequestContent, msgProcessingRequestContent)
		}
		buffer.Write(tmp)
	case core.ContentTypeGOB:
		enc := gob.NewEncoder(buffer)
		if err := enc.Encode(p.Content); err != nil {
			return nil, errors.Annotate(err, ErrProcessingRequestContent, msgProcessingRequestContent)
		}
	case core.ContentTypeURLEncoded:
		values, err := p.values()
		if err != nil {
			return nil, err
		}
		_, err = buffer.WriteString(values.Encode())
		if err != nil {
			return nil, errors.Annotate(err, ErrProcessingRequestContent, msgProcessingRequestContent)
		}
	}
	return buffer, nil
}

// values returns the content as URL encoded values.
func (p *Parameters) values() (url.Values, error) {
	if p.Content == nil {
		return url.Values{}, nil
	}
	// Check if type is already ok.
	urlvs, ok := p.Content.(url.Values)
	if ok {
		return urlvs, nil
	}
	// Check for simple key/values.
	kvs, ok := p.Content.(KeyValues)
	if !ok {
		return nil, errors.New(ErrInvalidContent, msgInvalidContent)
	}
	values := url.Values{}
	for key, value := range kvs {
		values.Set(key, value)
	}
	return values, nil
}

// EOF
