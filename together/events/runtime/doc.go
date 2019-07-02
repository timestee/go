// Tideland Go Library - Together - Events - Runtime
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package runtime is the core package of the Tideland event processing. It
// provides types for networked engines running individual processors.
//
// These processors are defined based on an interface and can be added
// to the runtime. Here they are running concurrently and can be networked
// and communicate via events. Several useful processors are already provided
// with the processors package.
//
// New runtimes are created with
//
//     runtime := runtime.New()
//
// and processors are added with
//
//    runtime.SpawnProcessors(
//        NewFooProc("a"),
//        NewBarProc("b"),
//        NewBazProc("c"),
//    )
//
// These processors can subscribe each other with
//
//    runtime.Subscribe("a", "b", "c")
//
// so that events which are emitted by the processor "a" will be
// received by the processors "b" and "c". Each processor can subscribe
// to multiple other subscribers and even circular subscriptions are
// no problem. But handle with care.
//
// Events from the outside are emitted using
//
//     runtime.Emit("foo", <my-event>)
//
// Sometimes it's needed to directly communicate with a cell to retrieve
// information. In this case the method
//
//     response, err := runtime.Call("foo", <my-request>, <timeout>)
//
// is to be used. Inside the ProcessEvent() of the addressed cell the
// event can be used to send the response with
//
//    switch event.Topic() {
//    case "myRequest?":
//        event.Respond(<my-data>)
//    case ...:
//        ...
//    }
//
// Instructions without a response are simply done by emitting an event
// to a cell.
package runtime // import "tideland.dev/go/together/events/runtime"

// EOF
