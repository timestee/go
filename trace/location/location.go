// Tideland Go Library - Trace - Location
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package location // import "tideland.dev/go/trace/location"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

//--------------------
// LOCATION
//--------------------

// Here returns package, file, and line at the given offset.
func Here(offset int) (string, string, string, int) {
	l := here(offset)
	return l.pkg, l.file, l.fun, l.line
}

// HereID returns package, file, and line at the given offset as ID.
func HereID(offset int) string {
	l := here(offset)
	return l.id
}

// HereCode returns returns a code based on the given offset.
func HereCode(prefix string, offset int) string {
	l := here(offset)
	pparts := strings.Split(l.pkg, "/")
	for _, ppart := range pparts {
		prefix += ppart[0:1]
	}
	prefix += l.file[0:1]
	prefix += strconv.Itoa(l.line)
	return strings.ToUpper(prefix)
}

//--------------------
// BACKEND
//--------------------

// location contains the details and the formatted ID
// of one location.
type location struct {
	pkg  string
	file string
	fun  string
	line int
	id   string
}

// Cached locations.
var (
	mu        sync.Mutex
	locations = make(map[uintptr]*location)
)

// here returns the location at the given offset.
func here(offset int) *location {
	mu.Lock()
	defer mu.Unlock()
	// Fix the offset.
	offset += 3
	if offset < 3 {
		offset = 3
	}
	// Retrieve program counters.
	pcs := make([]uintptr, 1)
	n := runtime.Callers(offset, pcs)
	if n == 0 {
		return nil
	}
	pcs = pcs[:n]
	// Check cache.
	pc := pcs[0]
	l, ok := locations[pc]
	if ok {
		return l
	}
	// Build ID based on program counters.
	frames := runtime.CallersFrames(pcs)
	for {
		frame, more := frames.Next()
		pkg, fun := path.Split(frame.Function)
		parts := strings.Split(fun, ".")
		pkg = path.Join(pkg, parts[0])
		fun = strings.Join(parts[1:], ".")
		_, file := path.Split(frame.File)
		id := fmt.Sprintf("(%s:%s:%s:%d)", pkg, file, fun, frame.Line)
		if !more {
			l := &location{
				pkg:  pkg,
				file: file,
				fun:  fun,
				line: frame.Line,
				id:   id,
			}
			locations[pc] = l
			return l
		}
	}
}

// EOF
