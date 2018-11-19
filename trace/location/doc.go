// Tideland Go Library - Trace - Location
//
// Copyright (C) 2017-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package location of the Tideland Go Library provides a way to retrieve
// the current location in code. This can be used in logging or
// monitoring. Passing an offset helps hiding calling wrappers.
//
//     pkg, file, fun, line := location.Here(0)
//     here := location.HereID(0)
//
// Internal caching fastens retrieval after first call.
package location

// EOF
