// Tideland Go Library - Trace - Monitor
//
// Copyright (C) 2009-2018 Frank Mueller / Tideland / Oldenburg / Germany
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
// MEASURINGS
//--------------------

// Measuring defines one execution time measuring containing the ID and
// the starting time of the measuring and able to pass this data after
// the end of the measuring to the measurer.
type Measuring interface {
	// End ends the measuring and passes it to the measurer.
	End() time.Duration
}

// measuring implements Measuring.
type measuring struct {
	owner *stopWatch
	id    string
	begin time.Time
}

// End implements Measuring.
func (m *measuring) End() time.Duration {
	duration := time.Now().Sub(m.begin)
	m.owner.end(m.id, duration)
	return duration
}

// WatchValue manages the value range for one watch.
type WatchValue struct {
	ID    string
	Count int
	Total time.Duration
	Min   time.Duration
	Max   time.Duration
	Avg   time.Duration
}

// String implements fmt.Stringer.
func (wv *WatchValue) String() string {
	factor := 1000000.0
	total := float64(wv.Total) / factor
	min := float64(wv.Min) / factor
	max := float64(wv.Max) / factor
	avg := float64(wv.Avg) / factor
	return fmt.Sprintf("%s: %d / total %.4f ms / min %.4f ms / max %.4f ms / avg %.4f ms", wv.ID, wv.Count, total, min, max, avg)
}

// update the value with a new measured duration.
func (wv *WatchValue) update(duration time.Duration) {
	// Check for initial values.
	if wv.Count == 0 {
		wv.Count = 1
		wv.Total = duration
		wv.Min = duration
		wv.Max = duration
		wv.Avg = duration
		return
	}
	// Regular update.
	wv.Count++
	wv.Total += duration
	if wv.Min > duration {
		wv.Min = duration
	}
	if wv.Max < duration {
		wv.Max = duration
	}
	wv.Avg = time.Duration(int64(wv.Total) / int64(wv.Count))
}

// WatchValues is a set of values.
type WatchValues []WatchValue

// Implement the sort interface.

func (wvs WatchValues) Len() int           { return len(wvs) }
func (wvs WatchValues) Swap(i, j int)      { wvs[i], wvs[j] = wvs[j], wvs[i] }
func (wvs WatchValues) Less(i, j int) bool { return wvs[i].ID < wvs[j].ID }

//--------------------
// STOP WATCH
//--------------------

// StopWatch allows to measure the execution time of
// code fragments.
type StopWatch interface {
	// Begin starts a new measuring with a given id.
	Begin(id string) Measuring

	// Measure measures the execution time of one function.
	Measure(id string, f func()) time.Duration

	// Read returns the measuring point for an id.
	Read(id string) (WatchValue, error)

	// Do performs the function f for all measuring points.
	Do(f func(WatchValue) error) error
}

// stopWatch implements StopWatch.
type stopWatch struct {
	act        *actor.Actor
	measurings map[string][]time.Duration
	values     map[string]*WatchValue
}

// newStopWatch creates a new stop watch.
func newStopWatch() *stopWatch {
	s := &stopWatch{
		act:        actor.New(actor.WithQueueLen(100)).Go(),
		measurings: make(map[string][]time.Duration),
		values:     make(map[string]*WatchValue),
	}
	go s.ticker()
	return s
}

// Begin implements StopWatch.
func (s *stopWatch) Begin(id string) Measuring {
	return &measuring{
		owner: s,
		id:    id,
		begin: time.Now(),
	}
}

// end returns a measuring to the collected ones.
func (s *stopWatch) end(id string, duration time.Duration) {
	s.act.DoAsync(func() error {
		s.measurings[id] = append(s.measurings[id], duration)
		return nil
	})
}

// Measure implements StopWatch.
func (s *stopWatch) Measure(id string, f func()) time.Duration {
	measuring := &measuring{
		owner: s,
		id:    id,
		begin: time.Now(),
	}
	f()
	return measuring.End()
}

// Read implements StopWatch.
func (s *stopWatch) Read(id string) (WatchValue, error) {
	var wv *WatchValue
	var err error
	s.act.DoSync(func() error {
		s.accumulateOne(id)
		wv = s.values[id]
		if wv == nil {
			err = errors.New(ErrInvalidWatchValue, "watch value '%s' does not exist", id)
		}
		return nil
	})
	if wv == nil {
		return WatchValue{}, err
	}
	return *wv, nil
}

// Do implements StopWatch.
func (s *stopWatch) Do(f func(WatchValue) error) error {
	var err error
	s.act.DoSync(func() error {
		s.accumulateAll()
		for _, wv := range s.values {
			if err = f(*wv); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// reset clears all values.
func (s *stopWatch) reset() {
	s.act.DoAsync(func() error {
		s.measurings = make(map[string][]time.Duration)
		s.values = make(map[string]*WatchValue)
		return nil
	})
}

// stop terminates the indicator.
func (s *stopWatch) stop() error {
	return s.act.Stop(nil)
}

// ticker makes the measurer accumulate all measuring points
// in intervals.
func (s *stopWatch) ticker() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		err := s.act.DoAsync(func() error {
			s.accumulateAll()
			return nil
		})
		if err != nil {
			return
		}
	}
}

// accumulateOne updates one watch value.
func (s *stopWatch) accumulateOne(id string) {
	measurings, ok := s.measurings[id]
	if ok {
		wv := s.values[id]
		if wv == nil {
			wv = &WatchValue{
				ID: id,
			}
			s.values[id] = wv
		}
		for _, duration := range measurings {
			wv.update(duration)
		}
		s.measurings[id] = []time.Duration{}
	}
}

// accumulateAll accumulates all watch values.
func (s *stopWatch) accumulateAll() {
	for id := range s.measurings {
		s.accumulateOne(id)
	}
}

// EOF
