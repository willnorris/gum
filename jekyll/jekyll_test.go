// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package jekyll

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"willnorris.com/go/gum"
)

func TestHandler_ServeHTTP(t *testing.T) {
	dir, cleanup, err := newTestSite("")
	if err != nil {
		t.Fatalf("error creating test site: %v", err)
	}
	defer cleanup()

	// write post
	os.MkdirAll(path.Join(dir, "_posts"), 0755)
	postPath := path.Join(dir, "_posts", "2014-05-28-test.md")
	if err := ioutil.WriteFile(postPath, []byte("---\nwordpress_id: 100\n---\n"), 0644); err != nil {
		t.Fatalf("error creating test post: %v", err)
	}

	handler, err := NewHandler(dir)
	if err != nil {
		t.Fatalf("error constructing handler: %v", err)
	}

	mappings := make(chan gum.Mapping, 3)
	handler.Mappings(mappings)

	want := gum.Mapping{ShortPath: "/b/1f", Permalink: "/2014/05/28/test.html"}
	select {
	case got := <-mappings:
		if got != want {
			t.Errorf("handler returned mapping %v, want %v", got, want)
		}
	default:
		t.Errorf("handler returned no mappings")
	}
}
