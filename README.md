# gum [![Build Status](https://travis-ci.org/willnorris/gum.svg?branch=master)](https://travis-ci.org/willnorris/gum) [![GoDoc](https://godoc.org/willnorris.com/go/gum?status.svg)](https://godoc.org/willnorris.com/go/gum) [![BSD License](https://img.shields.io/badge/license-BSD-blue.svg?style=flat)](LICENSE)

Gum is a personal short URL resolver written in Go.

Gum is a short URL resolver for personal content (blog posts, photos, checkins,
etc).  Short URLs typically consist of a content type (normally a single
letter) plus an identifier for the resource (often a date-based encoded ID).
Exact content types and identifier format are configurable, but the overall
design is strongly modelled after [Whistle][].

This is the URL resolver I run behind <https://willnorris.com/>.  So far, I've
only implemented those handlers I use on my own site, but everything should be
easily usable by others.  Pass the `-help` flag to the 
[gum binary](cmd/gum/main.go) to see the configuration options.

[Whistle]: http://tantek.com/w/Whistle

## License ##

Gum is copyright Google, but is not an official Google product.  It is
available under a [BSD License][].

[BSD License]: LICENSE
