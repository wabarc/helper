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
	"os"
	"strings"
	"testing"
	"time"
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
	text := `
foo bar https://example.org/ zoo
foo bar https://example.org/a_(b)?args=中文 zoo
foo bar https://example.org/せかい zoo
`
	matched := MatchURL(text)
	if len(matched) != 3 {
		t.Log(matched)
		t.Fail()
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

func TestFileNameWithoutPath(t *testing.T) {
	now := time.Now().Format("2006-01-02-150405")
	expect := now + "-example-org.htm"
	link := "https://example.org"
	ct := "text/html; charset=UTF-8"

	got := FileName(link, ct)
	if got != expect {
		t.Fail()
	}
}

func TestFileNameWithPath(t *testing.T) {
	now := time.Now().Format("2006-01-02-150405")
	expect := now + "-example-org-some-path.htm"
	link := "https://example.org/some-path?k=v"
	ct := "text/html; charset=UTF-8"

	got := FileName(link, ct)
	if got != expect {
		t.Fail()
	}
}

func TestFileNameIsPNG(t *testing.T) {
	now := time.Now().Format("2006-01-02-150405")
	expect := now + "-example-org-path-to-image.png"
	link := "https://example.org/path-to-image"
	ct := "image/png"

	got := FileName(link, ct)
	if got != expect {
		t.Fail()
	}
}

func TestFileNameIsJPG(t *testing.T) {
	now := time.Now().Format("2006-01-02-150405")
	expect := now + "-example-org-path-to-image.jpe"
	link := "https://example.org/path-to-image"
	ct := "image/jpeg"

	got := FileName(link, ct)
	if got != expect {
		t.Fail()
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
