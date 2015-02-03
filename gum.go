// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

// Package gum provides the gum personal short URL resolver.
package gum // import "willnorris.com/go/gum"

import (
	"net/http"
	"sync"

	"github.com/golang/glog"
)

// Server is a short URL redirection server.
type Server struct {
	// ServeMux which handles all incoming requests
	mux *http.ServeMux

	// mutex is a read/write lock for accessing the urls map
	mutex sync.RWMutex
	// map of short URL paths to destinations
	urls map[string]string

	// channel of static mappings of short URLs and their destinations.
	// Handlers can write to this channel to register new mappings; the
	// Server will read from this channel and handle serving the redirects.
	mappings chan Mapping
}

// NewServer constructs a new Server.
func NewServer() *Server {
	s := &Server{
		mux:      http.NewServeMux(),
		urls:     make(map[string]string),
		mappings: make(chan Mapping),
	}

	// the default handler serves redirects for registered Mapping values
	s.mux.HandleFunc("/", s.redirect)
	go s.readMappings()

	return s
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// redirect the request if a matching URL mapping has been configured.  If no
// mapping is found, a 404 status is returned.
func (s *Server) redirect(w http.ResponseWriter, r *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if url, ok := s.urls[r.URL.Path]; ok && url != "" {
		http.Redirect(w, r, url, http.StatusMovedPermanently)
		return
	}

	// no redirect found
	w.WriteHeader(http.StatusNotFound)
}

// readMappings reads values off the s.mappings channel and uses them to
// populate s.urls.  This method does not return.
func (s *Server) readMappings() {
	for {
		m := <-s.mappings

		s.mutex.Lock()
		if old, exists := s.urls[m.ShortPath]; exists {
			if m.Permalink == "" {
				glog.Infof("Deleting mapping: %v", m.ShortPath)
				delete(s.urls, m.ShortPath)
			} else if m.Permalink != old {
				glog.Warningf("Overwriting mapping: %v => %v (previously %q)", m.ShortPath, m.Permalink, old)
				s.urls[m.ShortPath] = m.Permalink
			}
		} else {
			glog.Infof("New mapping: %-7v => %v", m.ShortPath, m.Permalink)
			s.urls[m.ShortPath] = m.Permalink
		}
		s.mutex.Unlock()
	}
}

// AddHandler adds the provided Handler to the server.
func (s *Server) AddHandler(h Handler) {
	h.Register(s.mux)
	h.Mappings(s.mappings)
}

// Mapping represents a mapping between a short URL path and the permalink URL it is for.
type Mapping struct {
	// ShortPath is the path of the short URL (including leading slash) to
	// be mapped.
	ShortPath string

	// Permalink is the destination URL being mapped to.
	Permalink string
}

// A Handler serves requests for short URLs.  Typically, a handler will
// register itself for an entire path prefix using the Register func, or
// provide a list of static mappings using the Mappings func.
type Handler interface {
	// Register the handler with the provided ServeMux.  This method will be
	// called when the handler is added to the mux, and allows the
	// handler to specify the kinds of short URLs it can handle.
	// Typically, but not always, this will be URLs of the form "/x" and
	// /x/*" where x is a particular content type.
	Register(*http.ServeMux)

	// Mappings provides a write only channel for the handler to write
	// static Mapping values onto.  These mappings are then registered with
	// and the redirects handled by the Server.
	Mappings(chan<- Mapping)
}
