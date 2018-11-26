// Tideland Go Library - Network - REST - Handlers
//
// Copyright (C) 2009-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package handlers

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"path/filepath"
	"strings"

	"tideland.one/go/rest/core"
	"tideland.one/go/trace/errors"
	"tideland.one/go/trace/logger"
)

//--------------------
// FILE SERVE HANDLER
//--------------------

// fileServeHandler implements the file server.
type fileServeHandler struct {
	id  string
	dir string
}

// NewFileServeHandler creates a new handler serving the files names
// by the resource ID part out of the passed directory.
func NewFileServeHandler(id, dir string) core.ResourceHandler {
	pdir := filepath.FromSlash(dir)
	if !strings.HasSuffix(pdir, string(filepath.Separator)) {
		pdir += string(filepath.Separator)
	}
	return &fileServeHandler{id, pdir}
}

// ID is specified on the ResourceHandler interface.
func (h *fileServeHandler) ID() string {
	return h.id
}

// Init is specified on the ResourceHandler interface.
func (h *fileServeHandler) Init(env *core.Environment, domain, resource string) error {
	return nil
}

// Get is specified on the GetResourceHandler interface.
func (h *fileServeHandler) Get(job *core.Job) (bool, error) {
	cleanResourceID := filepath.Clean(strings.Replace(job.Path().ResourceID(), "../", "", -1))
	filename, err := filepath.Abs(h.dir + cleanResourceID)
	if err != nil {
		logger.Errorf("file '%s' does not exist", filename)
		return false, errors.Annotate(err, ErrDownloadingFile, msgDownloadingFile, filename)
	}
	logger.Infof("serving file '%s'", filename)
	http.ServeFile(job.ResponseWriter(), job.Request(), filename)
	return true, nil
}

// EOF
