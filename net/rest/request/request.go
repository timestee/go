// Tideland Go Library - Network - REST - Request
//
// Copyright (C) 2009-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package request

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"tideland.one/go/trace/errors"
)

//--------------------
// SERVERS
//--------------------

// key is to address the servers inside a context.
type key int

const (
	serversKey key = iota
)

// server contains the configuration of one server.
type server struct {
	URL       string
	Transport *http.Transport
}

// Servers maps IDs of domains to their server configurations.
// Multiple ones can be added per domain for spreading the
// load or provide higher availability.
type Servers struct {
	mu      sync.RWMutex
	servers map[string][]*server
}

// NewServers creates a new servers manager.
func NewServers() *Servers {
	rand.Seed(time.Now().Unix())
	return &Servers{
		servers: make(map[string][]*server),
	}
}

// Add adds a domain server configuration.
func (s *Servers) Add(domain, url string, transport *http.Transport) {
	s.mu.Lock()
	defer s.mu.Unlock()
	srvs, ok := s.servers[domain]
	if ok {
		s.servers[domain] = append(srvs, &server{url, transport})
		return
	}
	s.servers[domain] = []*server{{url, transport}}
}

// Caller retrieves a caller for a domain.
func (s *Servers) Caller(domain string) (*Caller, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	srvs, ok := s.servers[domain]
	if !ok {
		return nil, errors.New(ErrNoServerDefined, msgNoServerDefined, domain)
	}
	return newCaller(domain, srvs), nil
}

// NewContext returns a new context that carries configured servers.
func NewContext(ctx context.Context, servers *Servers) context.Context {
	return context.WithValue(ctx, serversKey, servers)
}

// FromContext returns the servers configuration stored in ctx, if any.
func FromContext(ctx context.Context) (*Servers, bool) {
	servers, ok := ctx.Value(serversKey).(*Servers)
	return servers, ok
}

// EOF
