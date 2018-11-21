// Tideland Go Library - Trace - Logger
//
// Copyright (C) 2012-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package logger

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"os"
	"sync"

	"tideland.one/go/trace/location"
)

//--------------------
// LEVEL
//--------------------

// LogLevel describes the chosen log level between
// debug and critical.
type LogLevel int

// Log levels to control the logging output.
const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelCritical
	LevelFatal
)

//--------------------
// EXIT
//--------------------

// FatalExiterFunc defines a functions that will be called
// in case of a Fatalf call.
type FatalExiterFunc func()

// OSFatalExiter exits the application with os.Exit and
// the return code -1.
func OSFatalExiter() {
	os.Exit(-1)
}

// PanicFatalExiter exits the application with a panic.
func PanicFatalExiter() {
	panic("program aborted after fatal situation, see log")
}

//--------------------
// FILTER
//--------------------

// FilterFunc allows to filter the output of the logging. Filters
// have to return true if the received entry shall be filtered and
// not output.
type FilterFunc func(level LogLevel, info, msg string) bool

//--------------------
// LOGGER
//--------------------

// Logger provides a flexible configurable logging system.
type Logger struct {
	mu          sync.RWMutex
	level       LogLevel
	out         Writer
	fatalExiter FatalExiterFunc
	shallWrite  FilterFunc
}

// NewStandard returns a standard logger.
// NewStandardWriter(os.Stdout),
func NewStandard(out Writer) *Logger {
	return &Logger{
		level:       LevelInfo,
		out:         out,
		fatalExiter: OSFatalExiter,
	}
}

// NewTest returns a testing logger.
func NewTest() (*Logger, Entries) {
	out := NewTestWriter()
	l := &Logger{
		level:       LevelInfo,
		out:         out,
		fatalExiter: OSFatalExiter,
	}
	return l, out
}

// Level returns the current log level.
func (l *Logger) Level() LogLevel {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// SetLevel switches to a new log level and returns
// the current one.
func (l *Logger) SetLevel(level LogLevel) LogLevel {
	l.mu.Lock()
	defer l.mu.Unlock()
	current := l.level
	switch {
	case level <= LevelDebug:
		l.level = LevelDebug
	case level >= LevelFatal:
		l.level = LevelFatal
	default:
		l.level = level
	}
	return current
}

// SetFatalExiter sets the fatal exiter function and
// returns the current one.
func (l *Logger) SetFatalExiter(fef FatalExiterFunc) FatalExiterFunc {
	l.mu.Lock()
	defer l.mu.Unlock()
	current := l.fatalExiter
	if fef != nil {
		l.fatalExiter = fef
	}
	return current
}

// SetFilter sets the global output filter and returns the current one.
func (l *Logger) SetFilter(ff FilterFunc) FilterFunc {
	l.mu.Lock()
	defer l.mu.Unlock()
	current := l.shallWrite
	l.shallWrite = ff
	return current
}

// UnsetFilter removes the global output filter and returns the current one.
func (l *Logger) UnsetFilter() FilterFunc {
	l.mu.Lock()
	defer l.mu.Unlock()
	current := l.shallWrite
	l.shallWrite = nil
	return current
}

// Debugf logs a message at debug level.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.level <= LevelDebug {
		here := location.HereID(1)

		l.writeChecked(LevelDebug, here, fmt.Sprintf(format, args...))
	}
}

// Infof logs a message at info level.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.level <= LevelInfo {
		here := location.HereID(1)

		l.writeChecked(LevelInfo, here, fmt.Sprintf(format, args...))
	}
}

// Warningf logs a message at warning level.
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.level <= LevelWarning {
		here := location.HereID(1)

		l.writeChecked(LevelWarning, here, fmt.Sprintf(format, args...))
	}
}

// Errorf logs a message at error level.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.level <= LevelError {
		here := location.HereID(1)

		l.writeChecked(LevelError, here, fmt.Sprintf(format, args...))
	}
}

// Criticalf logs a message at critical level.
func (l *Logger) Criticalf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.level <= LevelCritical {
		here := location.HereID(1)

		l.writeChecked(LevelCritical, here, fmt.Sprintf(format, args...))
	}
}

// Fatalf logs a message independent of any level. After logging the message the functions
// calls the fatal exiter function, which by default means exiting the application
// with error code -1.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	here := location.HereID(1)

	l.out.Write(LevelFatal, here, fmt.Sprintf(format, args...))
	l.fatalExiter()
}

// writeChecked is used to check if a specific logging is wanted
// and writes the log if it is.
func (l *Logger) writeChecked(level LogLevel, here, msg string) {
	if l.shallWrite != nil {
		if !l.shallWrite(level, here, msg) {
			return
		}
	}
	l.out.Write(level, here, msg)
}

// EOF
