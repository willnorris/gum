// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package static

import (
	"bytes"
	"testing"

	"willnorris.com/go/gum"
)

func TestParseFile(t *testing.T) {
	input := `<link rel="shortlink" href="/s"><link rel="canonical" href="/p">`
	buf := bytes.NewBufferString(input)

	mappings, err := parseFile(buf)
	if err != nil {
		t.Fatalf("error parsing file: %v", err)
	}
	if len(mappings) == 0 {
		t.Fatal("parseFile returned 0 mappings")
	}
	want := gum.Mapping{ShortPath: "/s", Permalink: "/p"}
	if got := mappings[0]; got != want {
		t.Fatalf("parseFile(%q) returned mapping %v, want %v", input, got, want)
	}
}
