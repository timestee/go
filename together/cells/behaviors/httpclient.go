// Tideland Go Library - Together - Cells - Behaviors
//
// Copyright (C) 2010-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package behaviors // import "tideland.dev/go/together/cells/behaviors"

//--------------------
// IMPORTS
//--------------------

import (
	"net/http"

	"tideland.dev/go/net/httpx"
	"tideland.dev/go/together/cells/event"
	"tideland.dev/go/together/cells/mesh"
)

//--------------------
// HTTP CLIENT BEHAVIOR
//--------------------

// httpClientBehavior performs HTTP requests.
type httpClientBehavior struct {
	id      string
	emitter mesh.Emitter
}

// NewHTTPClientBehavior performs HTTP request and transforms the
// response into emitted payload, depending on the content-type.
func NewHTTPClientBehavior(id string) mesh.Behavior {
	return &httpClientBehavior{
		id: id,
	}
}

// ID returns the individual identifier of a behavior instance.
func (b *httpClientBehavior) ID() string {
	return b.id
}

// Init the behavior.
func (b *httpClientBehavior) Init(emitter mesh.Emitter) error {
	b.emitter = emitter
	return nil
}

// Terminate the behavior.
func (b *httpClientBehavior) Terminate() error {
	return nil
}

// Process performs the HTTP request.
func (b *httpClientBehavior) Process(evt *event.Event) error {
	switch evt.Topic() {
	case TopicHTTPGet:
		return b.processGet(evt)
	}
	return nil
}

// Recover from an error.
func (b *httpClientBehavior) Recover(err interface{}) error {
	return nil
}

// processGet handles the GET request.
func (b *httpClientBehavior) processGet(evt *event.Event) error {
	id := evt.Payload().At("id").AsString("<none>")
	url := evt.Payload().At("url").AsString("")
	resp, err := http.Get(url)
	if err != nil {
		b.emitter.Broadcast(event.New(
			TopicHTTPGetReply,
			"id", id,
			"url", url,
			"code", resp.StatusCode,
			"error", err,
		))
		return nil
	}
	var data interface{}
	err = httpx.UnmarshalBody(resp.Body, resp.Header, &data)
	if err != nil {
		b.emitter.Broadcast(event.New(
			TopicHTTPGetReply,
			"id", id,
			"url", url,
			"error", err,
		))
		return nil
	}
	b.emitter.Broadcast(event.New(
		TopicHTTPGetReply,
		"id", id,
		"url", url,
		"code", resp.StatusCode,
		"type", resp.Header[httpx.HeaderContentType],
		"data", data,
	))
	return nil
}

// EOF
