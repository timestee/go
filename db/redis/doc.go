// Tideland Go Library - Database - Redis Client
//
// Copyright (C) 2017-2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package redis provides a very powerful as well as convenient
// client for the Redis database.
//
// A connection to the database can be established by calling
//
//     conn, err := redis.Open(cfg)
//
// The connection provides
//
//     resultSet, err := conn.Do(myCommand)
//
// to execute commands, optionally with arguments like in
// conn.Do("put", "foo", "bar"). It returns a result set with helpers
// to access the returned values and convert them into Go types. For
// typical returnings there are convenient conn.DoXxx() methods.
//
// All conn.Do() methods work atomically and are able to run all commands
// except subscriptions. Also the execution of scripts is possible that
// way. Additionally the execution of commands can be pipelined. A
// pipeline can be established by calling
//
//     ppl, err := redis.OpenPipeline(cfg)
//
// It provides a ppl.Do() method for the execution of individual commands.
// Their results can be collected with ppl.Collect(), which returns a
// slice of result sets containing the responses of the commands.
//
// Due to the nature of the subscription the client provides also here
// an own type which can be instantiated with
//
//     sub, err := redis.OpenSubscription(etc)
//
// Here channels, in the sense of the Redis Pub/Sub, can be subscribed
// or unsubscribed. Published values can be retrieved with sub.Pop().
//
// All three types, connection, pipeline, and subscription, can be
// closed by [conn|ppl|sub].Close().
package redis

// EOF
