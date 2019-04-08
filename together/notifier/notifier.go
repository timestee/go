// Tideland Go Library - Together - Notifier
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package notifier

//--------------------
// IMPORTS
//--------------------

import (
	"sync"
)

//--------------------
// STATUS
//--------------------

// Status describes the status of a background groroutine.
type Status int

// Different statuses of a background goroutine.
const (
	Unknown Status = iota
	Starting
	Ready
	Working
	Stopping
	Stopped
)

// statusStr contains the string representation of a status.
var statusStr = map[Status]string{
	Unknown:  "unknown",
	Starting: "starting",
	Ready:    "ready",
	Working:  "working",
	Stopping: "stopping",
	Stopped:  "stopped",
}

// String implements the fmt.Stringer interface.
func (s Status) String() string {
	if str, ok := statusStr[s]; ok {
		return str
	}
	return "invalid"
}

//--------------------
// NOTIFIER
//--------------------

// Notifier allows code to be notified about the internal
// status of a background goroutine.
type Notifier struct {
	mu        sync.Mutex
	status    Status
	readyC    chan struct{}
	workingC  chan struct{}
	stoppingC chan struct{}
	stoppedC  chan struct{}
}

// New creates a new Notifier instance.
func New() *Notifier {
	return &Notifier{
		status:    Starting,
		readyC:    make(chan struct{}),
		workingC:  make(chan struct{}),
		stoppingC: make(chan struct{}),
		stoppedC:  make(chan struct{}),
	}
}

// Status returns the current goroutine status.
func (n *Notifier) Status() Status {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.status
}

// Ready allows to wait until a goroutine has
// reached Ready status.
func (n *Notifier) Ready() <-chan struct{} {
	return n.readyC
}

// Working allows to wait until a goroutine has
// reached Working status.
func (n *Notifier) Working() <-chan struct{} {
	return n.workingC
}

// Stopping allows to wait until a goroutine has
// reached Stopping status.
func (n *Notifier) Stopping() <-chan struct{} {
	return n.stoppingC
}

// Stopped allows to wait until a goroutine has
// reached Stopped status.
func (n *Notifier) Stopped() <-chan struct{} {
	return n.stoppedC
}

// notify sets the new status and informs listener.
func (n *Notifier) notify(status Status) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if status <= n.status {
		return
	}
	n.status = status
	switch status {
	case Ready:
		close(n.readyC)
	case Working:
		close(n.workingC)
	case Stopping:
		close(n.stoppingC)
	case Stopped:
		close(n.stoppedC)
	default:
		panic("invalid loop status")
	}
}

//--------------------
// BUNDLE
//--------------------

// Bundle manages a set of Notifiers.
type Bundle struct {
	mu        sync.Mutex
	status    Status
	notifiers []*Notifier
	stoppedC  chan struct{}
}

// NewBundle creates a bundle of Notifiers.
func NewBundle() *Bundle {
	return &Bundle{
		status:   Starting,
		stoppedC: make(chan struct{}),
	}
}

// Add appends on or more Notifiers to the Bundle.
func (b *Bundle) Add(ns ...*Notifier) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.notifiers = append(b.notifiers, ns...)
}

// Notify sets the new status and informs all Notifiers.
func (b *Bundle) Notify(status Status) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if status <= b.status {
		return
	}
	b.status = status
	for _, n := range b.notifiers {
		n.notify(status)
	}
	if status == Stopped {
		close(b.stoppedC)
	}
}

// Status returns the current Status.
func (b *Bundle) Status() Status {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.status
}

// Stopped informs the owner of the Bundle that the status is now Stopped.
func (b *Bundle) Stopped() <-chan struct{} {
	return b.stoppedC
}

// EOF
