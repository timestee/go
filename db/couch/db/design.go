// Tideland Go Library - DB - CouchDB Client - Core
//
// Copyright (C) 2016-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package db

//--------------------
// DESIGN
//--------------------

// Design provides convenient access to a design document.
type Design struct {
	db       *DB
	id       string
	document *designDocument
}

// newDesign creates a design instance.
func newDesign(db *DB, id string) (*Design, error) {
	designID := "_design/" + id
	ok, err := db.HasDocument(designID)
	if err != nil {
		return nil, err
	}
	document := designDocument{}
	if ok {
		// Read the design document.
		resp := db.ReadDocument(designID)
		if !resp.IsOK() {
			return nil, resp.Error()
		}
		err = resp.Document(&document)
		if err != nil {
			return nil, err
		}
	} else {
		// Create the design document.
		document = designDocument{
			ID:       designID,
			Language: "javascript",
		}
	}
	d := &Design{
		db:       db,
		id:       id,
		document: &document,
	}
	return d, nil
}

// ID returns the ID of the design.
func (d *Design) ID() string {
	return d.id
}

// Language returns the language for views and shows.
func (d *Design) Language() string {
	return d.document.Language
}

// SetLanguage sets the language for views and shows.
func (d *Design) SetLanguage(language string) {
	d.document.Language = language
}

// View returns the map and the reduce functions of the
// view with the ID, otherwise false.
func (d *Design) View(id string) (string, string, bool) {
	if d.document.Views == nil {
		d.document.Views = designViews{}
	}
	view, ok := d.document.Views[id]
	if !ok {
		return "", "", false
	}
	return view.Map, view.Reduce, true
}

// SetView sets the map and the reduce functions of the
// view with the ID.
func (d *Design) SetView(id, mapf, reducef string) {
	if d.document.Views == nil {
		d.document.Views = designViews{}
	}
	d.document.Views[id] = designView{
		Map:    mapf,
		Reduce: reducef,
	}
}

// Show returns the show function with the ID, otherwise false.
func (d *Design) Show(id string) (string, bool) {
	if d.document.Shows == nil {
		d.document.Shows = map[string]string{}
	}
	show, ok := d.document.Shows[id]
	if !ok {
		return "", false
	}
	return show, true
}

// SetShow sets the show function with the ID.
func (d *Design) SetShow(id, showf string) {
	if d.document.Shows == nil {
		d.document.Shows = map[string]string{}
	}
	d.document.Shows[id] = showf
}

// Write creates a new design document or updates an
// existing one.
func (d *Design) Write(params ...Parameter) *ResultSet {
	if d.document.Revision == "" {
		return d.db.CreateDocument(d.document, params...)
	}
	return d.db.UpdateDocument(d.document, params...)
}

// Delete removes a design document.
func (d *Design) Delete(params ...Parameter) *ResultSet {
	return d.db.DeleteDocument(d.document, params...)
}

//--------------------
// DESIGN DOCUMENT
//--------------------

// designView defines a view inside a design document.
type designView struct {
	Map    string `json:"map,omitempty"`
	Reduce string `json:"reduce,omitempty"`
}

type designViews map[string]designView

// designAttachment defines an attachment inside a design document.
type designAttachment struct {
	Stub        bool   `json:"stub,omitempty"`
	ContentType string `json:"content_type,omitempty"`
	Length      int    `json:"length,omitempty"`
}

type designAttachments map[string]designAttachment

// designDocument contains the data of view design documents.
type designDocument struct {
	ID                     string            `json:"_id"`
	Revision               string            `json:"_rev,omitempty"`
	Language               string            `json:"language,omitempty"`
	ValidateDocumentUpdate string            `json:"validate_doc_update,omitempty"`
	Views                  designViews       `json:"views,omitempty"`
	Shows                  map[string]string `json:"shows,omitempty"`
	Attachments            designAttachments `json:"_attachments,omitempty"`
	Signatures             map[string]string `json:"signatures,omitempty"`
	Libraries              interface{}       `json:"libs,omitempty"`
}

// EOF
