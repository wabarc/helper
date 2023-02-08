// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package helper // import "github.com/wabarc/helper"

import (
	"bufio"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
)

const space = ` `

// Unsetenv unsets given envs
func Unsetenv(envs ...string) {
	for _, env := range envs {
		os.Unsetenv(env)
	}
}

// HasStdin determines if the user has piped input
func HasStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	mode := stat.Mode()

	isPipedFromChrDev := (mode & os.ModeCharDevice) == 0
	isPipedFromFIFO := (mode & os.ModeNamedPipe) != 0
	isTerminal := isatty.IsTerminal(os.Stdout.Fd()) ||
		isatty.IsCygwinTerminal(os.Stdout.Fd())

	return isTerminal && (isPipedFromChrDev || isPipedFromFIFO)
}

// ReadStdin reads stdin to slice.
func ReadStdin() (stdin []string) {
	if HasStdin() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			stdin = append(stdin, strings.Split(scanner.Text(), space)...)
		}
	}
	return
}
