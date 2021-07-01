// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package helper // import "github.com/wabarc/helper"

import (
	"math/rand"
	"time"
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
