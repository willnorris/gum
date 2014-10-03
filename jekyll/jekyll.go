// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

// Package jekyll parses jekyll post files.
package jekyll

import (
	"net/http"
	"net/url"

	"github.com/golang/glog"
	"willnorris.com/go/gum"
)

// Handler handles short URLs for jekyll posts.
type Handler struct {
	// Site is the Jekyll site this handler serves URLs for.
	site *Site
}

// NewHandler constructs a new Handler with the specified base path which
// contains the Jekyll site (that is, the directory containing the Jekyll
// _config.yml file).
func NewHandler(path string) (*Handler, error) {
	site, err := NewSite(path)
	if err != nil {
		return nil, err
	}

	return &Handler{site: site}, nil
}

// Mappings implements gum.Handler.
func (h *Handler) Mappings(mappings chan<- gum.Mapping) {
	glog.Infof("Jekyll handler added for site: %v", h.site.base)

	template := h.site.PermalinkTemplate()
	for _, p := range h.site.Posts {
		permalink := p.Permalink(template)
		if _, err := url.Parse(permalink); err != nil {
			glog.Errorf("Jekyll permalink is not a valid URL: %v", err)
			continue
		}

		shortURLs, err := p.ShortURLs()
		if err != nil {
			glog.Errorf("Error parsing Jekyll short URLs: %v", err)
			continue
		}

		for _, u := range shortURLs {
			if u == "" {
				continue
			}

			glog.Infof("  %v => %v", u, permalink)
			mappings <- gum.Mapping{ShortPath: u, Permalink: permalink}
		}

		// TODO: populate date-based short urls
	}
}

// Register is a noop for this handler.
func (h *Handler) Register(mux *http.ServeMux) {}
