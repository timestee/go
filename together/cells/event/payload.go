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
	"strconv"
	"time"

	"tideland.dev/go/trace/failure"
)

//--------------------
// CONSTANTS
//--------------------

// Standard value.
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

// For time converts.
const (
	maxInt = int64(^uint(0) >> 1)
)

//--------------------
// VALUE
//--------------------

// Value contains one payload value.
type Value struct {
	raw interface{}
	err error
}

// IsDefined returns true if this value is defined.
func (v *Value) IsDefined() bool {
	return v.raw != nil
}

// IsUndefined returns true if this value is undefined.
func (v *Value) IsUndefined() bool {
	return v.raw == nil
}

// AsString returns the value as string, dv is taken as default value.
func (v *Value) AsString(dv string) string {
	if v.IsUndefined() {
		// return dv
		return "<undefined>"
	}
	switch tv := v.raw.(type) {
	case string:
		return tv
	case int:
		return strconv.Itoa(tv)
	case float64:
		return strconv.FormatFloat(tv, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(tv)
	case time.Time:
		return tv.Format(time.RFC3339Nano)
	case time.Duration:
		return tv.String()
	}
	return dv
}

// AsInt returns the value as int, dv is taken as default value.
func (v *Value) AsInt(dv int) int {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.raw.(type) {
	case string:
		i, err := strconv.Atoi(tv)
		if err != nil {
			return dv
		}
		return i
	case int:
		return tv
	case float64:
		return int(tv)
	case bool:
		if tv {
			return 1
		}
		return 0
	case time.Time:
		ns := tv.UnixNano()
		if ns > maxInt {
			return dv
		}
		return int(ns)
	case time.Duration:
		ns := tv.Nanoseconds()
		if ns > maxInt {
			return dv
		}
		return int(ns)
	}
	return dv
}

// AsFloat64 returns the value as float64, dv is taken as default value.
func (v *Value) AsFloat64(dv float64) float64 {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.raw.(type) {
	case string:
		f, err := strconv.ParseFloat(tv, 64)
		if err != nil {
			return dv
		}
		return f
	case int:
		return float64(tv)
	case float64:
		return tv
	case bool:
		if tv {
			return 1.0
		}
		return 0.0
	case time.Time:
		ns := tv.UnixNano()
		return float64(ns)
	case time.Duration:
		ns := tv.Nanoseconds()
		return float64(ns)
	}
	return dv
}

// AsBool returns the value as bool, dv is taken as default value.
func (v *Value) AsBool(dv bool) bool {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.raw.(type) {
	case string:
		b, err := strconv.ParseBool(tv)
		if err != nil {
			return dv
		}
		return b
	case int:
		return tv == 1
	case float64:
		return tv == 1.0
	case bool:
		return tv
	case time.Time:
		return tv.UnixNano() > 0
	case time.Duration:
		return tv.Nanoseconds() > 0
	}
	return dv
}

// AsTime returns the value as time, dv is taken as default value.
func (v *Value) AsTime(dv time.Time) time.Time {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.raw.(type) {
	case string:
		for _, timeFormat := range timeFormats {
			t, err := time.Parse(timeFormat, tv)
			if err == nil {
				return t
			}
		}
		return dv
	case int:
		return time.Time{}.Add(time.Duration(tv))
	case float64:
		d := int64(tv)
		return time.Time{}.Add(time.Duration(d))
	case bool:
		if tv {
			return time.Time{}.Add(1)
		}
		return time.Time{}
	case time.Time:
		return tv
	case time.Duration:
		return time.Time{}.Add(tv)
	}
	return dv
}

// AsDuration returns the value as duration, dv is taken as default value.
func (v *Value) AsDuration(dv time.Duration) time.Duration {
	if v.IsUndefined() {
		return dv
	}
	switch tv := v.raw.(type) {
	case string:
		d, err := time.ParseDuration(tv)
		if err == nil {
			return d
		}
		return dv
	case int:
		return time.Duration(tv)
	case float64:
		d := int64(tv)
		return time.Duration(d)
	case bool:
		if tv {
			return 1
		}
		return 0
	case time.Time:
		return time.Duration(tv.UnixNano())
	case time.Duration:
		return tv
	}
	return dv
}

// AsPayloadChan returns the value as payload channel.
func (v *Value) AsPayloadChan() PayloadChan {
	if v.IsUndefined() {
		return nil
	}
	tv, ok := v.raw.(PayloadChan)
	if !ok {
		return nil
	}
	return tv
}

// String implements fmt.Stringer.
func (v *Value) String() string {
	return fmt.Sprintf("%v", v.raw)
}

// Error implements error.
func (v *Value) Error() string {
	return fmt.Sprintf("%v", v.err)
}

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
			key = fmt.Sprintf("%v", kv)
			pl.values[key] = DefaultValue
			continue
		}
		plv, ok := kv.(*Payload)
		if ok {
			// It's a payload value. Add all with joined keys.
			pl.nestMap(key, plv.values)
			continue
		}
		if reflect.TypeOf(kv).Kind() == reflect.Map {
			// It's a map value. Add all with joined keys.
			pl.nestMap(key, kv)
			continue
		}
		// It's a standard non-payload value.
		pl.values[key] = kv
	}
}

// nestMap recursively nests a map with joined keys.
func (pl *Payload) nestMap(key string, value interface{}) {
	iter := reflect.ValueOf(value).MapRange()
	for iter.Next() {
		rvKey := iter.Key()
		rvKeyStr := fmt.Sprintf("%v", rvKey.Interface())
		rvValue := iter.Value()
		rvKind := reflect.TypeOf(rvValue.Interface()).Kind()
		if rvKind == reflect.Map {
			pl.nestMap(key+"/"+rvKeyStr, rvValue.Interface())
			return
		}
		pl.values[key+"/"+rvKeyStr] = rvValue.Interface()
	}
}

// EOF
