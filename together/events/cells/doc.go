// Tideland Go Library - Together - Events - Cells
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package cells is the core package of the Tideland event processing. It
// provides type for networked cells running individual behaviors.
//
// These cell behaviors are defined based on an interface and can be added
// to the runtime. Here they are running as concurrent cells that
// can be networked and communicate via events. Several useful behaviors
// are already provided with the behaviors package.
//
// New runtimes are created with
//
//     runtime := cells.NewRuntime(<identifier>)
//
// and cells are added with
//
//    runtime.StartCell("foo", NewFooBehavior())
//
// Cells then can be subscribed with
//
//    runtime.Subscribe("foo", "bar")
//
// so that events which are emitted by the "foo" cell during the processing
// of received or created events will be received by the "bar" cell. Each cell
// can have multiple cells subscibed.
//
// Events from the outside are emitted using
//
//     runtime.Emit("foo", <my-event>)
//
// Behaviors have to implement the cells.Behavior interface.
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
package cells // import "tideland.dev/go/together/events/cells"

// EOF
