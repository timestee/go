// Tideland Go Library - Audit - Generators
//
// Copyright (C) 2013-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package audit/generators of Tideland Go Library helps to quickly
// generate data needed for unit tests. The generation of all supported
// different types is based on a passed rand.Rand. When using the same
// value here the generated data will be the same when repeating
// tests. So generators.FixedRand() delivers such a fixed value.
package generators

// EOF