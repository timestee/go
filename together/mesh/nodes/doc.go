// Tideland Go Library - Together - Mesh - Nodes
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package nodes is the core package of the Tideland meshed event processing.
// It provides types for meshed nodes running individual behaviors.
//
// These behaviors are defined based on an interface and can be added
// to the mesh. Here they are running concurrently and can be networked
// and communicate via events. Several useful behaviors are already provided
// with the behaviors package.
//
// New meshes are created with
//
//     mesh := nodes.NewMesh()
//
// and nodes are started with
//
//    mesh.SpawnNodes(
//        NewFooer("a"),
//        NewBarer("b"),
//        NewBazer("c"),
//    )
//
// These nodes can subscribe each other with
//
//    mesh.Subscribe("a", "b", "c")
//
// so that events which are emitted by the node "a" will be
// received by the nodes "b" and "c". Each node can subscribe
// to multiple other subscribers and even circular subscriptions are
// no problem. But handle with care.
//
// Events from the outside are emitted using
//
//     mesh.Emit("foo", <my-event>)
//
// Inside of Process() of the behaviors these can analyse
// Event.Topic to decide what to do.
package nodes // import "tideland.dev/go/together/mesh/nodes"

// EOF
