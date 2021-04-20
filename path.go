// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package helper // import "github.com/wabarc/helper"

import (
	"os"
	"os/exec"
	"path/filepath"
)

// FindExecPath tries to find the Chrome browser somewhere in the current
// system. It performs a rather agressive search, which is the same in all
// systems. That may make it a bit slow, but it will only be run when creating a
// new ExecAllocator.
// Fork from: https://github.com/chromedp/chromedp/blob/5b511f121cb2fc1e868b8667049d3cff32a6e0b1/allocate.go#L331
func FindChromeExecPath() string {
	for _, path := range [...]string{
		// Unix-like
		"headless_shell",
		"headless-shell",
		"chromium",
		"chromium-browser",
		"google-chrome",
		"google-chrome-stable",
		"google-chrome-beta",
		"google-chrome-unstable",
		"/usr/bin/google-chrome",

		// Windows
		"chrome",
		"chrome.exe", // in case PATHEXT is misconfigured
		`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
		`C:\Program Files\Google\Chrome\Application\chrome.exe`,
		filepath.Join(os.Getenv("USERPROFILE"), `AppData\Local\Google\Chrome\Application\chrome.exe`),

		// Mac
		"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome",
	} {
		found, err := exec.LookPath(path)
		if err == nil {
			return found
		}
	}
	// Fall back to something simple and sensible, to give a useful error
	// message.
	return "google-chrome"
}
