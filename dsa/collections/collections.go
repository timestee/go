// Tideland Go Library - DSA - Collections
//
// Copyright (C) 2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package collections

//--------------------
// EXCHANGE TYPES
//--------------------

// KeyValue wraps a key and a value for the key/value iterator.
type KeyValue struct {
	Keys  string
	Value interface{}
}

// KeyStringValue carries a combination of key and string value.
type KeyStringValue struct {
	Key   string
	Value string
}

// EOF
