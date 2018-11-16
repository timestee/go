// Tideland Go Library - Audit - Capture
//
// Copyright (C) 2017-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package capture of Tideland Go Library Audit assists testing code 
// writing output to stdout or stderr. Those will be temporarily 
// exchanged so that the written output will be catched and can be
// retrieved.
//
//     cout := capture.Stdout(func() {
//         fmt.Printf("Hello, World!")
//     })
//     cerr := capture.Stderr(func() { ... })
//
//     assert.Equal(cout.String(), "Hello, World!")
//
//     cout, cerr = capture.Both(func() { ... })
//
// The captured content data also can be retrieved as bytes.
package capture

// EOF
