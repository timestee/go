// Tideland Go Library - Trace - Monitor
//
// Copyright (C) 2009-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package monitor

//--------------------
// MONITOR
//--------------------

// Monitor combines StopWatch and StaySetIndicator.
type Monitor interface {
	// StopWatch returns the internal stop watch instance.
	StopWatch() StopWatch

	// StaySetIndicator returns a stay-set indicator instance.
	StaySetIndicator() StaySetIndicator

	// Reset clears all collected values so far.
	Reset()

	// Stop terminates the monitor.
	Stop()
}

// monitor implements Monitor.
type monitor struct {
	sw  *stopWatch
	ssi *staySetIndicator
}

// New creates a new monitor.
func New() Monitor {
	m := &monitor{
		sw:  newStopWatch(),
		ssi: newStaySetIndicator(),
	}
	return m
}

// StopWatch implements Monitor.
func (m *monitor) StopWatch() StopWatch {
	return m.sw
}

// StaySetIndicator implements Monitor.
func (m *monitor) StaySetIndicator() StaySetIndicator {
	return m.ssi
}

// Reset implements Monitor.
func (m *monitor) Reset() {
	m.sw.reset()
	m.ssi.reset()
}

// Stop implements Monitor.
func (m *monitor) Stop() {
	m.sw.stop()
	m.ssi.stop()
}

// EOF
