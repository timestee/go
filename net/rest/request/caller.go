// Tideland Go Library - Network - REST - Request
//
// Copyright (C) 2009-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package request

//--------------------
// IMPORTS
//--------------------

import (
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"

	"tideland.one/go/net/jwt/token"
	"tideland.one/go/net/rest/core"
	"tideland.one/go/trace/errors"
)

//--------------------
// CALLER
//--------------------

// Caller provides an interface to make calls to configured services.
type Caller struct {
	domain string
	srvs   []*server
}

// newCaller creates a configured Caller.
func newCaller(domain string, srvs []*server) *Caller {
	return &Caller{domain, srvs}
}

// Get performs a GET request on the defined resource.
func (c *Caller) Get(resource, resourceID string, params *Parameters) (*Response, error) {
	return c.request("GET", resource, resourceID, params)
}

// Head performs a HEAD request on the defined resource.
func (c *Caller) Head(resource, resourceID string, params *Parameters) (*Response, error) {
	return c.request("HEAD", resource, resourceID, params)
}

// Put performs a PUT request on the defined resource.
func (c *Caller) Put(resource, resourceID string, params *Parameters) (*Response, error) {
	return c.request("PUT", resource, resourceID, params)
}

// Post performs a POST request on the defined resource.
func (c *Caller) Post(resource, resourceID string, params *Parameters) (*Response, error) {
	return c.request("POST", resource, resourceID, params)
}

// Patch performs a PATCH request on the defined resource.
func (c *Caller) Patch(resource, resourceID string, params *Parameters) (*Response, error) {
	return c.request("PATCH", resource, resourceID, params)
}

// Delete performs a DELETE request on the defined resource.
func (c *Caller) Delete(resource, resourceID string, params *Parameters) (*Response, error) {
	return c.request("DELETE", resource, resourceID, params)
}

// Options performs a OPTIONS request on the defined resource.
func (c *Caller) Options(resource, resourceID string, params *Parameters) (*Response, error) {
	return c.request("OPTIONS", resource, resourceID, params)
}

// request performs all requests.
func (c *Caller) request(method, resource, resourceID string, params *Parameters) (*Response, error) {
	// Preparation.
	client, urlStr, err := c.prepareClient(resource, resourceID)
	if err != nil {
		return nil, err
	}
	request, err := c.prepareRequest(method, urlStr, params)
	if err != nil {
		return nil, err
	}
	// Perform request.
	response, err := client.Do(request)
	if err != nil {
		return nil, errors.Annotate(err, ErrHTTPRequestFailed, msgHTTPRequestFailed)
	}
	// Analyze response.
	return analyzeResponse(response)
}

// prepareClient prepares the client and the URL for the call.
func (c *Caller) prepareClient(resource, resourceID string) (*http.Client, string, error) {
	// TODO Mue 2016-10-28 Add more algorithms than just random selection.
	srv := c.srvs[rand.Intn(len(c.srvs))]
	client := &http.Client{}
	if srv.Transport != nil {
		client.Transport = srv.Transport
	}
	u, err := url.Parse(srv.URL)
	if err != nil {
		return nil, "", errors.Annotate(err, ErrCannotPrepareRequest, msgCannotPrepareRequest)
	}
	upath := strings.Trim(u.Path, "/")
	path := []string{upath, c.domain, resource}
	if resourceID != "" {
		path = append(path, resourceID)
	}
	u.Path = strings.Join(path, "/")
	return client, u.String(), nil
}

// prepareRequest prepares the request to perform.
func (c *Caller) prepareRequest(method, urlStr string, params *Parameters) (*http.Request, error) {
	if params == nil {
		params = &Parameters{}
	}
	var req *http.Request
	var err error
	if method == "GET" || method == "HEAD" {
		// These allow only URL encoded.
		req, err = http.NewRequest(method, urlStr, nil)
		if err != nil {
			return nil, errors.Annotate(err, ErrCannotPrepareRequest, msgCannotPrepareRequest)
		}
		values, err := params.values()
		if err != nil {
			return nil, err
		}
		req.URL.RawQuery = values.Encode()
		req.Header.Set("Content-Type", core.ContentTypeURLEncoded)
	} else {
		// Here use the body for content.
		body, err := params.body()
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, urlStr, body)
		if err != nil {
			return nil, errors.Annotate(err, ErrCannotPrepareRequest, msgCannotPrepareRequest)
		}
		req.Header.Set("Content-Type", params.ContentType)
	}
	if params.Version != nil {
		req.Header.Set("Version", params.Version.String())
	}
	if params.Token != nil {
		req = token.RequestAdd(req, params.Token)
	}
	if params.Accept == "" {
		params.Accept = params.ContentType
	}
	if params.Accept != "" {
		req.Header.Set("Accept", params.Accept)
	}
	return req, nil
}

// analyzeResponse creates a response struct out of the HTTP response.
func analyzeResponse(resp *http.Response) (*Response, error) {
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Annotate(err, ErrAnalyzingResponse, msgAnalyzingResponse)
	}
	resp.Body.Close()
	return &Response{
		httpResp:    resp,
		contentType: resp.Header.Get("Content-Type"),
		content:     content,
	}, nil
}

// EOF
