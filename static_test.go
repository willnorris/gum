// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://opensource.org/licenses/BSD-3-Clause

package gum

import (
	"bytes"
	"testing"
)

func TestParseFile(t *testing.T) {
	input := `<link rel="shortlink" href="/s1"><link rel="canonical" href="/p"><a rel="shortlink" href="/s2">`
	buf := bytes.NewBufferString(input)

	mappings, err := parseFile(buf)
	if err != nil {
		t.Fatalf("error parsing file: %v", err)
	}
	if got, want := len(mappings), 2; got != want {
		t.Fatalf("parseFile returned %d mappings, want %d", got, want)
	}

	want := Mapping{ShortPath: "/s1", Permalink: "/p"}
	if got := mappings[0]; got != want {
		t.Fatalf("parseFile(%q) returned mapping %v, want %v", input, got, want)
	}

	want = Mapping{ShortPath: "/s2", Permalink: "/p"}
	if got := mappings[1]; got != want {
		t.Fatalf("parseFile(%q) returned mapping %v, want %v", input, got, want)
	}
}
