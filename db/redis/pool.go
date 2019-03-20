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

	"tideland.dev/go/text/etc"
	"tideland.dev/go/together/wait"
	"tideland.dev/go/trace/errors"
)

//--------------------
// CONSTANTS
//--------------------

const (
	defaultPoolsize = 10
)

//--------------------
// POOL
//--------------------

// Pool manages a number of Redis connections.
type Pool struct {
	mu        sync.Mutex
	closed    bool
	cfg       etc.Etc
	poolsize  int
	available connSet
	inUse     connSet
}

// NewPool creates a redis connection pool.
func NewPool(cfg etc.Etc) *Pool {
	poolsize := cfg.ValueAsInt("poolsize", defaultPoolsize)
	p := &Pool{
		cfg:       cfg,
		closed:    false,
		poolsize:  poolsize,
		available: make(connSet),
		inUse:     make(connSet),
	}
	return p
}

// Pull retrieves a connection our of the pool. If forced is true a
// new one will created if the pool is empty.
func (p *Pool) Pull(forced bool) (conn Connection, err error) {
	// Forced mode.
	if forced {
		return p.pullForced()
	}
	// Normal mode.
	err = wait.WithTimeout(
		context.Background(),
		10*time.Millisecond,
		time.Second,
		func() (bool, error) {
			conn, err = p.pullRegular()
			if err != nil {
				return false, err
			}
			if conn != nil {
				return true, nil
			}
			return false, nil
		},
	)
	return conn, err
}

// pullForced allways opens a new connection and adds it
// to the connections in use.
func (p *Pool) pullForced() (Connection, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil, errors.New(ErrPoolClosed, errorMessages)
	}
	conn, err := redis.Open(p.cfg)
	if err != nil {
		return nil, err
	}
	p.inUse.push(conn)
	return conn, nil
}

// pullRegular tries to get an open connection from the available
// ones. If none is available but there are less in use than allowed
// a new one will be opened.
func (p *Pool) pullRegular() (Connection, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil, errors.New(ErrPoolClosed, errorMessages)
	}
	if len(p.available) == 0 {
		if len(p.inUse) >= p.poolsize {
			// All connections in use.
			return nil, nil
		}
		conn, err := redis.Open(p.cfg)
		if err != nil {
			return nil, err
		}
		p.inUse.push(conn)
		return conn, nil
	}
	conn := p.available.popMove(p.inUse)
	return conn, nil
}

// Push returns a connection back into the pool.
func (p *Pool) Push(conn Connection) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return errors.New(ErrPoolClosed, errorMessages)
	}
	if len(p.available)+len(p.inUse) < p.poolsize {
		p.inUse.pushMove(conn, p.available)
		return nil
	}
	p.inUse.remove(conn)
	return conn.Close()
}

// Close terminates the pool after closing all connections.
func (p *Pool) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.closed {
		return nil
	}
	for conn := range p.available {
		conn.Close()
	}
	for conn := range p.inUse {
		conn.Close()
	}
	return nil
}

//--------------------
// CONNECTION SET
//--------------------

// connSet manages Redis connections.
type connSet map[Connection]struct{}

// pop retrieves the first found connection out of
// the set. If empty nil will be returned.
func (cs connSet) pop() Connection {
	for conn := range cs {
		delete(cs, conn)
		return conn
	}
	return nil
}

// remove removes a connection from the set.
func (cs connSet) remove(conn Connection) {
	delete(cs, conn)
}

// push adds a new connection to the set.
func (cs connSet) push(conn Connection) {
	cs[conn] = struct{}{}
}

// popMove pops a connection from the set and pushes
// it into the target.
func (cs connSet) popMove(target connSet) Connection {
	conn := cs.pop()
	if conn != nil {
		target.push(conn)
	}
	return conn
}

// pushMove removes a connection from the set and pushes
// it into the target.
func (cs connSet) pushMove(conn Connection, target connSet) {
	cs.remove(conn)
	target.push(conn)
}

// EOF
