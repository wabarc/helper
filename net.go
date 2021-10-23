// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package helper // import "github.com/wabarc/helper"

import (
	"errors"
	"net"
	"os"
	"time"
)

// ViaTor checks the Tor proxy whether running. Host and port
// listening by Tor defaults to 127.0.0.1 and 9050, they can be
// specific with `TOR_HOST` and `TOR_SOCKS_PORT` environments.
//
// ViaTor returns address used by Tor, and an error if
// Tor proxy missing.
func ViaTor() (addr string, err error) {
	host := os.Getenv("TOR_HOST")
	port := os.Getenv("TOR_SOCKS_PORT")
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "9050"
	}

	addr = net.JoinHostPort(host, port)
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		return addr, err
	}
	if conn != nil {
		conn.Close()
		return addr, nil
	}

	return addr, errors.New("can not access tor proxy")
}
