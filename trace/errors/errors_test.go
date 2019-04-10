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

	"tideland.dev/go/audit/asserts"
	"tideland.dev/go/trace/errors"
)

//--------------------
// TESTS
//--------------------

// TestIsError tests the creation and checking of errors.
func TestIsError(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	ec := "ETEST"
	emsg := "test error %d"
	err := errors.New(ec, emsg, 1)

	assert.Equal(err.Error(), "(tideland.dev/go/trace/errors_test:errors_test.go:TestIsError:32) [ETEST_____] test error 1")
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
	assert := asserts.NewTesting(t, asserts.FailStop)

	// First a valid error.
	ec := "EINVALID"
	emsg := "valid"
	err := errors.New(ec, emsg)
	assert.True(errors.Valid(err))

	hereID, lerr := errors.Location(err)
	assert.Nil(lerr)
	assert.Equal(hereID, "(tideland.dev/go/trace/errors_test:errors_test.go:TestValidation:58)")

	// Now an invalid error.
	err = goerrors.New("ouch")
	assert.False(errors.Valid(err))

	hereID, lerr = errors.Location(err)
	assert.True(errors.IsInvalidTypeError(lerr))
	assert.Empty(hereID)
}

// TestAnnotation the annotation of errors with new errors.
func TestAnnotation(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	err1 := testError("wrapped")
	err2 := errors.Annotate(err1, "EFIRST", "1st annotated")
	err3 := errors.Annotate(err2, "ESECOND", "2nd annotated")

	assert.ErrorMatch(err3, `.* \[ESECOND___\] 2nd annotated: .* \[EFIRST____\] 1st annotated: wrapped`)
	assert.Equal(errors.Annotated(err3), err2)
	assert.Equal(errors.Annotated(err2), err1)
	assert.Length(errors.Stack(err3), 3)

	assert.True(errors.IsInvalidTypeError(errors.Annotated(err1)))
}

// TestCollection tests the collection of multiple errors to one.
func TestCollection(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)

	errA := testError("EONE")
	errB := testError("ETWO")
	errC := testError("ETHREE")
	errD := testError("EFOUR")
	cerr := errors.Collect(errA, errB, errC, errD)

	assert.ErrorMatch(cerr, "EONE\nETWO\nETHREE\nEFOUR")
}

// TestDoAll tests the iteration over errors.
func TestDoAll(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	msgs := []string{}
	f := func(err error) {
		msgs = append(msgs, err.Error())
	}

	// Test it on annotated errors.
	errX := testError("EINIT")
	errA := errors.Annotate(errX, "EFOO", "foo")
	errB := errors.Annotate(errA, "EBAR", "bar")
	errC := errors.Annotate(errB, "EBAZ", "baz")
	errD := errors.Annotate(errC, "EYADDA", "yadda")

	errors.DoAll(errD, f)

	assert.Length(msgs, 5)

	// Test it on collected errors.
	msgs = []string{}
	errA = testError("EFOO")
	errB = testError("EBAR")
	errC = testError("EBAZ")
	errD = testError("EYADDA")
	cerr := errors.Collect(errA, errB, errC, errD)

	errors.DoAll(cerr, f)

	assert.Equal(msgs, []string{"EFOO", "EBAR", "EBAZ", "EYADDA"})

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
