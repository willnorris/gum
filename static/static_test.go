// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package static

import (
	"bytes"
	"testing"
)

func TestParseFile(t *testing.T) {
	input := `<link rel="shortlink" href="s"><link rel="canonical" href="p">`
	buf := bytes.NewBufferString(input)

	shortlink, permalink, err := parseFile(buf)
	if err != nil {
		t.Fatalf("error parsing file: %v", err)
	}
	if got, want := shortlink, "s"; got != want {
		t.Fatalf("parseFile(%q) returned shortlink %v, want %v", input, got, want)
	}
	if got, want := permalink, "p"; got != want {
		t.Fatalf("parseFile(%q) returned permalink %v, want %v", input, got, want)
	}
}
