// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package gum

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirectHandler(t *testing.T) {
	tests := []struct {
		prefix, dest string
		in, out      string
	}{
		// common cases
		{
			prefix: "x", dest: "http://example/",
			in: "/x", out: "http://example/",
		},
		{
			prefix: "x", dest: "http://example/",
			in: "/x/", out: "http://example/",
		},
		{
			prefix: "x", dest: "http://example/",
			in: "/x/y", out: "http://example/y",
		},
		{
			prefix: "x", dest: "http://example/",
			in: "/x/y?a=b", out: "http://example/y?a=b",
		},

		// absolute input URL (rare)
		{
			prefix: "x", dest: "http://example/",
			in: "http://foo/x/y", out: "http://example/y",
		},

		// destination URL with path
		{
			prefix: "x", dest: "http://example/a/",
			in: "/x", out: "http://example/a/",
		},
		{
			prefix: "x", dest: "http://example/a/",
			in: "/x/", out: "http://example/a/",
		},
		{
			prefix: "x", dest: "http://example/a/",
			in: "/x/y", out: "http://example/a/y",
		},

		// destination URL with path (no trailing slash)
		{
			prefix: "x", dest: "http://example/a",
			in: "/x", out: "http://example/a",
		},
		{
			prefix: "x", dest: "http://example/a",
			in: "/x/", out: "http://example/a",
		},
		{
			prefix: "x", dest: "http://example/a",
			in: "/x/y", out: "http://example/y",
		},

		// relative destination URL
		{prefix: "x", dest: "/a/", in: "/x", out: "/a/"},
		{prefix: "x", dest: "/a/", in: "/x/", out: "/a/"},
		{prefix: "x", dest: "/a/", in: "/x/y", out: "/a/y"},

		// no prefix
		{
			prefix: "", dest: "http://example/",
			in: "/x", out: "http://example/x",
		},
		{
			prefix: "", dest: "http://example/",
			in: "/x/", out: "http://example/x/",
		},
		{
			prefix: "", dest: "http://example/",
			in: "/x/y", out: "http://example/x/y",
		},

		// root prefix
		{
			prefix: "/", dest: "http://example/",
			in: "/x", out: "http://example/x",
		},
		{
			prefix: "/", dest: "http://example/",
			in: "/x/", out: "http://example/x/",
		},
		{
			prefix: "/", dest: "http://example/",
			in: "/x/y", out: "http://example/x/y",
		},

		// no destination (redirects to root)
		{prefix: "x", dest: "", in: "/x", out: "/"},
		{prefix: "x", dest: "", in: "/x/", out: "/"},
		{prefix: "x", dest: "", in: "/x/y", out: "/y"},
	}

	for _, tt := range tests {
		handler, err := NewRedirectHandler(tt.prefix, tt.dest)
		if err != nil {
			t.Fatalf("error constructing handler: %v", err)
		}

		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", tt.in, nil)
		if err != nil {
			t.Fatalf("error constructing request: %v", err)
		}

		handler.ServeHTTP(w, r)

		if got, want := w.Code, http.StatusMovedPermanently; got != want {
			t.Errorf("Response status code got %v, want %v", got, want)
		}

		loc := w.Header().Get("Location")
		if loc == "" {
			t.Errorf("No location header set for input: %q", tt.in)
		}
		if got, want := loc, tt.out; got != want {
			t.Errorf("Location header for input %q got: %v, want: %v", tt.in, got, want)
		}
	}
}

// Test that RedirectHandler registers proper prefixes on ServeMux.
func TestRedirectHandler_Register(t *testing.T) {
	tests := []struct {
		prefix       string
		in, location string
		code         int
	}{
		{prefix: "x", in: "/x", location: "/", code: http.StatusMovedPermanently},
		{prefix: "x", in: "/x/", location: "/", code: http.StatusMovedPermanently},
		{prefix: "x", in: "/x/y", location: "/y", code: http.StatusMovedPermanently},
		{prefix: "x", in: "/xy", location: "", code: http.StatusNotFound},
	}

	for _, tt := range tests {
		mux := http.NewServeMux()
		handler, err := NewRedirectHandler(tt.prefix, "")
		if err != nil {
			t.Fatalf("error constructing handler: %v", err)
		}
		handler.Register(mux)

		req, err := http.NewRequest("GET", tt.in, nil)
		if err != nil {
			t.Errorf("error constructing request for %q: %v", tt.in, err)
		}

		resp := httptest.NewRecorder()
		mux.ServeHTTP(resp, req)

		if got, want := resp.Code, tt.code; got != want {
			t.Errorf("response for %q had status %v, want %v", tt.in, got, want)
		}
		if got, want := resp.Header().Get("Location"), tt.location; got != want {
			t.Errorf("response Location header was %v, want %v", got, want)
		}
	}
}
