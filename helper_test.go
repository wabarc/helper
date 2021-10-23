// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package helper // import "github.com/wabarc/helper"

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestStrip(t *testing.T) {
	link := "https://example.com/?utm_source=wabarc&utm_medium=cpc"
	if strings.Contains(strip(link), "utm") {
		t.Fail()
	}

	link = "https://example.com/t-55534999?at_custom1=link&at_campaign=64&at_custom3=Regional+East&at_custom2=twitter&at_medium=custom7&at_custom=691F31DA-4E9E-11EB-A68F-435816F31EAE"
	if strings.Contains(strip(link), "at_custom") {
		t.Fail()
	}

	link = "https://weibointl.api.weibo.cn/share/123456.html?weibo_id=101341001431"
	if !strings.EqualFold(strip(link), "https://weibointl.api.weibo.cn/share/123456.html") {
		t.Fail()
	}
}

func TestMatchURL(t *testing.T) {
	var tests = []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "Match Host",
			text:     "foo bar https://example.org/ zoo",
			expected: "https://example.org/",
		},
		{
			name:     "Match Host and Args",
			text:     "foo bar https://example.org/a_(b)?args=世界 zoo",
			expected: "https://example.org/a_(b)?args=%E4%B8%96%E7%95%8C",
		},
		{
			name:     "Match Path",
			text:     "foo bar https://example.org/せかい zoo",
			expected: "https://example.org/%E3%81%9B%E3%81%8B%E3%81%84",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matched := MatchURL(test.text)
			if len(matched) == 0 {
				t.Fatalf("Unexpected match URL number, got %d instead of 0", len(matched))
			}
			if matched[0] != test.expected {
				t.Errorf("Unexpected match URL, got %s instead of [%s]", matched, test.expected)
			}
		})
	}
}

func TestMatchURLFallback(t *testing.T) {
	var tests = []struct {
		name     string
		text     string
		expected string
	}{
		{
			name:     "Match Path",
			text:     "foo bar https://example.org/せかい zoo",
			expected: "https://webcache.googleusercontent.com/search?q=cache:https://example.org/%E3%81%9B%E3%81%8B%E3%81%84",
		},
		{
			name:     "Match and Use Google Cache",
			text:     "foo bar https://example.org/404 zoo",
			expected: "https://webcache.googleusercontent.com/search?q=cache:https://example.org/404",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			matched := MatchURLFallback(test.text)
			if len(matched) == 0 {
				t.Fatalf("Unexpected match URL number, got %d instead of 0", len(matched))
			}
			if matched[0] != test.expected {
				t.Errorf("Unexpected match URL, got %s instead of [%s]", matched, test.expected)
			}
		})
	}
}

func TestIsURL(t *testing.T) {
	allow := []string{
		"http://example.org",
		"https://example.org:443",
	}
	deny := []string{
		"",
		"https",
		"https://",
		"http://www",
		"/testing-path",
		"testing-path",
		"alskjff#?asf//dfas",
	}
	for _, u := range allow {
		if !IsURL(u) {
			t.Fail()
			t.Log(u)
		}
	}
	for _, u := range deny {
		if IsURL(u) {
			t.Fail()
			t.Log(u)
		}
	}
}

func TestFileName(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		link string
		ct   string

		suffix string
	}{
		{
			link:   "",
			ct:     "",
			suffix: "",
		},
		{
			link:   "https://example.org",
			ct:     "text/html; charset=UTF-8",
			suffix: "-example-org.htm",
		},
		{
			link:   "https://example.org/some-path?k=v",
			ct:     "text/html; charset=UTF-8",
			suffix: "-example-org-some-path.htm",
		},
		{
			link:   "https://example.org/path-to-image",
			ct:     "image/png",
			suffix: "-example-org-path-to-image.png",
		},
		{
			link:   "https://example.org/path-to-image",
			ct:     "image/jpeg",
			suffix: "-example-org-path-to-image.jpe",
		},
	}

	for _, test := range tests {
		t.Run(test.suffix, func(t *testing.T) {
			filename := FileName(test.link, test.ct)
			if !strings.HasSuffix(filename, test.suffix) {
				t.Errorf(`Unexpected generate file name, got %s instead of has suffix %s`, filename, test.suffix)
			}
		})
	}
}

func TestFileSize(t *testing.T) {
	tmpfile, err := ioutil.TempFile(".", "helper-testing")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	size := int64(10 * 1024)
	fd, err := os.Create(tmpfile.Name())
	if err != nil {
		t.Fatal("Failed to create output")
	}
	_, err = fd.Seek(size-1, 0)
	if err != nil {
		t.Fatal("Failed to seek")
	}
	_, err = fd.Write([]byte{0})
	if err != nil {
		t.Fatal("Write failed")
	}
	err = fd.Close()
	if err != nil {
		t.Fatal("Failed to close file")
	}

	got := FileSize(tmpfile.Name())
	if got != size {
		t.Fail()
	}
}

func TestRealURI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test in short mode.")
	}

	final := "https://example.com/"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, final, http.StatusSeeOther)
	}))
	defer ts.Close()

	got := RealURI(ts.URL)
	if got != final {
		t.Fatalf("Test get final URL failed, expect: %v got: %s", final, got)
	}
}

func TestTinyURL(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip test in short mode.")
	}

	link := "https://example.com/"
	got := TinyURL(link)
	if !strings.Contains(got, "tinyurl.com") {
		t.Fatalf("Tiny URL failed, got: %s", got)
	}
}

func TestRandString(t *testing.T) {
	got := RandString(36, "")
	if len(got) != 36 {
		t.Log(got)
		t.Fatalf("Test random string failed, expect: %d, got: %d", 36, len(got))
	}
}

func TestMockServer(t *testing.T) {
	httpClient, mux, server := MockServer()
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World.")
	})

	resp, err := httpClient.Get(server.URL)
	if err != nil {
		t.Fatalf(`Unexpected http get %s failed: %v`, server.URL, err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf(`Unexpected read body failed: %v`, err)
	}
	if string(bytes) != "Hello, World." {
		t.Error("Parsed content not match.")
	}
}

func TestNotFound(t *testing.T) {
	t.Parallel()

	_, mux, server := MockServer()
	defer server.Close()

	var tests = []struct {
		name     string
		code     int
		expected bool
	}{
		{
			name:     "HTTP 200",
			code:     http.StatusOK,
			expected: false,
		},
		{
			name:     "HTTP 404",
			code:     http.StatusOK,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := "/" + RandString(5, "lower")
			mux.HandleFunc(p, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.code)
				fmt.Fprintf(w, "Hello, World.")
			})

			f := NotFound(server.URL + p)
			if f != test.expected {
				t.Fatalf(`Unexpected check url status, got %v instead of %t`, f, test.expected)
			}
		})
	}
}

func TestWritable(t *testing.T) {
	t.Parallel()

	dir, err := ioutil.TempDir(os.TempDir(), "helper-")
	if err != nil {
		t.Fatalf(`Unexpected create temp dir: %v`, err)
	}
	defer os.RemoveAll(dir)

	var tests = []struct {
		name string
		perm os.FileMode
		expt error
	}{
		{
			name: "wrx",
			perm: 0777,
			expt: nil,
		},
		{
			name: "r",
			perm: 0400,
			expt: fmt.Errorf(`'%s' is not writable`, filepath.Join(dir, "r")),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(dir, test.name)
			if err := os.Mkdir(path, test.perm); err != nil {
				t.Fatalf(`Unexpected create sub dir: %v`, err)
			}
			if err := Writable(path); err != nil && err.Error() != test.expt.Error() {
				t.Fatalf(`Unexpected dir writable, got <%v> instead of <%v>`, err, test.expt)
			}
		})
	}
}

func TestSetField(t *testing.T) {
	type s struct {
		Key string
		Val string
	}

	var test s
	if err := SetField(&test, "Key", "foo"); err != nil {
		t.Fatalf(`Unexpected set field: %v`, err)
	}
	if test.Key != "foo" {
		t.Fail()
	}
}

func TestIsDir(t *testing.T) {
	content := []byte("Hello, Golang!")
	tmpfile, err := ioutil.TempFile("", "helper-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}

	if ok := IsDir(tmpfile.Name()); ok {
		t.Fatalf(`Unexpected check path is directory, got %t instread of false`, ok)
	}

	dir, err := ioutil.TempDir("", "helper")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	if ok := IsDir(dir); !ok {
		t.Fatalf(`Unexpected check path is directory, got %t instread of true`, ok)
	}
}

func TestExists(t *testing.T) {
	t.Parallel()

	content := []byte("Hello, Golang!")
	tmpfile, err := ioutil.TempFile("", "helper-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		name     string
		filepath string
		expected bool
	}{
		{
			name:     "file exist",
			filepath: tmpfile.Name(),
			expected: true,
		},
		{
			name:     "file not exist",
			filepath: RandString(5, ""),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if ok := Exists(test.filepath); ok != test.expected {
				t.Fatalf(`Unexpected check file exists, got %t instread of %t`, ok, test.expected)
			}
		})
	}
}

func TestMoveFile(t *testing.T) {
	dir, err := ioutil.TempDir("", "helper")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	content := []byte("Hello, Golang!")
	srcfile, err := ioutil.TempFile("", "helper-")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := srcfile.Write(content); err != nil {
		t.Fatal(err)
	}

	dstfile := filepath.Join(dir, RandString(10, ""))
	if err := MoveFile(srcfile.Name(), dstfile); err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat(dstfile)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWebPToPNG(t *testing.T) {
	if _, err := exec.LookPath("dwebp"); err != nil {
		t.Skip(err)
	}
	src := "testdata/1.webp"
	dst := "testdata/1.png"

	if err := WebPToPNG(src, dst); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(dst); os.IsNotExist(err) {
		t.Fatal(err)
	}
}

func TestViaTor(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.NewServeMux())
	defer server.Close()

	p, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	var tests = []struct {
		host string
		port string

		addr string
	}{
		{
			host: "",
			port: "",
			addr: "127.0.0.1:9050",
		},
		{
			host: p.Hostname(),
			port: p.Port(),
			addr: p.Host,
		},
	}

	for _, test := range tests {
		t.Run(test.addr, func(t *testing.T) {
			os.Clearenv()
			os.Setenv("TOR_HOST", test.host)
			os.Setenv("TOR_SOCKS_PORT", test.port)
			addr, err := ViaTor()
			if err != nil {
				t.Fatal(err)
			}
			if addr != test.addr {
				t.Errorf(`Unexpected via tor, got %s instead of %s`, addr, test.addr)
			}
		})
	}
}
