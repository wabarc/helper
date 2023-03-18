// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package helper // import "github.com/wabarc/helper"

import (
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Writable ensures the directory exists and is writable
func Writable(dir string) error {
	// Construct the dir if missing
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	// Check the directory is writable
	if f, err := os.Create(filepath.Join(dir, "._check_writable")); err == nil {
		f.Close()
		os.Remove(f.Name())
	} else {
		return errors.New("'" + dir + "' is not writable")
	}
	return nil
}

// IsDir ensures directory of given path
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}

// RetryRemoveAll will attempt to remove an item or directory up to the given number of retries.
// NOTE: This function is necessary because of a Windows bug with removing files that have been recently used: https://github.com/golang/go/issues/51442
func RetryRemoveAll(path string, retries int) error {
	for i := 0; i < retries; i++ {
		err := os.RemoveAll(path)
		if err == nil {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}

	return os.RemoveAll(path)
}
