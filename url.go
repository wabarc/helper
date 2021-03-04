// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

/*
Package helper handles common functions for the waybackk application in Golang.
*/

package helper // import "github.com/wabarc/helper"

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"mvdan.cc/xurls/v2"
)

// MatchURL is extract URL from text, returns []string always.
func MatchURL(text string) []string {
	urls := []string{}
	rx := xurls.Strict()
	matches := rx.FindAllString(text, -1)
	for _, el := range matches {
		urls = append(urls, strip(el))
	}

	return urls
}

// IsURL returns a result of validation for string.
func IsURL(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	return u.Scheme != "" && strings.Contains(u.Host, ".")
}

func strip(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}

	queries := u.Query()
	for key := range queries {
		if strings.HasPrefix(key, "utm_") || strings.HasPrefix(key, "at_custom") || strings.HasPrefix(key, "at_medium") || strings.EqualFold(key, "weibo_id") {
			queries.Del(key)
		}
	}

	u.RawQuery = queries.Encode()

	return u.String()
}

// RealURI returns final URL
func RealURI(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	return resp.Request.URL.String()
}

func TinyURL(link string) string {
	_, err := url.Parse(link)
	if err != nil {
		return ""
	}

	resp, err := http.Get("https://tinyurl.com/api-create.php?url=" + link)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	final := string(body)
	if final != "Error" {
		return final
	}

	return ""
}
