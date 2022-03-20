// Copyright 2022 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package helper // import "github.com/wabarc/helper"

import (
	"context"
	"time"

	"github.com/fortytw2/leaktest"
)

// CheckTest calls `leaktest.Check` to snapshots the currently-running goroutines
// to be run at the end of tests to see whether any goroutines leaked, waiting up
// to 5 seconds in error conditions.
func CheckTest(t leaktest.ErrorReporter) {
	leaktest.Check(t)()
}

// CheckTimeout calls `leaktest.CheckTimeout` which is same as Check, but with a configurable timeout
func CheckTimeout(t leaktest.ErrorReporter, d time.Duration) {
	leaktest.CheckTimeout(t, d)()
}

// CheckContext calls `leaktest.CheckContext` which is same as CheckTest, but uses a
// context.Context for cancellation and timeout control.
func CheckContext(ctx context.Context, t leaktest.ErrorReporter) {
	leaktest.CheckContext(ctx, t)()
}
