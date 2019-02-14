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

	"tideland.one/go/trace/errors"
)

//--------------------
// SUBSCRIPTION
//--------------------

// Subscription manages a Subscription to Redis channels and allows
// to subscribe and unsubscribe from channels.
type Subscription struct {
	database *Database
	resp     *resp
}

// newSubscription creates a new Subscription.
func newSubscription(db *Database) (*Subscription, error) {
	r, err := newResp(db)
	if err != nil {
		return nil, err
	}
	sub := &Subscription{
		database: db,
		resp:     r,
	}
	// Perform authentication and database selection.
	err = sub.resp.authenticate()
	if err != nil {
		return nil, err
	}
	return sub, nil
}

// Subscribe adds one or more channels to the Subscription.
func (sub *Subscription) Subscribe(channels ...string) error {
	return sub.subUnsub("subscribe", channels...)
}

// Unsubscribe removes one or more channels from the Subscription.
func (sub *Subscription) Unsubscribe(channels ...string) error {
	return sub.subUnsub("unsubscribe", channels...)
}

// subUnsub is the generic Subscription and unsubscription method.
func (sub *Subscription) subUnsub(cmd string, channels ...string) error {
	pattern := false
	args := []interface{}{}
	for _, channel := range channels {
		if containsPattern(channel) {
			pattern = true
		}
		args = append(args, channel)
	}
	if pattern {
		cmd = "p" + cmd
	}
	err := sub.resp.sendCommand(cmd, args...)
	logCommand(cmd, args, err, sub.database.logging)
	return err
}

// Pop waits for a published value and returns it.
func (sub *Subscription) Pop() (PublishedValue, error) {
	result, err := sub.resp.receiveResultSet()
	if err != nil {
		return nil, err
	}
	// Analyse the result.
	kind, err := result.StringAt(0)
	if err != nil {
		return nil, err
	}
	switch {
	case strings.Contains(kind, "message"):
		channel, err := result.StringAt(1)
		if err != nil {
			return nil, err
		}
		value, err := result.ValueAt(2)
		if err != nil {
			return nil, err
		}
		return &publishedValue{
			kind:    kind,
			channel: channel,
			value:   value,
		}, nil
	case strings.Contains(kind, "subscribe"):
		channel, err := result.StringAt(1)
		if err != nil {
			return nil, err
		}
		count, err := result.IntAt(2)
		if err != nil {
			return nil, err
		}
		return &publishedValue{
			kind:    kind,
			channel: channel,
			count:   count,
		}, nil
	default:
		return nil, errors.New(ErrInvalidResponse, msgInvalidResponse, result)
	}
}

// Close ends the Subscription.
func (sub *Subscription) Close() error {
	err := sub.resp.sendCommand("punsubscribe")
	if err != nil {
		return err
	}
	for {
		pv, err := sub.Pop()
		if err != nil {
			return err
		}
		if pv.Kind() == "punsubscribe" {
			break
		}
	}
	return sub.resp.close()
}

// EOF
