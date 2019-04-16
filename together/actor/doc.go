// Tideland Go Library - Together - Actor
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package actor supports the simple creation of concurrent applications
// following the idea of actor models. The work to be done has to be defined
// as func() inside your public methods or functions and sent to the actor
// running in the background.
//
//     type Counter struct {
//         counter int
//         act     actor.Actor
//     }
//
//     func NewCounter() *Counter {
//         c := &Counter{}
//         c.act = actor.New(actor.WithFinalizer(func(err error) error {
//             if err != nil {
//                 return err
//             }
//             c.counter = 0
//             return nil
//         }).Go()
//         return c
//     }
//
//     func (c *Counter) Incr(i int) int {
//         var counter int
//         c.act.DoSync(func() error {
//             c.counter += i
//             counter = c.counter
//             return nil
//         })
//         return counter
//     }
//
// Different options for the constructor allow to pass a context for stopping,
// how many actions are queued, and how panics in actions shall be handled.
package actor // import "tideland.dev/go/together/actor"

// EOF
