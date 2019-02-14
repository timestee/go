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
	"fmt"
	"strconv"
	"strings"

	"tideland.one/go/trace/errors"
)

//--------------------
// VALUE
//--------------------

// Value is simply a byte slice.
type Value []byte

// NewValue creates a value out of the passed data.
func NewValue(value interface{}) Value {
	return Value(valueToBytes(value))
}

// String implements the Stringer interface.
func (v Value) String() string {
	if v == nil {
		return "(nil)"
	}
	return string([]byte(v))
}

// IsOK returns true if the value is the Redis OK value.
func (v Value) IsOK() bool {
	return v.String() == "+OK"
}

// IsNil returns true if the value is the Redis nil value.
func (v Value) IsNil() bool {
	return v == nil
}

// Bool return the value as bool.
func (v Value) Bool() (bool, error) {
	b, err := strconv.ParseBool(v.String())
	if err != nil {
		return false, v.invalidTypeError(err, "bool")
	}
	return b, nil
}

// Int returns the value as int.
func (v Value) Int() (int, error) {
	i, err := strconv.Atoi(v.String())
	if err != nil {
		return 0, v.invalidTypeError(err, "int")
	}
	return i, nil
}

// Int64 returns the value as int64.
func (v Value) Int64() (int64, error) {
	i, err := strconv.ParseInt(v.String(), 10, 64)
	if err != nil {
		return 0, v.invalidTypeError(err, "int64")
	}
	return i, nil
}

// Uint64 returns the value as uint64.
func (v Value) Uint64() (uint64, error) {
	i, err := strconv.ParseUint(v.String(), 10, 64)
	if err != nil {
		return 0, v.invalidTypeError(err, "uint64")
	}
	return i, nil
}

// Float64 returns the value as float64.
func (v Value) Float64() (float64, error) {
	f, err := strconv.ParseFloat(v.String(), 64)
	if err != nil {
		return 0.0, v.invalidTypeError(err, "float64")
	}
	return f, nil
}

// Bytes returns the value as byte slice.
func (v Value) Bytes() []byte {
	return []byte(v)
}

// StringSlice returns the value as slice of strings when seperated by CRLF.
func (v Value) StringSlice() []string {
	return strings.Split(v.String(), "\r\n")
}

// StringMap returns the value as a map of strings when seperated by CRLF
// and colons between key and value.
func (v Value) StringMap() map[string]string {
	tmp := v.StringSlice()
	m := make(map[string]string, len(tmp))
	for _, s := range tmp {
		kv := strings.Split(s, ":")
		if len(kv) > 1 {
			m[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
		}
	}
	return m
}

// Unpack removes the braces of a list value.
func (v Value) Unpack() Value {
	if len(v) > 2 && v[0] == '[' && v[len(v)-1] == ']' {
		return Value(v[1 : len(v)-1])
	}
	return v
}

// invalidTypeError returns an annotated error if a value access has
// been unsuccessful.
func (v Value) invalidTypeError(err error, descr string) error {
	return errors.Annotate(err, ErrInvalidType, msgInvalidType, v.String(), descr)
}

// Values is a set of values.
type Values []Value

// Len returns the number of values.
func (vs Values) Len() int {
	return len(vs)
}

// Strings returns all values as strings.
func (vs Values) Strings() []string {
	ss := make([]string, len(vs))
	for i, v := range vs {
		ss[i] = v.String()
	}
	return ss
}

//--------------------
// KEY/VALUE
//--------------------

// KeyValue combines a pair of key and value
type KeyValue struct {
	Key   string
	Value Value
}

// String implements the Stringer interface.
func (kv *KeyValue) String() string {
	return fmt.Sprintf("%s = %v", kv.Key, kv.Value)
}

// KeyValues is a set of KeyValues.
type KeyValues []KeyValue

// Len returns the number of keys and values in the set.
func (kvs KeyValues) Len() int {
	return len(kvs)
}

// String implements the Stringer interface.
func (kvs KeyValues) String() string {
	kvss := []string{}
	for _, kv := range kvs {
		kvss = append(kvss, kv.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(kvss, " / "))
}

//--------------------
// SCORED VALUE
//--------------------

// ScoredValue helps to add a set member together with its score.
type ScoredValue struct {
	Score float64
	Value Value
}

// String implements the Stringer interface.
func (sv *ScoredValue) String() string {
	return fmt.Sprintf("%v (%f)", sv.Value, sv.Score)
}

// ScoredValues is a set of ScoreValues.
type ScoredValues []ScoredValue

// Len returns the number of scored values in the set.
func (svs ScoredValues) Len() int {
	return len(svs)
}

// String implements the Stringer interface.
func (svs ScoredValues) String() string {
	svss := []string{}
	for _, sv := range svs {
		svss = append(svss, sv.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(svss, " / "))
}

//--------------------
// HASH
//--------------------

// Hash maps multiple fields of a hash to the
// according result values.
type Hash map[string]Value

// CreateHash creates a hash with the passed keys and values.
func CreateHash(kvs map[string]interface{}) Hash {
	h := make(Hash)
	for k, v := range kvs {
		h.Set(k, v)
	}
	return h
}

// Len returns the number of elements in the hash.
func (h Hash) Len() int {
	return len(h)
}

// Set sets a key to the given value.
func (h Hash) Set(key string, value interface{}) Hash {
	h[key] = Value(valueToBytes(value))
	return h
}

// String implements the Stringer interface.
func (h Hash) String(key string) (string, error) {
	if value, ok := h[key]; ok {
		return value.String(), nil
	}
	return "", errors.New(ErrInvalidKey, msgInvalidKey, key)
}

// Bool returns the value of a key as bool.
func (h Hash) Bool(key string) (bool, error) {
	if value, ok := h[key]; ok {
		return value.Bool()
	}
	return false, errors.New(ErrInvalidKey, msgInvalidKey, key)
}

// Int returns the value of a key as int.
func (h Hash) Int(key string) (int, error) {
	if value, ok := h[key]; ok {
		return value.Int()
	}
	return 0, errors.New(ErrInvalidKey, msgInvalidKey, key)
}

// Int64 returns the value of a key as int64.
func (h Hash) Int64(key string) (int64, error) {
	if value, ok := h[key]; ok {
		return value.Int64()
	}
	return 0, errors.New(ErrInvalidKey, msgInvalidKey, key)
}

// Uint64 returns the value of a key as uint64.
func (h Hash) Uint64(key string) (uint64, error) {
	if value, ok := h[key]; ok {
		return value.Uint64()
	}
	return 0, errors.New(ErrInvalidKey, msgInvalidKey, key)
}

// Float64 returns the value of a key as float64.
func (h Hash) Float64(key string) (float64, error) {
	if value, ok := h[key]; ok {
		return value.Float64()
	}
	return 0.0, errors.New(ErrInvalidKey, msgInvalidKey, key)
}

// Bytes returns the value of a key as byte slice.
func (h Hash) Bytes(key string) []byte {
	if value, ok := h[key]; ok {
		return value.Bytes()
	}
	return []byte{}
}

// StringSlice returns the value of a key as string slice.
func (h Hash) StringSlice(key string) []string {
	if value, ok := h[key]; ok {
		return value.StringSlice()
	}
	return []string{}
}

// StringMap returns the value of a key as string map.
func (h Hash) StringMap(key string) map[string]string {
	if value, ok := h[key]; ok {
		return value.StringMap()
	}
	return map[string]string{}
}

// Hashable represents types for Redis hashes.
type Hashable interface {
	// Len returns the number of elements of the hashable.
	Len() int

	// GetHash returns the hashable as hash.
	GetHash() Hash

	// SetHash sets the hashable state by the values
	// of a hash.
	SetHash(h Hash)
}

//--------------------
// PUBLISHED VALUE
//--------------------

// PublishedValue contains a published value and its channel
// channel pattern.
type PublishedValue interface {
	// Kind tells if it is a message or a subscription
	// value.
	Kind() string

	// Channel returns the channel where the value has
	// been received.
	Channel() string

	// Count returns the number of currently subscribed
	// channels when receiving this value.
	Count() int

	// Value returns the published value.
	Value() Value
}

// publishedValue implements the PublishedValue interface.
type publishedValue struct {
	kind    string
	channel string
	count   int
	value   Value
}

// Kind implements the PublishedValue interface.
func (pv *publishedValue) Kind() string {
	return pv.kind
}

// Channel implements the PublishedValue interface.
func (pv *publishedValue) Channel() string {
	return pv.channel
}

// Count implements the PublishedValue interface.
func (pv *publishedValue) Count() int {
	return pv.count
}

// Value implements the PublishedValue interface.
func (pv *publishedValue) Value() Value {
	return pv.value
}

// EOF
