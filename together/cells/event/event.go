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
	"strconv"
	"time"

	"tideland.dev/go/trace/errors"
)

//--------------------
// DEFAULTS
//--------------------

const (
	DefaultValue = true
)

var timeFormats = []string{
	"Mon Jan 2 15:04:05 -0700 MST 2006",
	"2006-01-02 15:04:05.999999999 -0700 MST",
	time.ANSIC,
	time.Kitchen,
	time.RFC1123,
	time.RFC1123Z,
	time.RFC3339,
	time.RFC3339Nano,
	time.RFC822,
	time.RFC822Z,
	time.RFC850,
	time.RubyDate,
	time.Stamp,
	time.StampMicro,
	time.StampMilli,
	time.StampNano,
	time.UnixDate,
}

//--------------------
// PAYLOAD
//--------------------

// Payload contains key/value pairs of data.
type Payload struct {
	values map[string]interface{}
}

func NewPayload(kvs ...interface{}) *Payload {
	p := &Payload{
		values: map[string]interface{}{},
	}
	var key string
	for i, kv := range kvs {
		if i%2 == 0 {
			key = fmt.Sprintf("%v", kv)
			p.values[key] = DefaultValue
			continue
		}
		p.values[key] = kv
	}
	return p
}

// String tries to interpred the keyed payload as string.
func (p *Payload) String(key string) (string, error) {
	v, ok := p.values[key]
	if !ok {
		return "", errors.New(ErrNoValue, msgNoValue, key)
	}
	tv, ok := v.(string)
	if ok {
		return tv, nil
	}
	return fmt.Sprintf("%v", v), nil
}

// Bool tries to interpred the keyed payload as bool.
func (p *Payload) Bool(key string) (bool, error) {
	v, ok := p.values[key]
	if !ok {
		return false, errors.New(ErrNoValue, msgNoValue, key)
	}
	tv, ok := v.(bool)
	if ok {
		return tv, nil
	}
	bv, err := strconv.ParseBool(fmt.Sprintf("%v", v))
	if err != nil {
		return false, errors.Annotate(err, ErrConverting, msgConverting, key, "bool")
	}
	return bv, nil
}

// Float64 tries to interpred the keyed payload as float64.
func (p *Payload) Float64(key string) (float64, error) {
	v, ok := p.values[key]
	if !ok {
		return 0.0, errors.New(ErrNoValue, msgNoValue, key)
	}
	tv, ok := v.(float64)
	if ok {
		return tv, nil
	}
	fv, err := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
	if err != nil {
		return 0.0, errors.Annotate(err, ErrConverting, msgConverting, key, "float64")
	}
	return fv, nil
}

// Int tries to interpred the keyed payload as int.
func (p *Payload) Int(key string) (int, error) {
	v, ok := p.values[key]
	if !ok {
		return 0, errors.New(ErrNoValue, msgNoValue, key)
	}
	tv, ok := v.(int)
	if ok {
		return tv, nil
	}
	iv, err := strconv.Atoi(fmt.Sprintf("%v", v))
	if err != nil {
		return 0, errors.Annotate(err, ErrConverting, msgConverting, key, "int")
	}
	return iv, nil
}

// Int64 tries to interpred the keyed payload as int64.
func (p *Payload) Int64(key string) (int64, error) {
	v, ok := p.values[key]
	if !ok {
		return 0, errors.New(ErrNoValue, msgNoValue, key)
	}
	tv, ok := v.(int64)
	if ok {
		return tv, nil
	}
	iv, err := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
	if err != nil {
		return 0, errors.Annotate(err, ErrConverting, msgConverting, key, "int64")
	}
	return iv, nil
}

// Uint64 tries to interpred the keyed payload as uint64.
func (p *Payload) Uint64(key string) (uint64, error) {
	v, ok := p.values[key]
	if !ok {
		return 0, errors.New(ErrNoValue, msgNoValue, key)
	}
	tv, ok := v.(uint64)
	if ok {
		return tv, nil
	}
	iv, err := strconv.ParseUint(fmt.Sprintf("%v", v), 10, 64)
	if err != nil {
		return 0, errors.Annotate(err, ErrConverting, msgConverting, key, "uint64")
	}
	return iv, nil
}

// Time tries to interpred the keyed payload as time.Time.
func (p *Payload) Time(key string) (time.Time, error) {
	v, ok := p.values[key]
	if !ok {
		return time.Time{}, errors.New(ErrNoValue, msgNoValue, key)
	}
	tv, ok := v.(time.Time)
	if ok {
		return tv, nil
	}
	vs := fmt.Sprintf("%v", v)
	for _, timeFormat := range timeFormats {
		iv, err := time.Parse(timeFormat, vs)
		if err == nil {
			return iv, nil
		}
	}
	return time.Time{}, errors.New(ErrConverting, msgConverting, key, "time.Time")
}

// Duration tries to interpred the keyed payload as time.Duration.
func (p *Payload) Duration(key string) (time.Duration, error) {
	v, ok := p.values[key]
	if !ok {
		return 0, errors.New(ErrNoValue, msgNoValue, key)
	}
	tv, ok := v.(time.Duration)
	if ok {
		return tv, nil
	}
	iv, err := time.ParseDuration(fmt.Sprintf("%v", v))
	if err != nil {
		return 0, errors.Annotate(err, ErrConverting, msgConverting, key, "time.Duration")
	}
	return iv, nil
}

// Clone creates a new payload with the content of the current one and
// applies the given changes.
func (p *Payload) Clone(kvs ...interface{}) *Payload {
	cp := &Payload{
		values: map[string]interface{}{},
	}
	for key, value := range p.values {
		cp.values[key] = value
	}
	var key string
	for i, kv := range kvs {
		if i%2 == 0 {
			key = fmt.Sprintf("%v", kv)
			cp.values[key] = DefaultValue
			continue
		}
		cp.values[key] = kv
	}
	return cp
}

//--------------------
// EVENT
//--------------------

// Event describes an event of the cells. It contains a topic as well as
// a possible number of key/value pairs as payload.
type Event struct {
	timestamp time.Time
	topic     string
	payload   *Payload
}

// New creates a new event. The arguments after the topic are taken
// to create a new payload.
func New(topic string, kvs ...interface{}) *Event {
	return &Event{
		timestamp: time.Now(),
		topic:     topic,
		payload:   NewPayload(kvs...),
	}
}

// WithPayload creates a new event with a given external payload.
func WithPayload(topic string, p *Payload) *Event {
	return &Event{
		timestamp: time.Now(),
		topic:     topic,
		payload:   p,
	}
}

// Timestamp returns the event timestamp.
func (e *Event) Timestamp() time.Time {
	return e.timestamp
}

// Topic returns the event topic.
func (e *Event) Topic() string {
	return e.topic
}

// Payload returns the event payload.
func (e *Event) Payload() *Payload {
	return e.payload
}

// EOF
