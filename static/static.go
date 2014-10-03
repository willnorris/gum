// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

// Package static handles redirects parsed from static HTML files.  Files are
// parsed and searched for rel="shortlink" and rel="canonical" links.  If both
// are found, a redirect is registered for the pair.
package static

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"github.com/golang/glog"
	"willnorris.com/go/gum"
)

const (
	relShortlink = "shortlink"
	relCanonical = "canonical"
)

// Handler handles short URLs parsed from static HTML files.
type Handler struct {
	base string
}

// NewHandler constructs a new Handler with the specified base path of HTML
// files.
func NewHandler(base string) (*Handler, error) {
	if stat, err := os.Stat(base); err != nil {
		return nil, err
	} else if !stat.IsDir() {
		return nil, fmt.Errorf("Specified base path %q is not a directory", base)
	}

	return &Handler{base: base}, nil
}

// Mappings implements gum.Handler.
func (h *Handler) Mappings(mappings chan<- gum.Mapping) {
	loadFiles(h.base, mappings)
}

// Register is a noop for this handler.
func (h *Handler) Register(mux *http.ServeMux) {}

func loadFiles(base string, mappings chan<- gum.Mapping) {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			glog.Errorf("error reading file %q: %v", path, err)
			return nil
		}
		if info.IsDir() || filepath.Ext(path) != ".html" {
			// skip directories and non-HTML files
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			glog.Errorf("error opening file %q: %v", path, err)
			return nil
		}
		defer f.Close()

		shortlink, permalink, err := parseFile(f)
		if err != nil {
			glog.Errorf("error parsing file %q: %v", path, err)
			return nil
		} else if len(shortlink) == 0 || len(permalink) == 0 {
			// no shortlink or permalink
			return nil
		}

		shorturl, err := url.Parse(shortlink)
		if err != nil {
			glog.Errorf("error parsing shortlink %q: %v", shortlink, err)
			return nil
		}

		glog.Infof("  %v => %v", shorturl.Path, permalink)
		mappings <- gum.Mapping{ShortPath: shorturl.Path, Permalink: permalink}
		return nil
	}

	err := filepath.Walk(base, walkFn)
	if err != nil {
		glog.Errorf("Walk(%q) returned error: %v", base, err)
	}
}

// parseFile parses r as HTML and returns the URLs of the first links found
// with the "shortlink" and "canonical" rel values.
func parseFile(r io.Reader) (shortlink string, permalink string, err error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", "", err
	}

	var f func(*html.Node) bool
	f = func(n *html.Node) (done bool) {
		if n.Type == html.ElementNode {
			if n.DataAtom == atom.Link || n.DataAtom == atom.A {
				var href, rel string
				for _, a := range n.Attr {
					if a.Key == atom.Href.String() {
						href = a.Val
					}
					if a.Key == atom.Rel.String() {
						rel = a.Val
					}
				}
				if len(href) > 0 && len(rel) > 0 {
					for _, v := range strings.Split(rel, " ") {
						if v == relShortlink && len(shortlink) == 0 {
							shortlink = href
						}
						if v == relCanonical && len(permalink) == 0 {
							permalink = href
						}
					}
				}
				if len(shortlink) > 0 && len(permalink) > 0 {
					return true
				}
			}
		}
		for c := n.FirstChild; c != nil && !done; c = c.NextSibling {
			done = f(c)
		}
		return done
	}

	f(doc)

	return shortlink, permalink, nil
}
