// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package jekyll

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestPage_ParseFrontMatter(t *testing.T) {
	title := "t"
	date := "2014-05-28 13:50:27 -0700"

	input := fmt.Sprintf("---\ntitle: %s\ndate: '%s'\n---\nbody\n", title, date)
	buf := bytes.NewBufferString(input)

	p := new(Page)
	err := p.parseFrontMatter(buf)
	if err != nil {
		t.Fatalf("error parsing page: %v", err)
	}

	if got, want := p.FrontMatter["title"], title; got != want {
		t.Errorf("FrontMatter title got: %v, want: %v", got, want)
	}
	if got, want := p.FrontMatter["date"], date; got != want {
		t.Errorf("FrontMatter title got: %v, want: %v", got, want)
	}
}

func TestPage_Time(t *testing.T) {
	p := new(Page)
	var want time.Time

	if got := p.Time(); !got.IsZero() {
		t.Errorf("p.Time got: %v, want zero value", got)
	}

	// date from front matter
	p = &Page{
		FrontMatter: map[string]interface{}{
			"date": "2014-05-28 13:50:27 -0700",
		},
	}
	want = time.Date(2014, 5, 28, 13, 50, 27, 0, time.FixedZone("PDT", (-7*3600)))
	if got := p.Time(); !got.Equal(want) {
		t.Errorf("p.Time got: %v, want: %v", got, want)
	}

	// date from file name
	p = &Page{Name: "2014-05-28-test.md"}
	want = time.Date(2014, 5, 28, 0, 0, 0, 0, time.UTC)
	if got := p.Time(); !got.Equal(want) {
		t.Errorf("p.Time got: %v, want: %v", got, want)
	}
}
