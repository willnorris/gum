// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

// Package gum provides the gum personal short URL resolver.
package gum

import "net/http"

// Server is a short URL redirection server.
type Server struct {
	mux *http.ServeMux
}

// NewServer constructs a new Server.
func NewServer() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

func (g *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// AddHandler adds the provided Handler to the server.
func (g *Server) AddHandler(h Handler) {
	h.Register(g.mux)
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
}
