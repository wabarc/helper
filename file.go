// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

/*
Package helper handles common functions for the waybackk application in Golang.
*/

package helper // import "github.com/wabarc/helper"

import (
	"fmt"
	"mime"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// FileName returns filename from webpage's link and content type.
func FileName(link, contentType string) string {
	now := time.Now().Format("2006-01-02-150405")
	ext := "html"
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
