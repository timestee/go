// Tideland Go Library - Trace - Errors - Unit Tests
//
// Copyright (C) 2013-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package errors_test

//--------------------
// IMPORTS
//--------------------

import (
	goerrors "errors"
	"testing"

	"tideland.one/go/audit/asserts"
	"tideland.one/go/trace/errors"
)

//--------------------
// TESTS
//--------------------

// TestIsError tests the creation and checking of errors.
func TestIsError(t *testing.T) {
	assert := asserts.NewTesting(t, true)

	ec := "E001"
	emsg := "test error %d"
	err := errors.New(ec, emsg, 1)

	assert.Equal(err.Error(), "(tideland.one/go/trace/errors_test) errors_test.go:TestIsError:32 [E001] test error 1")
	assert.True(errors.IsError(err, ec))
	assert.False(errors.IsError(err, "0"))

	err = testError("test error 2")

	assert.ErrorMatch(err, "test error 2")
	assert.False(errors.IsError(err, ec))
	assert.False(errors.IsError(err, "0"))

	err = goerrors.New("42")

	assert.False(errors.IsError(err, ec))
	assert.False(errors.IsError(err, "0"))
}

// TestValidation checks the validation of errors and
// the retrieval of details.
func TestValidation(t *testing.T) {
	assert := asserts.NewTesting(t, true)

	// First a valid error.
	ec := "E001"
	emsg := "valid"
	err := errors.New(ec, emsg)
	assert.True(errors.Valid(err))

	hereID, lerr := errors.Location(err)
	assert.Nil(lerr)
	assert.Equal(hereID, "(tideland.one/go/trace/errors_test) errors_test.go:TestValidation:58")

	// Now an invalid error.
	err = goerrors.New("ouch")
	assert.False(errors.Valid(err))

	hereID, lerr = errors.Location(err)
	assert.True(errors.IsInvalidTypeError(lerr))
	assert.Empty(hereID)
}

// TestAnnotation the annotation of errors with new errors.
func TestAnnotation(t *testing.T) {
	assert := asserts.NewTesting(t, true)

	err1 := testError("wrapped")
	err2 := errors.Annotate(err1, "E123", "1st annotated")
	err3 := errors.Annotate(err2, "E456", "2nd annotated")

	assert.ErrorMatch(err3, `.* \[E456\] 2nd annotated: .* \[E123\] 1st annotated: wrapped`)
	assert.Equal(errors.Annotated(err3), err2)
	assert.Equal(errors.Annotated(err2), err1)
	assert.Length(errors.Stack(err3), 3)

	assert.True(errors.IsInvalidTypeError(errors.Annotated(err1)))
}

// TestCollection tests the collection of multiple errors to one.
func TestCollection(t *testing.T) {
	assert := asserts.NewTesting(t, true)

	errA := testError("E001")
	errB := testError("E002")
	errC := testError("E003")
	errD := testError("E004")
	cerr := errors.Collect(errA, errB, errC, errD)

	assert.ErrorMatch(cerr, "E001\nE002\nE003\nE004")
}

// TestDoAll tests the iteration over errors.
func TestDoAll(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	msgs := []string{}
	f := func(err error) {
		msgs = append(msgs, err.Error())
	}

	// Test it on annotated errors.
	errX := testError("E000")
	errA := errors.Annotate(errX, "E001", "foo")
	errB := errors.Annotate(errA, "E002", "bar")
	errC := errors.Annotate(errB, "E003", "baz")
	errD := errors.Annotate(errC, "E004", "yadda")

	errors.DoAll(errD, f)

	assert.Length(msgs, 5)

	// Test it on collected errors.
	msgs = []string{}
	errA = testError("E001")
	errB = testError("E002")
	errC = testError("E003")
	errD = testError("E004")
	cerr := errors.Collect(errA, errB, errC, errD)

	errors.DoAll(cerr, f)

	assert.Equal(msgs, []string{"E001", "E002", "E003", "E004"})

	// Test it on a single error.
	msgs = []string{}
	errA = testError("foo")

	errors.DoAll(errA, f)

	assert.Equal(msgs, []string{"foo"})
}

//--------------------
// HELPERS
//--------------------

type testError string

func (e testError) Error() string {
	return string(e)
}

// EOF
