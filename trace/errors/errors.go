// Tideland Go Library - Trace - Errors
//
// Copyright (C) 2013-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package errors

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"strings"

	"tideland.dev/go/trace/location"
)

//--------------------
// CONSTANTS
//--------------------

// Error codes of the errors package.
const (
	ErrInvalidType       = "EINVTYPE"
	ErrNotYetImplemented = "ENOTYETI"
	ErrDeprecated        = "EFEATDEPR"

	msgInvalidType       = "passed error has type %T: '%s'"
	msgNotYetImplemented = "feature is not yet implemented: '%s'"
	msgDeprecated        = "feature is deprecated: '%s'"
)

//--------------------
// ERROR
//--------------------

// errorBox encapsulates an error.
type errorBox struct {
	err    error
	code   string
	msg    string
	hereID string
}

// newErrorBox creates an initialized error box.
func newErrorBox(err error, code string, msg string, args ...interface{}) *errorBox {
	return &errorBox{
		err:    err,
		code:   harmonize(code),
		msg:    fmt.Sprintf(msg, args...),
		hereID: location.HereID(2),
	}
}

// Error implements the error interface.
func (eb *errorBox) Error() string {
	if eb.err != nil {
		return fmt.Sprintf("%s [%s] %s: %v", eb.hereID, eb.code, eb.msg, eb.err)
	}
	return fmt.Sprintf("%s [%s] %s", eb.hereID, eb.code, eb.msg)
}

// errorCollection bundles multiple errors.
type errorCollection struct {
	errs []error
}

// Error implements the error interface.
func (ec *errorCollection) Error() string {
	errMsgs := make([]string, len(ec.errs))
	for i, err := range ec.errs {
		errMsgs[i] = err.Error()
	}
	return strings.Join(errMsgs, "\n")
}

// Annotate creates an error wrapping another one together with a
// a code.
func Annotate(err error, code string, msg string, args ...interface{}) error {
	return newErrorBox(err, code, msg, args...)
}

// New creates an error with the given code.
func New(code string, msg string, args ...interface{}) error {
	return newErrorBox(nil, code, msg, args...)
}

// Collect collects multiple errors into one.
func Collect(errs ...error) error {
	return &errorCollection{
		errs: errs,
	}
}

// Valid returns true if it is a valid error generated by
// this package.
func Valid(err error) bool {
	_, ok := err.(*errorBox)
	return ok
}

// IsError checks if an error is one created by this
// package and has the passed code
func IsError(err error, code string) bool {
	if e, ok := err.(*errorBox); ok {
		return e.code == harmonize(code)
	}
	return false
}

// Annotated returns the possibly annotated error. In case of
// a different error an invalid type error is returned.
func Annotated(err error) error {
	if e, ok := err.(*errorBox); ok {
		return e.err
	}
	return New(ErrInvalidType, msgInvalidType, err, err)
}

// Location returns the package and the file name as well as the line
// number of the error.
func Location(err error) (string, error) {
	if e, ok := err.(*errorBox); ok {
		return e.hereID, nil
	}
	return "", New(ErrInvalidType, msgInvalidType, err, err)
}

// Stack returns a slice of errors down to the lowest
// not annotated error.
func Stack(err error) []error {
	if eb, ok := err.(*errorBox); ok {
		return append([]error{eb}, Stack(eb.err)...)
	}
	return []error{err}
}

// All returns a slice of errors in case of collected errors.
func All(err error) []error {
	if ec, ok := err.(*errorCollection); ok {
		all := make([]error, len(ec.errs))
		copy(all, ec.errs)
		return all
	}
	return []error{err}
}

// DoAll iterates the passed function over all stacked
// or collected errors or simply the one that's passed.
func DoAll(err error, f func(error)) {
	switch terr := err.(type) {
	case *errorBox:
		for _, serr := range Stack(err) {
			f(serr)
		}
	case *errorCollection:
		for _, aerr := range All(err) {
			f(aerr)
		}
	default:
		f(terr)
	}
}

// IsInvalidTypeError checks if an error signals an invalid
// type in case of testing for an annotated error.
func IsInvalidTypeError(err error) bool {
	return IsError(err, ErrInvalidType)
}

// NotYetImplementedError returns the common error for a not yet
// implemented feature.
func NotYetImplementedError(feature string) error {
	return New(ErrNotYetImplemented, msgNotYetImplemented, feature)
}

// IsNotYetImplementedError checks if an error signals a not yet
// implemented feature.
func IsNotYetImplementedError(err error) bool {
	return IsError(err, ErrNotYetImplemented)
}

// DeprecatedError returns the common error for a deprecated
// feature.
func DeprecatedError(feature string) error {
	return New(ErrDeprecated, msgDeprecated, feature)
}

// IsDeprecatedError checks if an error signals deprecated
// feature.
func IsDeprecatedError(err error) bool {
	return IsError(err, ErrDeprecated)
}

//--------------------
// HELPER
//--------------------

// harmonzie ensures upper-case and a length of 8 by padding.
func harmonize(code string) string {
	var b strings.Builder
	code = strings.ToUpper(code) + "__________"
	for _, r := range code {
		if r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '_' {
			b.WriteRune(r)
		}
		if b.Len() == 10 {
			break
		}
	}
	return b.String()
}

// EOF
