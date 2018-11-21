// Tideland Go Library - Trace - Logger
//
// Copyright (C) 2012-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package logger of the Tideland Go Library provides a flexible way
// to log information with different levels and on different backends.
// A logger is created with
//
//     log := logger.NewStandard(logger.NewStandardOutWriter())
//
// or for tests with possible access to logged entries with
//
//     log, entries := logger.NewTest()
//
// The levels are Debug, Info, Warning, Error, Critical, and Fatal.
// Here log.Debugf() also logs information about file name, function
// name, and line number while log.Fatalf() may end the program
// depending on the set FatalExiterFunc.
//
// Different backends may be set. The StandardWriter writes to an
// io.Writer (initially os.Stdout), the GoWriter uses the Go log
// package, and the SysWriter uses the Go syslog package on the
// according operating systems. For testing the TestWriter exists.
// When created also access to the entries is returned. These can be
// used inside tests.
//
// Changes to the standard behavior can be made with log.SetLevel()
// and log.SetFatalExiter(). Own logger backends and exiter can be
// defined. Additionally a filter function allows to drill down the
// logged entries.
package logger

// EOF
