// Tideland Go Library - Together - Notifier
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package notifier // import "tideland.dev/go/together/notifier"

//--------------------
// IMPORTS
//--------------------

import (
	"reflect"
	"sync"
)

//--------------------
// CLOSER
//--------------------

// Closer signals a typical for-select-loop to terminate. It listens
// to multiple structs channels itself, e.g. from a context or other
// termination signalling functions.
type Closer struct {
	once  sync.Once
	cases []reflect.SelectCase
	doneC chan struct{}
}

// NewCloser creates a new Closer instance.
func NewCloser(closeCs ...<-chan struct{}) *Closer {
	c := &Closer{
		doneC: make(chan struct{}),
	}
	for _, closeC := range closeCs {
		c.cases = append(c.cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(closeC),
		})
	}
	return c
}

// Go starts the backend of the Closer.
func (c *Closer) Go() *Closer {
	c.once.Do(c.goWaiting)
	return c
}

// Done returns a channel that closes the Closer user has to end.
func (c *Closer) Done() <-chan struct{} {
	dC := c.doneC
	return dC
}

// goWaiting starts wait() as goroutine.
func (c *Closer) goWaiting() {
	go c.wait()
}

// wait is the backend goroutine waiting for closing of
// one of the channels.
func (c *Closer) wait() {
	reflect.Select(c.cases)
	close(c.doneC)
}

// EOF
