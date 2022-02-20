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
	"time"

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

// MatchURLFallback is extract URL from text, and convert to
// Google cache endpoint if not found, returns []string always.
func MatchURLFallback(text string) []string {
	urls := []string{}
	rx := xurls.Strict()
	matches := rx.FindAllString(text, -1)
	cache := "https://webcache.googleusercontent.com/search?q=cache:"
	for _, el := range matches {
		uri := strip(el)
		if NotFound(uri) {
			uri = cache + uri
		}
		urls = append(urls, uri)
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

// NotFound returns a result of URI status is 404
func NotFound(uri string) bool {
	if _, err := url.Parse(uri); err != nil {
		return true
	}

	req, err := http.NewRequest(http.MethodHead, uri, nil)
	if err != nil {
		return true
	}
	ua := `Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.7113.093 Safari/537.36`
	req.Header.Set("User-Agent", ua)

	noRedirect := func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	client := &http.Client{Timeout: 10 * time.Second, CheckRedirect: noRedirect}

	resp, err := client.Do(req)
	if err != nil {
		return true
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusNotFound
}

func strip(link string) string {
	u, err := url.Parse(link)
	if err != nil {
		return ""
	}

	var p = strings.HasPrefix
	var e = strings.EqualFold
	var maps = map[string]func(string, string) bool{
		"utm_":      p,
		"at_custom": p,
		"at_medium": p,
		"weibo_id":  e,
		"fbclid":    e,
		"chksm":     e,
	}
	queries := u.Query()
	for key := range queries {
		for prefix, v := range maps {
			if v(key, prefix) {
				queries.Del(key)
			}
		}
	}

	u.RawQuery = queries.Encode()

	return u.String()
}

// RealURI returns final URL
func RealURI(u *url.URL) *url.URL {
	resp, err := http.Head(u.String())
	if err != nil {
		return u
	}
	defer resp.Body.Close()

	return resp.Request.URL
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
