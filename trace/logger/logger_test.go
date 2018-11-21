// Tideland Go Library - Trace - Logger - Unit Tests
//
// Copyright (C) 2012-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package logger_test

//--------------------
// IMPORTS
//--------------------

import (
	"testing"

	"tideland.one/go/audit/asserts"
	"tideland.one/go/trace/logger"
)

//--------------------
// TESTS
//--------------------

// TestGetSetLevel tests the setting of the logging level.
func TestGetSetLevel(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	log, entries := logger.NewTest()

	log.SetLevel(logger.LevelDebug)
	log.Debugf("Debug.")
	log.Infof("Info.")
	log.Warningf("Warning.")
	log.Errorf("Error.")
	log.Criticalf("Critical.")

	assert.Length(entries, 5)
	entries.Reset()

	log.SetLevel(logger.LevelError)
	log.Debugf("Debug.")
	log.Infof("Info.")
	log.Warningf("Warning.")
	log.Errorf("Error.")
	log.Criticalf("Critical.")

	assert.Length(entries, 2)
	assert.Contents("TestGetSetLevel:44", entries.Entries()[0])
	assert.Contents("TestGetSetLevel:45", entries.Entries()[1])
	entries.Reset()
}

// TestFiltering tests the filtering of the logging.
func TestFiltering(t *testing.T) {
	assert := asserts.NewTesting(t, true)

	log, entries := logger.NewTest()

	log.SetLevel(logger.LevelDebug)
	log.SetFilter(func(level logger.LogLevel, here, msg string) bool {
		return level >= logger.LevelWarning && level <= logger.LevelError
	})

	log.Debugf("Debug.")
	log.Infof("Info.")
	log.Warningf("Warning.")
	log.Errorf("Error.")
	log.Criticalf("Critical.")

	assert.Length(entries, 2)
	entries.Reset()

	log.UnsetFilter()

	log.Debugf("Debug.")
	log.Infof("Info.")
	log.Warningf("Warning.")
	log.Errorf("Error.")
	log.Criticalf("Critical.")

	assert.Length(entries, 5)
	entries.Reset()
}

// TestGoLogger tests logging with the go logger.
func TestGoLogger(t *testing.T) {
	log := logger.NewStandard(logger.NewGoWriter())

	log.Debugf("Debug.")
	log.Infof("Info.")
	log.Warningf("Warning.")
	log.Errorf("Error.")
	log.Criticalf("Critical.")
}

// TestSysLogger tests logging with the syslogger.
func TestSysLogger(t *testing.T) {
	assert := asserts.NewTesting(t, true)
	out, err := logger.NewSysWriter("GOTRACE")
	assert.Nil(err)
	log := logger.NewStandard(out)

	log.SetLevel(logger.LevelDebug)

	log.Debugf("Debug.")
	log.Infof("Info.")
	log.Warningf("Warning.")
	log.Errorf("Error.")
	log.Criticalf("Critical.")
}

// TestFatalExit tests the call of the fatal exiter after a
// fatal error log.
func TestFatalExit(t *testing.T) {
	assert := asserts.NewTesting(t, true)

	log, entries := logger.NewTest()

	exited := false
	fatalExiter := func() {
		exited = true
	}

	log.SetFatalExiter(fatalExiter)

	log.Fatalf("fatal")

	assert.Length(entries, 1)
	assert.True(exited)
	entries.Reset()
}

// EOF
