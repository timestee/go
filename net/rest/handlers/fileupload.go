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
	"mime/multipart"

	"tideland.one/go/net/rest/core"
	"tideland.one/go/trace/errors"
)

//--------------------
// FILE UPLOAD HANDLER
//--------------------

const defaultMaxMemory = 32 << 20 // 32 MB

// FileUploadProcessor defines the function used for the processing
// of the uploaded file. It has to be specified by the user of the
// handler and e.g. persists the received data in the file system or
// a database.
type FileUploadProcessor func(job *core.Job, header *multipart.FileHeader, file multipart.File) error

// fileUploadHandler handles uploading POST requests.
type fileUploadHandler struct {
	id        string
	processor FileUploadProcessor
}

// NewFileUploadHandler creates a new handler for the uploading of files.
func NewFileUploadHandler(id string, processor FileUploadProcessor) core.ResourceHandler {
	return &fileUploadHandler{
		id:        id,
		processor: processor,
	}
}

// Init is specified on the ResourceHandler interface.
func (h *fileUploadHandler) ID() string {
	return h.id
}

// ID is specified on the ResourceHandler interface.
func (h *fileUploadHandler) Init(env *core.Environment, domain, resource string) error {
	return nil
}

// Post is specified on the PostResourceHandler interface.
func (h *fileUploadHandler) Post(job *core.Job) (bool, error) {
	if err := job.Request().ParseMultipartForm(defaultMaxMemory); err != nil {
		return false, errors.Annotate(err, ErrUploadingFile, msgUploadingFile)
	}
	for _, headers := range job.Request().MultipartForm.File {
		for _, header := range headers {
			job.Environment().Log().Infof("receiving file %q", header.Filename)
			// Open file and process it.
			if infile, err := header.Open(); err != nil {
				return false, errors.Annotate(err, ErrUploadingFile, msgUploadingFile)
			} else if err := h.processor(job, header, infile); err != nil {
				return false, errors.Annotate(err, ErrUploadingFile, msgUploadingFile)
			}
		}
	}
	return true, nil
}

// EOF
