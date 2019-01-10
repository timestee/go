// Tideland Go Library - Trace - Logger - No SysLogger
//
// Copyright (C) 2012-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// +build windows plan9 nacl

package logger

//--------------------
// IMPORTS
//--------------------

import (
	"log"
)

//--------------------
// SYSWRITER
//--------------------

// sysWriter uses the Go syslog package. It does not work
// on Windows or Plan9.
type sysWriter struct {
	tag string
}

// NewSysWriter creates a writer using the Go syslog package.
// It does not work on Windows or Plan9. Here the Go log
// package is used.
func NewSysWriter(tag string) (Writer, error) {
	if len(tag) > 0 {
		tag = "(" + tag + ")"
	}
	return &sysWriter{tag}, nil
}

// Write implements Writer.
func (w *sysWriter) Write(level LogLevel, here, msg string) {
	text, ok := levelText[level]
	if !ok {
		text = "INVALID"
	}
	log.Println("["+text+"]", info, msg)
}

// EOF
