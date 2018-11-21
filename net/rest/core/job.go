// Tideland Go Library - Network - REST - Core
//
// Copyright (C) 2009-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package core

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"tideland.one/go/dsa/version"
	// "tideland.one/go/trace/logger"
)

//--------------------
// JOB
//--------------------

// Job encapsulates all the needed information for handling
// a request.
type Job struct {
	environment    *Environment
	ctx            context.Context
	request        *http.Request
	responseWriter http.ResponseWriter
	version        version.Version
	path           *Path
}

// newJob parses the URL and returns the prepared job.
func newJob(env *Environment, r *http.Request, rw http.ResponseWriter) *Job {
	// Init the job.
	j := &Job{
		environment:    env,
		request:        r,
		responseWriter: rw,
		path:           newPath(env, r),
	}
	// Retrieve the requested version of the API.
	vsnstr := j.request.Header.Get("Version")
	if vsnstr == "" {
		j.version = version.New(1, 0, 0)
	} else {
		vsn, err := version.Parse(vsnstr)
		if err != nil {
			// logger.Errorf("invalid request version: %v", err)
			j.version = version.New(1, 0, 0)
		} else {
			j.version = vsn
		}
	}
	return j
}

// Environment returns the server environment.
func (j *Job) Environment() *Environment {
	return j.environment
}

// Request returns the used Go HTTP request.
func (j *Job) Request() *http.Request {
	return j.request
}

// ResponseWriter returns the used Go HTTP response writer.
func (j *Job) ResponseWriter() http.ResponseWriter {
	return j.responseWriter
}

// Path returns access to the request path inside the URL.
func (j *Job) Path() *Path {
	return j.path
}

// Context returns a job context also containing the job itself.
func (j *Job) Context() context.Context {
	// Lazy init.
	if j.ctx == nil {
		j.ctx = newJobContext(j.environment.ctx, j)
	}
	return j.ctx
}

// EnhanceContext allows to enhance the job context values, a deadline, a timeout, or a cancel. So
// e.g. a first handler in a handler queue can store authentication information in the context
// and a following handler can use it (see the JWTAuthorizationHandler).
func (j *Job) EnhanceContext(f func(ctx context.Context) context.Context) {
	ctx := j.Context()
	j.ctx = f(ctx)
}

// Version returns the requested API version for this job. If none is set the version 1.0.0
// will be returned as default. It will be retrieved aut of the header Version.
func (j *Job) Version() version.Version {
	return j.version
}

// SetVersion allows to set an API version for the response. If none is set the version 1.0.0
// will be set as default. It will be set in the header Version.
func (j *Job) SetVersion(vsn version.Version) {
	j.version = vsn
}

// AcceptsContentType checks if the requestor accepts a given content type.
func (j *Job) AcceptsContentType(contentType string) bool {
	return strings.Contains(j.request.Header.Get("Accept"), contentType)
}

// HasContentType checks if the sent content has the given content type.
func (j *Job) HasContentType(contentType string) bool {
	return strings.Contains(j.request.Header.Get("Content-Type"), contentType)
}

// Languages returns the accepted language with the quality values.
func (j *Job) Languages() Languages {
	accept := j.request.Header.Get("Accept-Language")
	languages := Languages{}
	for _, part := range strings.Split(accept, ",") {
		lv := strings.Split(part, ";")
		if len(lv) == 1 {
			languages = append(languages, Language{lv[0], 1.0})
		} else {
			value, err := strconv.ParseFloat(lv[1], 64)
			if err != nil {
				value = 0.0
			}
			languages = append(languages, Language{lv[0], value})
		}
	}
	languages = sort.Reverse(languages).(Languages)
	return languages
}

// InternalPath builds an internal path out of the passed parts. It can be used to build
// requests to other systems using the same library.
func (j *Job) InternalPath(domain, resource, resourceID string, query ...KeyValue) string {
	path := j.createPath(domain, resource, resourceID)
	if len(query) > 0 {
		path += "?" + KeyValues(query).String()
	}
	return path
}

// Redirect to a domain, resource and resource ID (optional).
func (j *Job) Redirect(domain, resource, resourceID string) {
	path := j.createPath(domain, resource, resourceID)
	http.Redirect(j.responseWriter, j.request, path, http.StatusTemporaryRedirect)
}

// Renderer returns a template renderer.
func (j *Job) Renderer() *Renderer {
	return &Renderer{j.responseWriter, j.environment.templatesCache}
}

// GOB returns a GOB formatter.
func (j *Job) GOB() Formatter {
	return &gobFormatter{j}
}

// JSON returns a JSON formatter.
func (j *Job) JSON(html bool) Formatter {
	return &jsonFormatter{j, html}
}

// XML returns a XML formatter.
func (j *Job) XML() Formatter {
	return &xmlFormatter{j}
}

// Query returns a convenient access to query values.
func (j *Job) Query() Values {
	return &values{j.request.URL.Query()}
}

// Form returns a convenient access to form values.
func (j *Job) Form() Values {
	return &values{j.request.PostForm}
}

// String implements the fmt.Stringer interface.
func (j *Job) String() string {
	path := j.createPath(j.Path().Domain(), j.Path().Resource(), j.Path().ResourceID())
	return fmt.Sprintf("%s %s", j.request.Method, path)
}

// createPath creates a path out of the major URL parts.
func (j *Job) createPath(domain, resource, resourceID string) string {
	parts := append(j.environment.baseparts, domain, resource)
	if resourceID != "" {
		parts = append(parts, resourceID)
	}
	path := strings.Join(parts, "/")
	return "/" + path
}

// EOF
