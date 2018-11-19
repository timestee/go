// Tideland Go Library - Network - JSON Web Token - Cache
//
// Copyright (C) 2016-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package cache

//--------------------
// IMPORTS
//--------------------

import (
	"sync"
	"time"

	"tideland.one/go/net/jwt/token"
	"tideland.one/go/together/loop"
	"tideland.one/go/together/notifier"
)

//--------------------
// CACHE
//--------------------

// Cache provides a caching for tokens so that these
// don't have to be decoded or verified multiple times.
type Cache interface {
	// Get tries to retrieve a token from the cache.
	Get(token string) (token.JWT, bool)

	// Put adds a token to the cache.
	Put(jwt token.JWT) int

	// Cleanup manually tells the cache to cleanup.
	Cleanup()

	// Stop tells the cache to end working.
	Stop() error
}

// cacheEntry manages a token and its access time.
type cacheEntry struct {
	jwt      token.JWT
	accessed time.Time
}

// cache implements Cache.
type cache struct {
	mu         sync.Mutex
	entries    map[string]*cacheEntry
	ttl        time.Duration
	leeway     time.Duration
	interval   time.Duration
	maxEntries int
	cleanupc   chan time.Duration
	loop       *loop.Loop
}

// New creates a new JWT caching. The ttl value controls
// the time a cached token may be unused before cleanup. The
// leeway is used for the time validation of the token itself.
// The duration of the interval controls how often the background
// cleanup is running. Final configuration parameter is the maximum
// number of entries inside the cache. If these grow too fast the
// ttl will be temporarily reduced for cleanup.
func New(ttl, leeway, interval time.Duration, maxEntries int) Cache {
	c := &cache{
		entries:    map[string]*cacheEntry{},
		ttl:        ttl,
		leeway:     leeway,
		interval:   interval,
		maxEntries: maxEntries,
		cleanupc:   make(chan time.Duration, 5),
	}
	c.loop = loop.New(c.backendLoop).Go()
	return c
}

// Get implements the Cache interface.
func (c *cache) Get(token string) (token.JWT, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[token]
	if !ok {
		return nil, false
	}
	if entry.jwt.IsValid(c.leeway) {
		entry.accessed = time.Now()
		return entry.jwt, true
	}
	// Remove invalid token.
	delete(c.entries, token)
	return nil, false
}

// Put implements the Cache interface.
func (c *cache) Put(jwt token.JWT) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	if jwt.IsValid(c.leeway) {
		c.entries[jwt.String()] = &cacheEntry{jwt, time.Now()}
		lenEntries := len(c.entries)
		if lenEntries > c.maxEntries {
			ttl := int64(c.ttl) / int64(lenEntries) * int64(c.maxEntries)
			c.cleanupc <- time.Duration(ttl)
		}
	}
	return len(c.entries)
}

// Cleanup implements the Cache interface.
func (c *cache) Cleanup() {
	c.cleanupc <- c.ttl
}

// Stop implements the Cache interface.
func (c *cache) Stop() error {
	return c.loop.Stop(nil)
}

// backendLoop runs a cleaning session every five minutes.
func (c *cache) backendLoop(nc *notifier.Closer) error {
	defer func() {
		// Cleanup entries map after stop or error.
		c.entries = nil
	}()
	ticker := time.NewTicker(c.interval)
	for {
		select {
		case <-nc.Done():
			return nil
		case ttl := <-c.cleanupc:
			c.cleanup(ttl)
		case <-ticker.C:
			c.cleanup(c.ttl)
		}
	}
}

// cleanup checks for invalid or unused tokens.
func (c *cache) cleanup(ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	valids := map[string]*cacheEntry{}
	now := time.Now()
	for token, entry := range c.entries {
		if entry.jwt.IsValid(c.leeway) {
			if entry.accessed.Add(ttl).After(now) {
				// Everything fine.
				valids[token] = entry
			}
		}
	}
	c.entries = valids
}

// EOF
