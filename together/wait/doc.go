// Tideland Go Library - Together - Wait
//
// Copyright (C) 2019 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package wait provides different ways to wait for conditions
// and to poll for those with timeouts.
//
//     wait.Poll(ctx, time.Second, time.Minute, func() (bool, error) {
//         _, err := os.Stat(myFile)
//         if err == nil {
//             return true, nil
//         }
//         if os.IsNotExist(err) {
//             return false, nil
//         }
//         return false, err
//     })
package wait

// EOF
