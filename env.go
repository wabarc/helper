// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package helper // import "github.com/wabarc/helper"

import (
	"os"
)

// Unsetenv unsets given envs
func Unsetenv(envs ...string) {
	for _, env := range envs {
		os.Unsetenv(env)
	}
}
