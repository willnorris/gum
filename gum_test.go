// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://opensource.org/licenses/BSD-3-Clause

package gum

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMappings(t *testing.T) {
	g := NewServer()
	srv := httptest.NewServer(g)
	tr := http.DefaultTransport
	defer srv.Close()

	req, err := http.NewRequest("GET", srv.URL+"/foo", nil)
	if err != nil {
		t.Errorf("error construction request: %v", err)
	}

	mappings := []*Mapping{
		nil,
		{Permalink: "/bar"},
		{Permalink: "/baz"},
		{Permalink: ""},
	}

	for i, m := range mappings {
		if m != nil {
			m.ShortPath = "/foo"
			g.mappings <- *m
			// sleep long enough for mapping to be processed
			time.Sleep(3 * time.Millisecond)
		}

		resp, err := tr.RoundTrip(req)
		if err != nil {
			t.Errorf("error fetching /foo: %v", err)
		}

		wantCode := http.StatusNotFound
		if m != nil && m.Permalink != "" {
			wantCode = http.StatusMovedPermanently
		}
		if got, want := resp.StatusCode, wantCode; got != want {
			t.Errorf("%d. GET /foo returned status %v, want %v", i, got, want)
		}

		wantLocation := ""
		if m != nil {
			wantLocation = m.Permalink
		}
		if got, want := resp.Header.Get("Location"), wantLocation; got != want {
			t.Errorf("%d. GET /foo returned Location header %q, want %q", i, got, want)
		}

	}
}
