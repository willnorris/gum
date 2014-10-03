// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

// Package gum provides the gum personal short URL resolver.
package gum

import (
	"net/http"

	"github.com/golang/glog"
)

// Server is a short URL redirection server.
type Server struct {
	// ServeMux which handles all incoming requests
	mux *http.ServeMux

	// map of short URL paths to destinations
	urls map[string]string
}

// NewServer constructs a new Server.
func NewServer() *Server {
	s := &Server{
		mux:  http.NewServeMux(),
		urls: make(map[string]string),
	}

	// default handler
	s.mux.HandleFunc("/", s.redirect)

	return s
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// redirect the request if a matching URL mapping has been configured.  If no
// mapping is found, a 404 status is returned.
func (s *Server) redirect(w http.ResponseWriter, r *http.Request) {
	if url, ok := s.urls[r.URL.Path]; ok {
		http.Redirect(w, r, url, http.StatusMovedPermanently)
		return
	}

	// no redirect found
	w.WriteHeader(http.StatusNotFound)
}

// AddHandler adds the provided Handler to the server.
func (s *Server) AddHandler(h Handler) {
	h.Register(s.mux)

	// append redirect URLs from the handler
	if urls := h.URLs(); len(urls) > 0 {
		for k, v := range urls {
			if link, ok := s.urls[k]; ok && link != v {
				glog.Warningf("redirect already exists for path %q (destination: %v). skipping", k, link)
				continue
			}
			s.urls[k] = v
		}
	}
}

// A Handler serves requests for short URLs.  Typically, a handler will
// register itself for a single content type prefix so that matching requests
// are routed to it.
type Handler interface {
	// Register the handler with the provided Router.  This method will be
	// called when the handler is added to the router, and allows the
	// handler to specify the kinds of short URLs it can handle.
	// Typically, but not always, this will be URLs of the form "/x" and
	// /x/*" where x is a particular content type.
	Register(*http.ServeMux)

	// URLs returns a map of URL paths (the shortened URLs) to the
	// destination URL it should be redirected to.
	URLs() map[string]string
}
