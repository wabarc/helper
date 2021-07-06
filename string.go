// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package helper // import "github.com/wabarc/helper"

import (
	"bufio"
	"io"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandString(length int, letter string) string {
	alphabet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	switch letter {
	case "capital", "upper":
		alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "lower":
		alphabet = "abcdefghijklmnopqrstuvwxyz"
	}

	bytes := make([]byte, length)
	rand.Seed(time.Now().UnixNano())
	rand.Read(bytes)
	for i := range bytes {
		bytes[i] = alphabet[rand.Int63()%int64(len(alphabet))]
	}

	return string(bytes)
}

func UTF8Encoding(s string) (r io.Reader, err error) {
	buf := strings.NewReader(s)
	e, name, err := determineEncodingFromReader(buf)
	if err == io.EOF {
		return buf, nil
	}
	if err != nil {
		return
	}
	rd, err := charset.NewReader(buf, name)
	if err != nil {
		return
	}
	r = transform.NewReader(rd, e.NewDecoder())
	return
}

func determineEncodingFromReader(r io.Reader) (e encoding.Encoding, name string, err error) {
	buf, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		return
	}

	e, name, _ = charset.DetermineEncoding(buf, "")
	return
}
