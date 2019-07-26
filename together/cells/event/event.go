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
// CONSTANTS
//--------------------

// Standard topics.
const (
	TopicCollected = "collected"
	TopicCounted   = "counted"
	TopicProcess   = "process"
	TopicProcessed = "processed"
	TopicReset     = "reset"
	TopicResult    = "result"
	TopicStatus    = "status"
	TopicTick      = "tick"
)

// Standard keys.
const (
	// Some standard keys.
	KeyReply = "reply"
)

// Standard values.
const (
	DefaultValue = true
)

// Different formats for the parsing of strings
// into times.
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
// PAYLOAD CHANNEL
//--------------------

// PayloadChan is intended to be sent with an event as payload
// so that a behavior can use it to answer a request.
type PayloadChan chan *Payload

// MakePayloadChan create a single buffered payload channel.
func MakePayloadChan() PayloadChan {
	return make(PayloadChan, 1)
}

//--------------------
// PAYLOAD
//--------------------

// Payload contains key/value pairs of data.
type Payload struct {
	values map[string]interface{}
}

// NewPayload creates a new payload with the given pairs of
// keys and values.
func NewPayload(kvs ...interface{}) *Payload {
	pl := &Payload{
		values: map[string]interface{}{},
	}
	var key string
	for i, kv := range kvs {
		if i%2 == 0 {
			key = fmt.Sprintf("%v", kv)
			pl.values[key] = DefaultValue
			continue
		}
		pl.values[key] = kv
	}
	return pl
}

// Keys returns the keys of the payload.
func (pl *Payload) Keys() []string {
	var keys []string
	for key := range pl.values {
		keys = append(keys, key)
	}
	return keys
}

// String tries to interpred the keyed payload as string.
func (pl *Payload) String(key string) (string, error) {
	v, ok := pl.values[key]
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
func (pl *Payload) Bool(key string) (bool, error) {
	v, ok := pl.values[key]
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
func (pl *Payload) Float64(key string) (float64, error) {
	v, ok := pl.values[key]
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
func (pl *Payload) Int(key string) (int, error) {
	v, ok := pl.values[key]
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
func (pl *Payload) Int64(key string) (int64, error) {
	v, ok := pl.values[key]
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
func (pl *Payload) Uint64(key string) (uint64, error) {
	v, ok := pl.values[key]
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
func (pl *Payload) Time(key string) (time.Time, error) {
	v, ok := pl.values[key]
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
func (pl *Payload) Duration(key string) (time.Duration, error) {
	v, ok := pl.values[key]
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

// Payload tries to interpred the keyed payload as nested payload.
func (pl *Payload) Payload(key string) (*Payload, error) {
	v, ok := pl.values[key]
	if !ok {
		return nil, errors.New(ErrNoValue, msgNoValue, key)
	}
	tv, ok := v.(*Payload)
	if ok {
		return tv, nil
	}
	return nil, errors.New(ErrConverting, msgConverting, key, "event.Payload")
}

// PayloadChan tries to interpred the keyed payload as payload channel.
func (pl *Payload) PayloadChan(key string) (PayloadChan, error) {
	v, ok := pl.values[key]
	if !ok {
		return nil, errors.New(ErrNoValue, msgNoValue, key)
	}
	tv, ok := v.(PayloadChan)
	if ok {
		return tv, nil
	}
	return nil, errors.New(ErrConverting, msgConverting, key, "event.PayloadChan")
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
	var key string
	for i, kv := range kvs {
		if i%2 == 0 {
			key = fmt.Sprintf("%v", kv)
			plc.values[key] = DefaultValue
			continue
		}
		plc.values[key] = kv
	}
	return plc
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
func WithPayload(topic string, pl *Payload) *Event {
	return &Event{
		timestamp: time.Now(),
		topic:     topic,
		payload:   pl,
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

//--------------------
// USEFUL EVENT CREATOR
//--------------------

// NewStatusEvent creates an event containing the returned
// payload channel for status requests to cells.
func NewStatusEvent() (*Event, PayloadChan) {
	plc := MakePayloadChan()
	evt := New(TopicStatus, KeyReply, plc)
	return evt, plc
}

// EOF
