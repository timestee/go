// Tideland Go Library - Together - Notifier
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package notifier help at the coordination of multiple goroutines. First little
// helper is the Closer for the aggregation of typical closer channels <-chan struct{}
// into one Closer.Done() <-chan struct{}. This way for-select-loops don't need to
// each channel individually.
//
//     ca := make(chan struct{})
//     cb := make(chan struct{})
//     c := notifier.NewCloser(ca, cb)
//
//     go func() {
//          for {
//              select {
//              case <-c.Done():
//                  return
//              case foo := <-someOtherChannel:
//                  ...
//              }
//          }
//     }()
//
//     close(cb)  // Or any other of the passed channels.
//
// Second and third type are Notifier and Bundle. They maintain the goroutine
// statuses Starting, Ready, Working, Stopping, and Stopped. Those are kept in each
// Notifier. So an owner of a Notifier can retrieve the status and wait on according
// channels for notification. A set of notifiers is managed by one Bundle. It is
// responsible to take a new status and publish it to the registered Notifiers.
//
//     // Different goroutines interested in the status of the one following.
//     na := notifier.New()
//     go func() {
//         <-na.Working()
//         // Start something dependent.
//         ...
//     }()
//
//     nb := notifier.New()
//     go func() {
//         <-nb.Stopped()
//         // Start something dependent.
//         ...
//     }()
//
//     // Goroutine where the owner of the notifiers are interested in.
//     b := notifier.NewBundle()
//     b.Add(na, nb)
//
//     b.Notify(notifier.Ready)
//     ...
//     b.Notify(notifier.Working)
//     ...
//     b.Notify(notifier.Stopping)
//     ...
//     b.Notify(notifier.Stopped)
//
package notifier

// EOF
