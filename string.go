// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package helper // import "github.com/wabarc/helper"

import (
	"bufio"
	"crypto/rand"
	"io"
	"math/big"
	"reflect"
	"strings"
	"unsafe"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

func RandString(length int, letter string) string {
	alphabet := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	switch letter {
	case "capital", "upper":
		alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "lower":
		alphabet = "abcdefghijklmnopqrstuvwxyz"
	}

	bytes := make([]byte, length)
	rand.Read(bytes)
	for i := range bytes {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return ""
		}
		bytes[i] = alphabet[num.Int64()%int64(len(alphabet))]
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

// String2Byte converts string to a byte slice without memory allocation.
func String2Byte(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

// Byte2String converts byte slice to a string without memory allocation.
func Byte2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
