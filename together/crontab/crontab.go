// Tideland Go Library - Together - CronTab
//
// Copyright (C) 2017-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package crontab

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"tideland.one/go/together/actor"
	"tideland.one/go/together/loop"
	"tideland.one/go/together/notifier"
)

//--------------------
// GLOBAL
//--------------------

var (
	mu sync.Mutex
	ct *crontab
)

//--------------------
// CRONTAB
//--------------------

// crontab implements the tanle for all cronjobs.
type crontab struct {
	actor *actor.Actor
	jobs  map[string]*cronjob
}

// goCrontab instantiates and starts the crontab if it is not already running.
func goGrontab() {
	mu.Lock()
	defer mu.Unlock()
	if ct != nil {
		return
	}
	ct = &crontab{
		actor: actor.New().Go(),
		jobs:  make(map[string]*cronjob),
	}
}

//--------------------
// CRONJOB
//--------------------

// cronjob is responsible to run one job.
type cronjob struct {
	id       string
	start    *time.Time
	interval *time.Duration
	job      func() error
	loop     *loop.Loop
	notifier *notifier.Notifier
	rs       loop.Reasons
}

// newCronjob creates a new cronjob and starts its goroutine.
func newCronjob(id string, s *time.Time, i *time.Duration, j func() error) *cronjob {
	cj := &cronjob{
		id:       id,
		start:    s,
		interval: i,
		job:      j,
		notifier: notifier.New(),
		rs:       loop.MakeReasons(),
	}
	cj.loop = loop.New(cj.worker,
		loop.WithRecoverer(cj.recoverer),
		loop.WithNotifier(cj.notifier)).Go()
	<-cj.notifier.Working()
	return cj
}

// stop ends the cronjob goroutine.
func (cj *cronjob) stop() error {
	err := cj.loop.Stop(nil)
	if loop.IsErrLoopNotWorking(err) {
		return nil
	}
	return err
}

// status returns the status of the cronjob.
func (cj *cronjob) status() notifier.Status {
	return cj.loop.Status()
}

// worker runs the cronjob.
func (cj *cronjob) worker(c *notifier.Closer) error {
	// Init.
	var interval time.Duration
	if cj.start != nil {
		interval = time.Until(*cj.start)
	} else {
		interval = *cj.interval
	}
	// Loop.
	for {
		select {
		case <-c.Done():
			return nil
		case <-time.After(interval):
			if err := cj.job(); err != nil {
				return err
			}
			if cj.interval == nil {
				// Only once.
				return nil
			}
			// In intervals.
			interval = *cj.interval
		}
	}
}

// recoverer allows the cronjob to survive panics.
func (cj *cronjob) recoverer(reason interface{}) error {
	cj.rs = cj.rs.Append(reason)
	if cj.rs.Frequency(5, 10*time.Millisecond) {
		return fmt.Errorf("too high error frequency: %v", cj.rs)
	}
	if cj.rs.Len() >= 10 {
		return fmt.Errorf("too many errors: %v", cj.rs)
	}
	return nil
}

//--------------------
// API
//--------------------

// SubmitAt adds a function running only once at a given time.
func SubmitAt(id string, at time.Time, j func() error) error {
	goGrontab()
	var err error
	if actErr := ct.actor.DoSync(func() error {
		if ct.jobs[id] != nil {
			err = fmt.Errorf("job ID '%s' already exists", id)
			return nil
		}
		ct.jobs[id] = newCronjob(id, &at, nil, j)
		return nil
	}); actErr != nil {
		return actErr
	}
	return err
}

// SubmitEvery adds a function running every interval.
func SubmitEvery(id string, every time.Duration, j func() error) error {
	goGrontab()
	var err error
	if actErr := ct.actor.DoSync(func() error {
		if ct.jobs[id] != nil {
			err = fmt.Errorf("job ID '%s' already exists", id)
			return nil
		}
		ct.jobs[id] = newCronjob(id, nil, &every, j)
		return nil
	}); actErr != nil {
		return actErr
	}
	return err
}

// SubmitAtEvery adds a function running every interval starting at a given time.
func SubmitAtEvery(id string, at time.Time, every time.Duration, j func() error) error {
	goGrontab()
	var err error
	if actErr := ct.actor.DoSync(func() error {
		if ct.jobs[id] != nil {
			err = fmt.Errorf("job ID '%s' already exists", id)
			return nil
		}
		ct.jobs[id] = newCronjob(id, &at, &every, j)
		return nil
	}); actErr != nil {
		return actErr
	}
	return err
}

// SubmitAfterEvery adds a function running every interval after a given pause.
func SubmitAfterEvery(id string, pause, every time.Duration, j func() error) error {
	return SubmitAtEvery(id, time.Now().Add(pause), every, j)
}

// List returns all currently submitted IDs.
func List() ([]string, error) {
	goGrontab()
	var ids []string
	var err error
	if actErr := ct.actor.DoSync(func() error {
		for id := range ct.jobs {
			ids = append(ids, id)
		}
		return nil
	}); actErr != nil {
		return ids, actErr
	}
	sort.Strings(ids)
	return ids, err
}

// Status returns the status of a cronjob.
func Status(id string) (notifier.Status, error) {
	goGrontab()
	var status notifier.Status
	var err error
	if actErr := ct.actor.DoSync(func() error {
		if job, ok := ct.jobs[id]; ok {
			status = job.status()
			return nil
		}
		err = fmt.Errorf("job ID '%s' does not exist", id)
		return nil
	}); actErr != nil {
		return notifier.Unknown, actErr
	}
	return status, err
}

// Revoke stops a cronjob and removes it from the table.
func Revoke(id string) error {
	goGrontab()
	var err error
	if actErr := ct.actor.DoSync(func() error {
		if job, ok := ct.jobs[id]; ok {
			delete(ct.jobs, id)
			err = job.stop()
			return nil
		}
		err = fmt.Errorf("job ID '%s' does not exist", id)
		return nil
	}); actErr != nil {
		return actErr
	}
	return err
}

// EOF
