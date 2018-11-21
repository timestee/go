// Tideland Go Library - Network - REST - Core
//
// Copyright (C) 2009-2017 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package core

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"
	"strings"

	"tideland.one/go/text/stringex"
)

//--------------------
// CONSTANTS
//--------------------

// Path indexes for the different parts.
const (
	PathDomain     = 0
	PathResource   = 1
	PathResourceID = 2
)

//--------------------
// PATH
//--------------------

// Path provides access to the parts of a request path interesting for handling
// a job.
type Path struct {
	parts []string
}

// newPath returns the analyzed path.
func newPath(env *Environment, r *http.Request) *Path {
	parts := stringex.SplitMap(r.URL.Path, "/", func(part string) (string, bool) {
		if part == "" {
			return "", false
		}
		return part, true
	})[env.basepartsLen:]
	switch len(parts) {
	case 1:
		parts = append(parts, env.defaultResource)
	case 0:
		parts = append(parts, env.defaultDomain, env.defaultResource)
	}
	return &Path{
		parts: parts,
	}
}

// Length returns the number of parts of the path.
func (p *Path) Length() int {
	return len(p.parts)
}

// ContainsSubResourceIDs returns true, if the path doesn't end after the resource ID,
// e.g. to address items of an order.
//
// Example: /shop/orders/12345/item/1
func (p *Path) ContainsSubResourceIDs() bool {
	return len(p.parts) > 3
}

// Part returns the parts of the URL path based on the index or an empty string.
func (p *Path) Part(index int) string {
	if len(p.parts) <= index {
		return ""
	}
	return p.parts[index]
}

// Domain returns the requests domain.
func (p *Path) Domain() string {
	return p.parts[PathDomain]
}

// Resource returns the requests resource.
func (p *Path) Resource() string {
	return p.parts[PathResource]
}

// ResourceID returns the requests resource ID.
func (p *Path) ResourceID() string {
	if len(p.parts) > 2 {
		return p.parts[PathResourceID]
	}
	return ""
}

// JoinedResourceID returns the requests resource ID together with all following
// parts of the path.
func (p *Path) JoinedResourceID() string {
	if len(p.parts) > 2 {
		return strings.Join(p.parts[PathResourceID:], "/")
	}
	return ""
}

// EOF
