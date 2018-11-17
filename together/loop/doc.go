// Tideland Go Library - Together - Loop
//
// Copyright (C) 2017-2018 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package together/loop of Tideland Go Library supports the developer
// implementing the typical Go idiom for concurrent applications
// running in a loop in the background and doing a select on one
// or more channels. Stopping those loops or getting aware of
// internal errors requires extra efforts. The loop package helps
// to control this kind of goroutines.
//
//     type Printer struct {
//	       printC chan string
//         loop   loop.Loop
//     }
//
//     func (p *printer) worker(l loop.Loop) error {
//         for {
//             select {
//             case <-l.Done():
//                 return nil
//             case str := <-printC:
//                 println(str)
//         }
//     }
//
//     func NewPrinter() *Printer {
//         p := &Printer{
//             printC: make(chan string),
//         }
//         p.loop = loop.New(p.worker).Go()
//         return p
//     }
//
// The worker here now can be stopped with p.loop.Stop() returning
// a possible internal error or p.loop.Terminate(err) and p.loop.Wait().
// Also recovering of internal errors or panics by starting
// the loop with a recoverer function is possible. See the
// code examples.
package loop

// EOF
