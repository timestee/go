// Tideland Go Library - Network - REST - Audit
//
// Copyright (C) 2009-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package audit of Tideland Go Library is a little helper package
// for the unit testing of the rest package and the resource handlers.
// Requests can easily be created, marshalling data based on the
// content-type is done automatically. Response also provides assert
// methods for the tests.
//
// So first step is to create a test server and register the handler(s)
// to test. Could best be done with a little helper function, depending
// on own needs, e.g. when the context shall contain more information.
//
//     assert := asserts.NewTesting(t, true)
//     cfgStr := "{etc {basepath /}{default-domain testing}{default-resource index}}"
//     cfg, err := etc.ReadString(cfgStr)
//     assert.Nil(err)
//     mux := core.NewMultiplexer(context.Background(), cfg)
//     ts := audit.StartServer(mux, assert)
//     defer ts.Close()
//     err := mux.Register("my-domain", "my-resource", NewMyHandler())
//     assert.Nil(err)
//
// During the tests you create the requests with
//
//     req := audit.NewRequest("GET", "/my-domain/my-resource/4711")
//     req.AddHeader(audit.HeaderAccept, audit.ApplicationJSON)
//
// The request the is done with
//
//     resp := ts.DoRequest(req)
//     resp.AssertStatusEquals(200)
//     rest.AssertHeaderContains(audit.HeaderContentType, audit.ApplicationJSON)
//     resp.AssertBodyContains(`"ResourceID":"4711"`)
//
// Also data can be marshalled including setting the content type
// and the response can be unmarshalled based on that type.
//
//     req.MarshalBody(assert, audit.ApplicationJSON, myInData)
//     ...
//     var myOutData MyType
//     resp.AssertUnmarshalledBody(&myOutData)
//     assert.Equal(myOutData.MyField, "foo")
//
// There are more helpers for a convenient test, but the fields of
// Request and Response can also be accessed directly.
package audit

// EOF
