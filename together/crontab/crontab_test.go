// Tideland Go Library - Together - CronTab - Unit Tests
//
// Copyright (C) 2017-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package crontab_test

//--------------------
// IMPORTS
//--------------------

import (
	"errors"
	"testing"
	"time"

	"tideland.one/go/audit/asserts"
	"tideland.one/go/together/crontab"
	"tideland.one/go/together/notifier"
)

//--------------------
// TESTS
//--------------------

// TestSubmitStatusRevoke tests a simple submitting, status retrieval,
// and revoking.
func TestSubmitStatusRevoke(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	at := time.Now().Add(250 * time.Millisecond)
	beenThereDoneThat := false

	// Test.
	err := crontab.SubmitAt("foo-1", at, func() error {
		beenThereDoneThat = true
		return nil
	})
	assert.NoError(err)
	status, err := crontab.Status("foo-1")
	assert.NoError(err)
	assert.Equal(status, notifier.Working)
	err = crontab.Revoke("foo-1")
	assert.NoError(err)
	time.Sleep((500 * time.Millisecond))
	assert.False(beenThereDoneThat)

	err = crontab.SubmitAt("foo-2", at, func() error {
		return errors.New("ouch")
	})
	assert.NoError(err)
	time.Sleep(time.Second)
	err = crontab.Revoke("foo-2")
	assert.ErrorMatch(err, "ouch")
}

// TestList tests the listing of submitted jobs.
func TestList(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	at := time.Now().Add(time.Second)

	// Test.
	jobs, err := crontab.List()
	assert.NoError(err)
	assert.Empty(jobs)

	err = crontab.SubmitAt("fuddel-1", at, func() error {
		return nil
	})
	assert.NoError(err)
	jobs, err = crontab.List()
	assert.NoError(err)
	assert.Length(jobs, 1)

	err = crontab.SubmitAt("fuddel-2", at, func() error {
		return nil
	})
	assert.NoError(err)
	jobs, err = crontab.List()
	assert.NoError(err)
	assert.Length(jobs, 2)
	assert.Contents("fuddel-1", jobs)
	assert.Contents("fuddel-2", jobs)

	err = crontab.Revoke("fuddel-1")
	assert.NoError(err)
	jobs, err = crontab.List()
	assert.NoError(err)
	assert.Length(jobs, 1)
	assert.False(func() bool {
		for _, job := range jobs {
			if job == "fuddel-1" {
				return true
			}
		}
		return false
	}())
	assert.Contents("fuddel-2", jobs)

	err = crontab.Revoke("fuddel-2")
	assert.NoError(err)
	jobs, err = crontab.List()
	assert.NoError(err)
	assert.Empty(jobs)
}

// TestSubmitAt tests if a job is executed only once.
func TestSubmitAt(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	atOne := time.Now().Add(100 * time.Millisecond)
	oneDiff := 0 * time.Millisecond
	oneDiffExtend := float64(10 * time.Millisecond)
	atTwo := atOne.Add(200 * time.Millisecond)
	count := 0

	// Test.
	err := crontab.SubmitAt("bar-1", atOne, func() error {
		oneDiff = time.Now().Sub(atOne)
		count++
		return nil
	})
	assert.NoError(err)
	err = crontab.SubmitAt("bar-2", atTwo, func() error {
		count *= 5
		return nil
	})
	assert.NoError(err)
	time.Sleep(200 * time.Millisecond)
	assert.About(float64(oneDiff), 0.0, oneDiffExtend)
	assert.Equal(count, 1)
	time.Sleep(time.Second)
	assert.Equal(count, 5)

	status, err := crontab.Status("bar-1")
	assert.NoError(err)
	assert.Equal(status, notifier.Stopped)
	status, err = crontab.Status("bar-2")
	assert.NoError(err)
	assert.Equal(status, notifier.Stopped)

	err = crontab.Revoke("bar-1")
	assert.NoError(err)
	err = crontab.Revoke("bar-2")
	assert.NoError(err)
}

// TestSubmitEvery tests if a job is executed every given interval.
func TestSubmitEvery(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	every := 200 * time.Millisecond
	count := 0

	// Test.
	err := crontab.SubmitEvery("baz-1", every, func() error {
		count++
		return nil
	})
	assert.NoError(err)
	time.Sleep(700 * time.Millisecond)
	assert.Equal(count, 3)
	time.Sleep(400 * time.Millisecond)
	assert.Equal(count, 5)

	err = crontab.Revoke("baz-1")
	assert.NoError(err)
}

// TestSubmitAtEvery tests if a job is executed every given interval
// after a given time.
func TestSubmitAtEvery(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	at := time.Now().Add(500 * time.Millisecond)
	every := 100 * time.Millisecond
	count := 0

	// Test.
	err := crontab.SubmitAtEvery("babbel-1", at, every, func() error {
		count++
		return nil
	})
	assert.NoError(err)
	time.Sleep(750 * time.Millisecond)
	assert.Equal(count, 3)

	err = crontab.Revoke("babbel-1")
	assert.NoError(err)
}

// TestSubmitAfterEvery tests if a job is executed every given interval
// after a given pause.
func TestSubmitAfterEvery(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	pause := 500 * time.Millisecond
	every := 100 * time.Millisecond
	count := 0

	// Test.
	err := crontab.SubmitAfterEvery("daddel-1", pause, every, func() error {
		count++
		return nil
	})
	assert.NoError(err)
	time.Sleep(750 * time.Millisecond)
	assert.Equal(count, 3)

	err = crontab.Revoke("daddel-1")
	assert.NoError(err)
}

// TestIllegal tests double id submit and illegal id revoke.
func TestIllegal(t *testing.T) {
	// Init.
	assert := asserts.NewTesting(t, true)
	at := time.Now().Add(time.Second)
	job := func() error {
		return nil
	}

	// Test.
	err := crontab.SubmitAt("yadda-1", at, job)
	assert.NoError(err)
	err = crontab.SubmitAt("yadda-1", at, job)
	assert.ErrorMatch(err, `job ID 'yadda-1' already exists`)

	err = crontab.Revoke("yadda-2")
	assert.ErrorMatch(err, `job ID 'yadda-2' does not exist`)
}

// EOF
