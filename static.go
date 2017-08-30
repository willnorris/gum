// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

package gum

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	fsnotify "gopkg.in/fsnotify.v1"
)

const (
	relShortlink = "shortlink"
	relCanonical = "canonical"
	attrAltHref  = "data-alt-href"
)

// StaticHandler handles short URLs parsed from static HTML files.  Files are
// parsed and searched for rel="shortlink" and rel="canonical" links.  If both
// are found, a redirect is registered for the pair.
type StaticHandler struct {
	base    string
	watcher *fsnotify.Watcher
}

// NewStaticHandler constructs a new StaticHandler with the specified base path
// of HTML files.
func NewStaticHandler(base string) (*StaticHandler, error) {
	if stat, err := os.Stat(base); err != nil {
		return nil, err
	} else if !stat.IsDir() {
		return nil, fmt.Errorf("Specified base path %q is not a directory", base)
	}

	return &StaticHandler{base: base}, nil
}

// Mappings implements Handler.
func (h *StaticHandler) Mappings(mappings chan<- Mapping) error {
	if err := loadFiles(h.base, mappings); err != nil {
		return err
	}

	var err error
	h.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return errors.Wrap(err, "error creating file watcher")
	}

	go func() {
		for {
			select {
			case ev := <-h.watcher.Events:
				if ev.Op&(fsnotify.Remove|fsnotify.Rename) != 0 {
					// ignore Remove and Rename events
					continue
				}

				stat, err := os.Stat(ev.Name)
				if err != nil {
					log.Printf("Error reading file stats for %q: %v", ev.Name, err)
					continue
				}

				// add watcher for newly created directories
				if ev.Op&fsnotify.Create == fsnotify.Create && stat.IsDir() {
					h.watcher.Add(ev.Name)
				}

				// if event is Create or Write, reload files
				if ev.Op&(fsnotify.Create|fsnotify.Write) != 0 {
					if err := loadFiles(ev.Name, mappings); err != nil {
						log.Print(err)
					}
				}
			case err := <-h.watcher.Errors:
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	// setup initial file watchers for h.base and all sub-directories
	err = filepath.Walk(h.base, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() {
			err = h.watcher.Add(path)
			if err != nil {
				return errors.Wrapf(err, "error watching path %q", path)
			}
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "error setting up watchers for %q", h.base)
	}
	return nil
}

// Register is a noop for this handler.
func (h *StaticHandler) Register(mux *http.ServeMux) error { return nil }

func loadFiles(base string, mappings chan<- Mapping) error {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".html" {
			// skip directories and non-HTML files
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		fileMappings, err := parseFile(f)
		if err != nil {
			return errors.Wrapf(err, "error parsing file %q", path)
		}

		for _, m := range fileMappings {
			mappings <- m
		}
		return nil
	}

	err := filepath.Walk(base, walkFn)
	if err != nil {
		return errors.Wrapf(err, "error walking %q", base)
	}
	return nil
}

// parseFile parses r as HTML and returns the URLs of the first links found
// with the "shortlink" and "canonical" rel values.
func parseFile(r io.Reader) (mappings []Mapping, err error) {
	var permalink string
	var shortlinks []string

	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.DataAtom == atom.Link || n.DataAtom == atom.A {
				var href, rel, altHref string
				for _, a := range n.Attr {
					if a.Key == atom.Href.String() {
						href = a.Val
					}
					if a.Key == atom.Rel.String() {
						rel = a.Val
					}
					if a.Key == attrAltHref {
						altHref = a.Val
					}
				}
				if href != "" && rel != "" {
					for _, v := range strings.Split(rel, " ") {
						if v == relShortlink {
							shortlinks = append(shortlinks, href)
							shortlinks = append(shortlinks, strings.Split(altHref, " ")...)
						}
						if v == relCanonical && permalink == "" {
							permalink = href
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	if len(shortlinks) > 0 && permalink != "" {
		for _, link := range shortlinks {
			shorturl, err := url.Parse(link)
			if err != nil {
				log.Printf("error parsing shortlink %q: %v", link, err)
			}
			if path := shorturl.Path; len(path) > 1 {
				mappings = append(mappings, Mapping{ShortPath: path, Permalink: permalink})
			}
		}
	}

	return mappings, nil
}
