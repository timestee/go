// Tideland Go Library - Together - Wait
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package wait provides different ways to wait for conditions by
// polling. The conditions are checked by user defined functions
// with the signature
//
//     func() (bool, error)
//
// Here the bool return value signals if the condition is fulfilled,
// e.g. a file you're waiting for has been written into the according
// directory.
//
// The different pre-defined tickers support
//
// - simple constant intervals,
// - a maximum number of constant intervals,
// - constant intervals with a deadline,
// - constant intervals with a timeout,
// - changing intervals, and
// - jittering intervals.
//
// The behaviour of changing intervals can be user defined by
// functions with the signature
//
//     func(time.Duration) (time.Duration, bool)
//
// Here the argument is the current interval, return values are the
// wanted interval and if the polling shall continue. For the predefined
// tickers according convenience functions named With...() exist.
//
// Example (waiting for a file to exist):
//
//     wait.Poll(
//         ctx,
//         wait.MakeExpiringIntervalTicker(time.Second, 30*time.Second),
//         func() (bool, error) {
//             _, err := os.Stat(myFile)
//             if err == nil {
//                 return true, nil
//             }
//             if os.IsNotExist(err) {
//                 return false, nil
//             }
//             return false, err
//         },
//     )
//
// From external the polling can be stopped by cancelling the context.
package wait

// EOF
