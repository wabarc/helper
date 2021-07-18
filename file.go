// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

/*
Package helper handles common functions for the waybackk application in Golang.
*/

package helper // import "github.com/wabarc/helper"

import (
	"fmt"
	"io"
	"mime"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

// FileName returns filename from webpage's link and content type.
func FileName(link, contentType string) string {
	now := time.Now().Format("2006-01-02-150405.000")
	ext := ".html"
	if exts, _ := mime.ExtensionsByType(contentType); len(exts) > 0 {
		ext = exts[0]
	}

	u, err := url.ParseRequestURI(link)
	if err != nil || u.Scheme == "" || u.Hostname() == "" {
		return now + ext
	}

	domain := strings.ReplaceAll(u.Hostname(), ".", "-")
	if u.Path == "" || u.Path == "/" {
		return fmt.Sprintf("%s-%s%s", now, domain, ext)
	}

	baseName := path.Base(u.Path)
	if parts := strings.Split(baseName, "-"); len(parts) > 4 {
		baseName = strings.Join(parts[:4], "-")
	}

	return fmt.Sprintf("%s-%s-%s%s", now, domain, baseName, ext)
}

// FileSeze returns file attritubes of size about an inode, and
// it's unit alway is bytes.
func FileSize(filepath string) int64 {
	f, err := os.Stat(filepath)
	if err != nil {
		return 0
	}

	return f.Size()
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// MoveFile move file to another directory.
func MoveFile(src, dst string) error {
	if src == dst {
		return nil
	}

	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}

	si, err := in.Stat()
	if err != nil {
		return fmt.Errorf("Stat error: %s", err)
	}
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	perm := si.Mode() & os.ModePerm
	out, err := os.OpenFile(dst, flag, perm)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	in.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}

	err = out.Sync()
	if err != nil {
		return fmt.Errorf("Sync error: %s", err)
	}

	err = os.Remove(src)
	if err != nil {
		return fmt.Errorf("Failed removing original file: %s", err)
	}
	return nil
}

// WebPToPNG convert WebP to PNG
func WebPToPNG(src, dst string) error {
	dwebp, err := exec.LookPath("dwebp")
	if err != nil {
		return err
	}
	args := []string{src, "-o", dst}
	cmd := exec.Command(dwebp, args...)
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
