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
	"sync"
	"time"

	"tideland.one/go/text/etc"
	"tideland.one/go/trace/errors"
)

//--------------------
// DATABASE
//--------------------

// Database controls the parameters for the
// access to a Database.
type Database struct {
	mux        sync.Mutex
	cfg        *Configuration
	address    string
	network    string
	timeout    time.Duration
	index      int
	password   string
	logging    bool
	monitoring bool
}

// Open opens the connection to a Redis Database based on the
// passed configuration.
func Open(cfg *etc.Etc) (*Connection, error) {
	db, err := configureDatabase(cfg)
	if err != nil {
		return nil, err
	}
	return newConnection(db)
}

// OpenPipeline opens the connection to a Redis Database
// especially for a pipeline.
func OpenPipeline(cfg *etc.Etc) (*Pipeline, error) {
	db, err := configureDatabase(cfg)
	if err != nil {
		return nil, err
	}
	return newPipeline(db)
}

// OpenSubscription opens the connection to a Redis Database
// especially for a subscription.
func (db *Database) OpenSubscription(cfg *etc.Etc) (*Subscription, error) {
	db, err := configureDatabase(cfg)
	if err != nil {
		return nil, err
	}
	return newSubscription(db)
}

//--------------------
// PRIVATE
//--------------------

// configureDatabase retrieves the Database settings out
// of a configuration.
func configureDatabase(cfg *etc.Etc) (*Database, error) {
	if cfg == nil {
		return nil, errors.New(ErrNoConfiguration, msgNoConfiguration)
	}
	db := &Database{
		address:  cfg.ValueAsString("address", defaultSocket),
		network:  cfg.ValueAsString("network", defaultNetwork),
		timeout:  cfg.ValueAsDuration("timeout", defaultTimeout),
		index:    cfg.ValueAsInt("index", defaultIndex),
		password: cfg.ValueAsString("password", defaultPassword),
		logging:  cfg.ValueAsBool("logging", defaultLogging),
	}
	return db, nil
}

// EOF
