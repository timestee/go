// Tideland Go Library - Network - REST - Core
//
// Copyright (C) 2009-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package core

//--------------------
// IMPORTS
//--------------------

import (
	"io/ioutil"
	"net/http"
	"sync"
	"text/template"
	"time"

	"tideland.one/go/trace/errors"
)

//--------------------
// TEMPLATES CACHE ITEM
//--------------------

// templatesCacheItem stores the parsed template and the
// content type.
type templatesCacheItem struct {
	id             string
	timestamp      time.Time
	parsedTemplate *template.Template
	contentType    string
}

// isValid checks if the the entry is younger than the
// passed validity period.
func (ti *templatesCacheItem) isValid(validityPeriod time.Duration) bool {
	return ti.timestamp.Add(validityPeriod).After(time.Now())
}

// render the cached entry.
func (ti *templatesCacheItem) render(rw http.ResponseWriter, data interface{}) error {
	rw.Header().Set("Content-Type", ti.contentType)
	if err := ti.parsedTemplate.Execute(rw, data); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}

//--------------------
// TEMPLATES CACHE
//--------------------

// TemplatesCache caches and renders templates.
type TemplatesCache struct {
	mutex sync.RWMutex
	items map[string]*templatesCacheItem
}

// newTemplatesCache creates a new template cache.
func newTemplatesCache() *TemplatesCache {
	return &TemplatesCache{
		items: make(map[string]*templatesCacheItem),
	}
}

// Parse parses a raw template an stores it.
func (t *TemplatesCache) Parse(id, rawTemplate, contentType string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	parsedTemplate, err := template.New(id).Parse(rawTemplate)
	if err != nil {
		return err
	}
	t.items[id] = &templatesCacheItem{id, time.Now(), parsedTemplate, contentType}
	return nil
}

// LoadAndParse loads a template from filesystem, parses it, and stores it.
func (t *TemplatesCache) LoadAndParse(id, filename, contentType string) error {
	rawTemplate, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return t.Parse(id, string(rawTemplate), contentType)
}

// Render executes the pre-parsed template with the data. It also sets the content type header.
func (t *TemplatesCache) Render(rw http.ResponseWriter, id string, data interface{}) error {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	entry, ok := t.items[id]
	if !ok {
		return errors.New(ErrNoCachedTemplate, msgNoCachedTemplate, id)
	}
	return entry.render(rw, data)
}

// LoadAndRender checks if the template with the given id has already been parsed.
// In this case it will use it, otherwise the template will be loaded, parsed, added
// to the cache, and used then.
func (t *TemplatesCache) LoadAndRender(rw http.ResponseWriter, id, filename, contentType string, data interface{}) error {
	t.mutex.RLock()
	_, ok := t.items[id]
	t.mutex.RUnlock()
	if !ok {
		if err := t.LoadAndParse(id, filename, contentType); err != nil {
			return err
		}
	}
	return t.Render(rw, id, data)
}

//--------------------
// RENDERER
//--------------------

// Renderer renders templates. It is returned by a Job and knows
// where to render it.
type Renderer struct {
	rw http.ResponseWriter
	tc *TemplatesCache
}

// Render executes the pre-parsed template with the data. It also sets the content type header.
func (r *Renderer) Render(id string, data interface{}) error {
	return r.tc.Render(r.rw, id, data)
}

// LoadAndRender checks if the template with the given ID has already been parsed. In this case
// it will use it, otherwise the template will be loaded, parsed, added to the cache, and used then.
func (r *Renderer) LoadAndRender(id, filename, contentType string, data interface{}) error {
	return r.tc.LoadAndRender(r.rw, id, filename, contentType, data)
}

// EOF
