// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

/*
Package helper handles common functions for the waybackk application in Golang.
*/
package helper // import "github.com/wabarc/helper"

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strconv"
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
func RealURI(u *url.URL) (*url.URL, error) {
	return realURI(u, safeDialContext)
}

func realURI(u *url.URL, dialContext func(context.Context, string, string) (net.Conn, error)) (*url.URL, error) {
	if err := validateURL(u); err != nil {
		return nil, err
	}

	transport := &http.Transport{DialContext: dialContext}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return validateURL(req.URL)
		},
	}

	req, err := http.NewRequest(http.MethodHead, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("invalid request URI: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.Request == nil || resp.Request.URL == nil {
		return nil, fmt.Errorf("response has no final URL")
	}
	return resp.Request.URL, nil
}

func safeDialContext(
	ctx context.Context,
	network, addr string,
) (net.Conn, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid host:port: %w", err)
	}

	ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil || len(ips) == 0 {
		return nil, fmt.Errorf("could not resolve host")
	}

	// Reject the hostname if any resolved address is non-public.
	for _, ip := range ips {
		if isBlockedIP(ip.IP) {
			return nil, fmt.Errorf("request to internal or non-public address is blocked")
		}
	}

	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return nil, fmt.Errorf("invalid port")
	}

	// Connect directly to a checked IP, avoiding a second DNS lookup.
	dialer := &net.Dialer{Timeout: 5 * time.Second}
	return dialer.DialContext(
		ctx,
		network,
		net.JoinHostPort(ips[0].IP.String(), port),
	)
}

func validateURL(u *url.URL) error {
	if u == nil || (u.Scheme != "http" && u.Scheme != "https") {
		return fmt.Errorf("only http and https URLs are allowed")
	}
	if u.Hostname() == "" || u.User != nil {
		return fmt.Errorf("invalid URL")
	}
	return nil
}

func isBlockedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}
	if v4 := ip.To4(); v4 != nil {
		ip = v4
	}

	return ip == nil ||
		ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsUnspecified() ||
		ip.IsMulticast()
}

func ipString(ip net.IP) string {
	if ip == nil {
		return ""
	}
	return ip.String()
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
