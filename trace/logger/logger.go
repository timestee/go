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
	"io"
	"log"
	"os"
	"sync"
	"time"

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

// Logger defines the common interface for the different
// loggers.
type Logger interface {
	// Level returns the current log level.
	Level() LogLevel

	// SetLevel switches to a new log level and returns
	// the current one.
	SetLevel(level LogLevel) LogLevel

	// SetWriter sets a new Writer.
	SetWriter(out Writer) Writer

	// SetFatalExiter sets the fatal exiter function and
	// returns the current one.
	SetFatalExiter(fef FatalExiterFunc) FatalExiterFunc

	// SetFilter sets the global output filter and returns
	// the current one.
	SetFilter(ff FilterFunc) FilterFunc

	// UnsetFilter removes the global output filter and
	// returns the current one.
	UnsetFilter() FilterFunc

	// Debugf logs a message at debug level.
	Debugf(format string, args ...interface{})

	// Infof logs a message at info level.
	Infof(format string, args ...interface{})

	// Warningf logs a message at warning level.
	Warningf(format string, args ...interface{})

	// Errorf logs a message at error level.
	Errorf(format string, args ...interface{})

	// Criticalf logs a message at critical level.
	Criticalf(format string, args ...interface{})

	// Fatalf logs a message independent of any level. After
	// logging the message the functions calls the fatal exiter
	// function, which by default means exiting the application
	// with error code -1.
	Fatalf(format string, args ...interface{})
}

// logger implements Logger.
type logger struct {
	mu          sync.RWMutex
	level       LogLevel
	out         Writer
	fatalExiter FatalExiterFunc
	shallWrite  FilterFunc
}

// NewStandard returns a standard logger.
func NewStandard() Logger {
	return &logger{
		level:       LevelInfo,
		out:         NewStandardWriter(os.Stdout),
		fatalExiter: OSFatalExiter,
	}
}

// NewTest returns a testing logger.
func NewTest() (Logger, Entries) {
	w := NewTestWriter()
	l := &logger{
		level:       LevelInfo,
		out:         w,
		fatalExiter: OSFatalExiter,
	}
	return l, w
}

// Level implements Logger.
func (l *logger) Level() LogLevel {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// SetLevel implements Logger.
func (l *logger) SetLevel(level LogLevel) LogLevel {
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

// SetWriter implements Logger.
func (l *logger) SetWriter(out Writer) Writer {
	l.mu.Lock()
	defer l.mu.Unlock()
	current := l.out
	l.out = out
	return current
}

// SetFatalExiter implements Logger.
func (l *logger) SetFatalExiter(fef FatalExiterFunc) FatalExiterFunc {
	l.mu.Lock()
	defer l.mu.Unlock()
	current := l.fatalExiter
	if fef != nil {
		l.fatalExiter = fef
	}
	return current
}

// SetFilter implements Logger.
func (l *logger) SetFilter(ff FilterFunc) FilterFunc {
	l.mu.Lock()
	defer l.mu.Unlock()
	current := l.shallWrite
	l.shallWrite = ff
	return current
}

// UnsetFilter implements Logger.
func (l *logger) UnsetFilter() FilterFunc {
	l.mu.Lock()
	defer l.mu.Unlock()
	current := l.shallWrite
	l.shallWrite = nil
	return current
}

// Debugf logs a message at debug level.
func (l *logger) Debugf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.level <= LevelDebug {
		here := location.HereID(1)

		l.writeChecked(LevelDebug, here, fmt.Sprintf(format, args...))
	}
}

// Infof implements Logger.
func (l *logger) Infof(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.level <= LevelInfo {
		here := location.HereID(1)

		l.writeChecked(LevelInfo, here, fmt.Sprintf(format, args...))
	}
}

// Warningf implements Logger.
func (l *logger) Warningf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.level <= LevelWarning {
		here := location.HereID(1)

		l.writeChecked(LevelWarning, here, fmt.Sprintf(format, args...))
	}
}

// Errorf implements Logger.
func (l *logger) Errorf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.level <= LevelError {
		here := location.HereID(1)

		l.writeChecked(LevelError, here, fmt.Sprintf(format, args...))
	}
}

// Criticalf implements Logger.
func (l *logger) Criticalf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if l.level <= LevelCritical {
		here := location.HereID(1)

		l.writeChecked(LevelCritical, here, fmt.Sprintf(format, args...))
	}
}

// Fatalf implements logger.
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	here := location.HereID(1)

	l.out.Write(LevelFatal, here, fmt.Sprintf(format, args...))
	l.fatalExiter()
}

// writeChecked is used to check if a specific logging is wanted
// and writes the log if it is.
func (l *logger) writeChecked(level LogLevel, here, msg string) {
	if l.shallWrite != nil {
		if !l.shallWrite(level, here, msg) {
			return
		}
	}
	l.out.Write(level, here, msg)
}

//--------------------
// WRITER
//--------------------

// defaultTimeFormat controls how the timestamp of the standard
// logger is printed by default.
const defaultTimeFormat = "2006-01-02 15:04:05 Z07:00"

// levelText maps log levels to the according display texts.
var levelText = map[LogLevel]string{
	LevelDebug:    "DEBUG",
	LevelInfo:     "INFO",
	LevelWarning:  "WARNING",
	LevelError:    "ERROR",
	LevelCritical: "CRITICAL",
	LevelFatal:    "FATAL",
}

// Writer is the interface for different log writers.
type Writer interface {
	// Write writes the given message with additional
	// information at the specific log level.
	Write(level LogLevel, here, msg string)
}

// standardWriter is a simple writer writing to the given I/O
// writer. Beside the output it doesn't handle the levels differently.
type standardWriter struct {
	mu         sync.Mutex
	out        io.Writer
	timeFormat string
}

// NewTimeformatWriter creates a writer writing to the passed
// output and with the specified time format.
func NewTimeformatWriter(out io.Writer, timeFormat string) Writer {
	return &standardWriter{
		out:        out,
		timeFormat: timeFormat,
	}
}

// NewStandardWriter creates the standard writer writing
// to the passed output.
func NewStandardWriter(out io.Writer) Writer {
	return NewTimeformatWriter(out, defaultTimeFormat)
}

// Write implements Writer.
func (w *standardWriter) Write(level LogLevel, here, msg string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	text, ok := levelText[level]
	if !ok {
		text = "INVALID"
	}
	io.WriteString(w.out, time.Now().Format(w.timeFormat))
	io.WriteString(w.out, " [")
	io.WriteString(w.out, text)
	io.WriteString(w.out, "] ")
	io.WriteString(w.out, here)
	io.WriteString(w.out, " ")
	io.WriteString(w.out, msg)
	io.WriteString(w.out, "\n")
}

// goWriter just uses the standard go log package.
type goWriter struct{}

// NewGoWriter creates a writer using the Go log package.
func NewGoWriter() Writer {
	return &goWriter{}
}

// Write implements Writer.
func (w *goWriter) Write(level LogLevel, here, msg string) {
	text, ok := levelText[level]
	if !ok {
		text = "INVALID"
	}
	log.Println("["+text+"]", here, msg)
}

// Entries contains the collected entries of a test writer.
type Entries interface {
	// Len returns the number of collected entries.
	Len() int

	// Entries returns the collected entries.
	Entries() []string

	// Reset clears the collected entries.
	Reset()
}

// TestWriter extends the Writer interface with methods to
// retrieve and reset the collected data for testing purposes.
type TestWriter interface {
	Writer
	Entries
}

// testWriter simply collects logs to be evaluated inside of tests.
type testWriter struct {
	mu      sync.Mutex
	entries []string
}

// NewTestWriter returns a special writer for testing purposes.
func NewTestWriter() TestWriter {
	return &testWriter{}
}

// Write implements Writer.
func (w *testWriter) Write(level LogLevel, here, msg string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	text, ok := levelText[level]
	if !ok {
		text = "INVALID"
	}
	entry := fmt.Sprintf("%d [%s] %s %s", time.Now().UnixNano(), text, here, msg)
	w.entries = append(w.entries, entry)
}

// Len implements TestWriter.
func (w *testWriter) Len() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.entries)
}

// Entries implements TestWriter.
func (w *testWriter) Entries() []string {
	w.mu.Lock()
	defer w.mu.Unlock()
	entries := make([]string, len(w.entries))
	copy(entries, w.entries)
	return entries
}

// Reset implements TestWriter.
func (w *testWriter) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries = nil
}

// EOF
