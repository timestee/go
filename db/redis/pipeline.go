// Tideland Go Library - Database - Redis Client
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package redis

//--------------------
// IMPORTS
//--------------------

import (
	"strings"

	"tideland.dev/go/trace/errors"
)

//--------------------
// PIPELINE
//--------------------

// Pipeline manages a Redis connection executing pipelined commands.
type Pipeline struct {
	database *Database
	resp     *resp
	counter  int
}

// newPipeline creates a new Pipeline instance.
func newPipeline(db *Database) (*Pipeline, error) {
	r, err := newResp(db)
	if err != nil {
		return nil, err
	}
	ppl := &Pipeline{
		database: db,
		resp:     r,
	}
	// Perform authentication and database selection.
	err = ppl.resp.authenticate()
	if err != nil {
		ppl.Close()
		return nil, err
	}
	err = ppl.resp.selectDatabase()
	if err != nil {
		ppl.Close()
		return nil, err
	}
	return ppl, nil
}

// Do implements the Pipeline interface.
func (ppl *Pipeline) Do(cmd string, args ...interface{}) error {
	cmd = strings.ToLower(cmd)
	if strings.Contains(cmd, "subscribe") {
		return errors.New(ErrUseSubscription, msgUseSubscription)
	}
	err := ppl.resp.sendCommand(cmd, args...)
	logCommand(cmd, args, err, ppl.database.logging)
	if err != nil {
		return err
	}
	ppl.counter++
	return err
}

// Collect implements the Pipeline interface.
func (ppl *Pipeline) Collect() ([]*ResultSet, error) {
	results := []*ResultSet{}
	for i := ppl.counter; i > 0; i-- {
		result, err := ppl.resp.receiveResultSet()
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

// Close implements the Pipeline interface.
func (ppl *Pipeline) Close() error {
	return ppl.resp.close()
}

// EOF
