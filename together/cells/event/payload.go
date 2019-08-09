// Tideland Go Library - Together - Cells - Event
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license

package event // import "tideland.dev/go/together/cells/event"

//--------------------
// IMPORTS
//--------------------

import (
	"fmt"
	"reflect"

	"tideland.dev/go/trace/failure"
)

//--------------------
// PAYLOAD CHANNEL
//--------------------

// PayloadChan is intended to be sent with an event as payload
// so that a behavior can use it to answer a request.
type PayloadChan chan *Payload

//--------------------
// PAYLOAD
//--------------------

// Payload contains key/value pairs of data.
type Payload struct {
	values map[string]interface{}
	replyc PayloadChan
}

// NewPayload creates a new payload with the given pairs of
// keys and values. In case of payloads as values their keys
// will be joined with their higher level keys by a slash
// and so later can be accessed by "{key}/{payload-key}", even
// multiple levels deep.
func NewPayload(kvs ...interface{}) *Payload {
	pl := &Payload{
		values: map[string]interface{}{},
	}
	pl.setKeyValues(kvs...)
	return pl
}

// NewReplyPayload creates a new payload like NewPayload() but
// also a reply channel.
func NewReplyPayload(kvs ...interface{}) (*Payload, PayloadChan) {
	pl := &Payload{
		values: map[string]interface{}{},
		replyc: make(PayloadChan, 1),
	}
	pl.setKeyValues(kvs...)
	return pl, pl.replyc
}

// Keys returns the keys of the payload.
func (pl *Payload) Keys() []string {
	var keys []string
	for key := range pl.values {
		keys = append(keys, key)
	}
	return keys
}

// At returns the value at the given key. This value may
// be empty.
func (pl *Payload) At(key string) *Value {
	v, ok := pl.values[key]
	if !ok {
		return &Value{
			err: failure.New("no payload value at key %q", key),
		}
	}
	return &Value{
		raw: v,
	}
}

// Do performs a function for all key/value pairs.
func (pl *Payload) Do(f func(key string, value *Value) error) error {
	var errs []error
	for key, rawValue := range pl.values {
		value := &Value{
			raw: rawValue,
		}
		errs = append(errs, f(key, value))
	}
	return failure.Collect(errs...)
}

// Reply allows the receiver of a payload to reply via a channel.
func (pl *Payload) Reply(rpl *Payload) error {
	if pl.replyc == nil {
		return failure.New("payload contains no reply channel")
	}
	select {
	case pl.replyc <- rpl:
		return nil
	default:
		return failure.New("payload reply channel is closed")
	}
}

// Clone creates a new payload with the content of the current one and
// applies the given changes.
func (pl *Payload) Clone(kvs ...interface{}) *Payload {
	plc := &Payload{
		values: map[string]interface{}{},
	}
	for key, value := range pl.values {
		plc.values[key] = value
	}
	plc.setKeyValues(kvs...)
	return plc
}

// setKeyValues iterates over the key/value values and adds
// them to the payloads values.
func (pl *Payload) setKeyValues(kvs ...interface{}) {
	var key string
	for i, kv := range kvs {
		if i%2 == 0 {
			// Talking about a key.
			if plk, ok := kv.(*Payload); ok {
				// A payload, merge values.
				pl.mergeMap(plk.values)
				pl.replyc = plk.replyc
				continue
			}
			if reflect.TypeOf(kv).Kind() == reflect.Map {
				// A map, merge it.
				pl.mergeMap(kv)
				continue
			}
			// Any other key.
			key = fmt.Sprintf("%v", kv)
			pl.values[key] = DefaultValue
			continue
		}
		// Talking about a value.
		if reflect.TypeOf(kv).Kind() == reflect.Map {
			// It's a map, take it as nested payload.
			pl.values[key] = NewPayload(kv)
			continue
		}
		// Any other value.
		pl.values[key] = kv
	}
}

// mergeMap maps any map into the own values.
func (pl *Payload) mergeMap(m interface{}) {
	var kvs []interface{}
	iter := reflect.ValueOf(m).MapRange()
	for iter.Next() {
		key := iter.Key().Interface()
		value := iter.Value().Interface()
		kvs = append(kvs, key, value)
	}
	pl.setKeyValues(kvs...)
}

// EOF
