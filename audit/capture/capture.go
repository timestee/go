// Tideland Go Library - Audit - Capture
//
// Copyright (C) 2017-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package capture

//--------------------
// IMPORTS
//--------------------

import (
	"bytes"
	"io"
	"os"
)

//--------------------
// CAPTURED
//--------------------

// Captured provides access to the captured output in
// multiple ways.
type Captured interface {
	// Bytes returns the captured content as bytes.
	Bytes() []byte

	// String returns the captured content as string.
	String() string

	// Len returns the number of captured bytes.
	Len() int
}

// captured implements Captured.
type captured struct {
	buffer []byte
}

// Bytes implements Captured.
func (c *captured) Bytes() []byte {
	buf := make([]byte, c.Len())
	copy(buf, c.buffer)
	return buf
}

// String implements Stringer.
func (c *captured) String() string {
	return string(c.Bytes())
}

// Len implements Captured.
func (c *captured) Len() int {
	return len(c.buffer)
}

//--------------------
// CAPTURING
//--------------------

// Stdout captures Stdout.
func Stdout(f func()) Captured {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	outC := make(chan []byte)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.Bytes()
	}()

	w.Close()
	os.Stdout = old
	cptrd := captured{
		buffer: <-outC,
	}
	return &cptrd
}

// Stderr captures Stderr.
func Stderr(f func()) Captured {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	f()

	outC := make(chan []byte)

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.Bytes()
	}()

	w.Close()
	os.Stderr = old
	cptrd := captured{
		buffer: <-outC,
	}
	return &cptrd
}

// Both captures Stdout and Stderr.
func Both(f func()) (Captured, Captured) {
	var cerr Captured
	ff := func() {
		cerr = Stderr(f)
	}
	cout := Stdout(ff)
	return cout, cerr
}

// EOF
