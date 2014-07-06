// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

// Package jekyll parses jekyll post files.
package jekyll

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
)

// JekyllHandler handles short URLs for jekyll posts.
type JekyllHandler struct {
	// Prefix is the path component prefix this handler should handle.
	// Prefix should not contain leading or trailing slashes.
	Prefix string

	// Site is the Jekyll site this handler serves URLs for.
	site *Site

	urls map[string]*url.URL
}

// NewHandler constructs a new JekyllHandler with the specified prefix and base
// path which contains the Jekyll site (that is, the directory containing the
// Jekyll _config.yml file).
func NewHandler(prefix, path string) (*JekyllHandler, error) {
	h := &JekyllHandler{
		Prefix: prefix,
		urls:   make(map[string]*url.URL),
	}

	var err error
	h.site, err = NewSite(path)
	if err != nil {
		return nil, err
	}

	err = h.populateURLs()
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *JekyllHandler) populateURLs() error {
	template := h.site.PermalinkTemplate()

	for _, p := range h.site.Posts {
		permalink, err := url.Parse(p.Permalink(template))
		if err != nil {
			return err
		}

		shortURLs, err := p.ShortURLs()
		if err != nil {
			return err
		}

		for _, u := range shortURLs {
			if u == "" {
				continue
			}

			if link, ok := h.urls[u]; ok && link != permalink {
				return fmt.Errorf("short url %q is already registered for permalink %q", u, permalink)
			}
			h.urls[u] = permalink
		}

		// TODO: populate date-based short urls
	}

	return nil
}

func (h *JekyllHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if u, ok := h.urls[r.URL.Path]; ok {
		http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
	}
}

// Register this handler with the provided Router.
func (h *JekyllHandler) Register(router *mux.Router) {
	router.Handle("/"+h.Prefix, h)
	router.PathPrefix("/" + h.Prefix + "/").Handler(h)

	router.PathPrefix("/t/").Handler(h)
	router.PathPrefix("/p/").Handler(h)

	// TODO: handle different possible prefixes
}
