// Tideland Go Library - Text - Generic JSON Processor
//
// Copyright (C) 2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package gjp

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.one/go/trace/errors"
)

//--------------------
// DIFFERENCE
//--------------------

// Diff manages the two parsed documents and their differences.
type Diff interface {
	// FirstDocument returns the first document passed to Diff().
	FirstDocument() Document

	// SecondDocument returns the second document passed to Diff().
	SecondDocument() Document

	// Differences returns a list of paths where the documents
	// have different content.
	Differences() []string

	// DifferenceAt returns the differences at the given path by
	// returning the first and the second value.
	DifferenceAt(path string) (Value, Value)
}

// diff implements Diff.
type diff struct {
	first  Document
	second Document
	paths  []string
}

// Compare parses and compares the documents and returns their differences.
func Compare(first, second []byte, separator string) (Diff, error) {
	fd, err := Parse(first, separator)
	if err != nil {
		return nil, err
	}
	sd, err := Parse(second, separator)
	if err != nil {
		return nil, err
	}
	d := &diff{
		first:  fd,
		second: sd,
	}
	err = d.compare()
	if err != nil {
		return nil, err
	}
	return d, nil
}

// CompareDocuments compares the documents and returns their differences.
func CompareDocuments(first, second Document, separator string) (Diff, error) {
	fd, ok := first.(*document)
	if !ok {
		return nil, errors.New(ErrInvalidDocument, "invalid first document implementation")
	}
	fd.separator = separator
	sd, ok := second.(*document)
	if !ok {
		return nil, errors.New(ErrInvalidDocument, "invalid second document implementation")
	}
	sd.separator = separator
	d := &diff{
		first:  fd,
		second: sd,
	}
	err := d.compare()
	if err != nil {
		return nil, err
	}
	return d, nil
}

// FirstDocument implements Diff.
func (d *diff) FirstDocument() Document {
	return d.first
}

// SecondDocument implements Diff.
func (d *diff) SecondDocument() Document {
	return d.second
}

// Differences implements Diff.
func (d *diff) Differences() []string {
	return d.paths
}

// DifferenceAt implements Diff.
func (d *diff) DifferenceAt(path string) (Value, Value) {
	firstValue := d.first.ValueAt(path)
	secondValue := d.second.ValueAt(path)
	return firstValue, secondValue
}

// compare iterates over the both documents looking for different
// values or even paths.
func (d *diff) compare() error {
	firstPaths := map[string]struct{}{}
	firstProcessor := func(path string, value Value) error {
		firstPaths[path] = struct{}{}
		if !value.Equals(d.second.ValueAt(path)) {
			d.paths = append(d.paths, path)
		}
		return nil
	}
	err := d.first.Process(firstProcessor)
	if err != nil {
		return err
	}
	secondProcessor := func(path string, value Value) error {
		_, ok := firstPaths[path]
		if ok {
			// Been there, done that.
			return nil
		}
		d.paths = append(d.paths, path)
		return nil
	}
	return d.second.Process(secondProcessor)
}

// EOF
