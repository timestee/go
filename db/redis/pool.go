// Tideland Go Library - DB - Redis Client
//
// Copyright (C) 2009-2019 Frank Mueller / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package redis

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"time"

	"tideland.dev/go/together/actor"
	"tideland.dev/go/together/wait"
)

//--------------------
// CONNECTION POOL
//--------------------

const (
	forcedPull   = true
	unforcedPull = false

	forcedPullRequest = iota
	unforcedPullRequest
	pushRequest
	killRequest
	closeRequest
)

// pool manages a number of Redis resp instances.
type pool struct {
	ctx       context.Context
	database  *Database
	available map[*resp]*resp
	inUse     map[*resp]*resp
	act       *actor.Actor
}

// newPool creates a connection pool with uninitialized
// protocol instances.
func newPool(db *Database) *pool {
	p := &pool{
		database:  db,
		available: make(map[*resp]*resp),
		inUse:     make(map[*resp]*resp),
		act:       actor.New(actor.WithContext(db.ctx)).Go(),
	}
	return p
}

func (p *pool) stop() error {
	return p.act.Stop(p.close())
}

// pullForced retrieves a new created protocol.
func (p *pool) pullForced() (resp *resp, err error) {
	if aerr := p.act.DoSync(func() error {
		resp, err = newResp(p.database)
		if err == nil {
			p.inUse[resp] = resp
		}
		return nil
	}); aerr != nil {
		return nil, aerr
	}
	return
}

// pullRetry retrieves a protocol out of the pool. It tries to
// do it multiple times.
func (p *pool) pullRetry() (resp *resp, err error) {
	if werr := wait.WithTimeout(
		p.database.ctx,
		10*time.Millisecond,
		30*time.Second,
		func() (bool, error) {
			resp, err = p.pull()
			if err != nil {
				return false, err
			}
			if resp != nil {
				return true, nil
			}
			return false, nil
		},
	); werr != nil {
		return nil, werr
	}
	return
}

// pull retrieves a protocol out of the pool.
func (p *pool) pull() (resp *resp, err error) {
	if aerr := p.act.DoSync(func() error {
		switch {
		case len(p.available) > 0:
		fetched:
			for resp = range p.available {
				delete(p.available, resp)
				p.inUse[resp] = resp
				break fetched
			}
		case len(p.inUse) < p.database.poolsize:
			resp, err = newResp(p.database)
			if err != nil {
				p.inUse[resp] = resp
			}
		}
		return nil
	}); aerr != nil {
		return nil, aerr
	}
	return
}

// push returns a protocol back into the pool.
func (p *pool) push(resp *resp) (err error) {
	if aerr := p.act.DoSync(func() error {
		delete(p.inUse, resp)
		if len(p.available) < p.database.poolsize {
			p.available[resp] = resp
		} else {
			err = resp.close()
		}
		return nil
	}); aerr != nil {
		return aerr
	}
	return
}

// kill closes the connection and removes it from the pool.
func (p *pool) kill(resp *resp) (err error) {
	if aerr := p.act.DoSync(func() error {
		delete(p.inUse, resp)
		err = resp.close()
		return nil
	}); aerr != nil {
		return aerr
	}
	return
}

// close closes all pooled protocol instances, first the available ones,
// then the ones in use.
func (p *pool) close() (err error) {
	if aerr := p.act.DoSync(func() error {
		for resp := range p.available {
			resp.close()
		}
		for resp := range p.inUse {
			resp.close()
		}
		return nil
	}); aerr != nil {
		return aerr
	}
	return
}

// EOF
