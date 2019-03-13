// Tideland Go Library - Together - Wait
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package wait

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"time"
)

//--------------------
// CONDITION
//--------------------

// Condition has to be implemented for checking the wanted condition. A positive
// condition will return true and nil, a negative false and nil. In case of errors
// during the check false and the error have to be returned. The function will
// be used by the poll functions.
type Condition func() (bool, error)

//--------------------
// POLL
//--------------------

// Poll checks in intervals if the condition is true or the timeout is reached. Also
// in case of a cancelled context the polling will stop.
func Poll(ctx context.Context, interval, timeout time.Duration, condition Condition) error {
	return wait(ctx, MakeLimitedIntervalWaiter(interval, timeout), condition)
	return nil
}

// EOF
