// Tideland Go Library - Trace - Monitor
//
// Copyright (C) 2009-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package monitor

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"time"

	"tideland.one/go/together/actor"
	"tideland.one/go/trace/errors"
)

//--------------------
// INDICATOR VALUE
//--------------------

// IndicatorValue manages the value range for one indicator.
type IndicatorValue struct {
	ID      string
	Count   int
	Current int
	Min     int
	Max     int
}

// String implements fmt.Stringer.
func (iv *IndicatorValue) String() string {
	return fmt.Sprintf("%s: %d / act %d / min %d / max %d", iv.ID, iv.Count, iv.Current, iv.Min, iv.Max)
}

// update the indicator value.
func (iv *IndicatorValue) update(shallIncr bool) {
	// Check for initial values.
	if iv.Count == 0 {
		iv.Count = 1
		iv.Current = 1
		iv.Min = 1
		iv.Max = 1
	}
	// Regular update.
	iv.Count++
	if shallIncr {
		iv.Current++
	} else {
		iv.Current--
	}
	if iv.Current < iv.Min {
		iv.Min = iv.Current
	}
	if iv.Current > iv.Max {
		iv.Max = iv.Current
	}
}

// IndicatorValues is a set of stay-set values.
type IndicatorValues []IndicatorValue

// Implement the sort interface.

func (ivs IndicatorValues) Len() int           { return len(ivs) }
func (ivs IndicatorValues) Swap(i, j int)      { ivs[i], ivs[j] = ivs[j], ivs[i] }
func (ivs IndicatorValues) Less(i, j int) bool { return ivs[i].ID < ivs[j].ID }

//--------------------
// STAY-SET INDICATOR
//--------------------

// Describing increment or decrement of stay-set values.
const (
	up   = true
	down = false
)

// StaySetIndicator allows to increase and decrease stay-set values.
type StaySetIndicator interface {
	// Increase increases a stay-set staySetIndicator.
	Increase(id string)

	// Decrease decreases a stay-set staySetIndicator.
	Decrease(id string)

	// Read returns a stay-set staySetIndicator.
	Read(id string) (IndicatorValue, error)

	// Do performs the function f for all values.
	Do(f func(IndicatorValue) error) error
}

// staySetIndicator implements StaySetIndicator.
type staySetIndicator struct {
	act     *actor.Actor
	changes map[string][]bool
	values  map[string]*IndicatorValue
}

// newStaySetIndicator creates a new StaySetIndicator.
func newStaySetIndicator() *staySetIndicator {
	i := &staySetIndicator{
		act:     actor.New(actor.WithQueueLen(100)).Go(),
		changes: make(map[string][]bool),
		values:  make(map[string]*IndicatorValue),
	}
	go i.ticker()
	return i
}

// Increase implements StaySetIndicator.
func (i *staySetIndicator) Increase(id string) {
	i.act.DoAsync(func() error {
		i.changes[id] = append(i.changes[id], up)
		return nil
	})
}

// Decrease implements StaySetIndicator.
func (i *staySetIndicator) Decrease(id string) {
	i.act.DoAsync(func() error {
		i.changes[id] = append(i.changes[id], down)
		return nil
	})
}

// Read implemets StaySetIndicator.
func (i *staySetIndicator) Read(id string) (IndicatorValue, error) {
	var iv *IndicatorValue
	var err error
	i.act.DoSync(func() error {
		i.accumulateOne(id)
		iv = i.values[id]
		if iv == nil {
			err = errors.New(ErrInvalidIndicatorValue, "indicator value '%s' does not exist", id)
		}
		return nil
	})
	if iv == nil {
		return IndicatorValue{}, err
	}
	return *iv, nil
}

// Do implements StaySetIndicator.
func (i *staySetIndicator) Do(f func(IndicatorValue) error) error {
	var err error
	i.act.DoSync(func() error {
		i.accumulateAll()
		for _, ssi := range i.values {
			if err = f(*ssi); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// reset clears all values.
func (i *staySetIndicator) reset() {
	i.act.DoAsync(func() error {
		i.changes = make(map[string][]bool)
		i.values = make(map[string]*IndicatorValue)
		return nil
	})
}

// stop terminates the indicator.
func (i *staySetIndicator) stop() error {
	return i.act.Stop(nil)
}

// ticker makes the monitor accumulate all measuring points
// in intervals.
func (i *staySetIndicator) ticker() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		err := i.act.DoAsync(func() error {
			i.accumulateAll()
			return nil
		})
		if err != nil {
			return
		}
	}
}

// accumulateOne updates the indicator value for one ID.
func (i *staySetIndicator) accumulateOne(id string) {
	changes, ok := i.changes[id]
	if ok {
		iv := i.values[id]
		if iv == nil {
			iv = &IndicatorValue{
				ID: id,
			}
			i.values[id] = iv
		}
		for _, increment := range changes {
			iv.update(increment)
		}
		i.changes[id] = []bool{}
	}
}

// accumulateAll updates all indicator values.
func (i *staySetIndicator) accumulateAll() {
	for id := range i.changes {
		i.accumulateOne(id)
	}
}

// EOF
