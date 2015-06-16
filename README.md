# gum [![Build Status](https://travis-ci.org/willnorris/gum.svg?branch=master)](https://travis-ci.org/willnorris/gum) [![GoDoc](https://godoc.org/willnorris.com/go/gum?status.svg)](https://godoc.org/willnorris.com/go/gum) [![BSD License](https://img.shields.io/badge/license-BSD-blue.svg?style=flat)](LICENSE)

Gum is a personal short URL resolver written in Go.  It is primarily designed
to be used with statically generated sites.

This is the URL resolver I run behind [willnorris.com][] (and my short domain
[wjn.me][]).  For example, [wjn.me/c/gum](https://wjn.me/c/gum) redirects to
this project on GitHub.  So far, I've only implemented those handlers I use on
my own site, but everything should be easily usable by others.

[willnorris.com]: https://willnorris.com/
[wjn.me]: https://wjn.me/

## Getting Started ##

Install the package using:

    go get willnorris.com/go/gum/cmd/gum

Once installed, ensure `$GOPATH/bin` is in your `$PATH`, then run using:

    gum

This will start gum on port 4594.  Test this by navigating to
<http://localhost:4594/> and you should see a blank page (since no redirects
have been configured yet).

### Path Redirects ###

Path redirects allow a specified path prefix, and all URLs under that prefix,
to be redirected to a destination URL.  The portion of the URL after the prefix
is resolved relative to the destination.  Note that this means that the
destination should normally include a trailing slash.

Path redirects are specified with the `redirect` flag, which takes a value of
the form `prefix=destination`.  For example, to redirect all `/w` URLs to the
English version of Wikipedia, run:

    gum -redirect "w=https://en.wikipedia.org/wiki/"

Load <http://localhost:4594/w/URL_shortening> and you should be redirected to
the appropriate Wikipedia article.  The `redirect` flag can be repeated
multiple times.

### Static File Redirects ###

Gum can parse HTML file and automatically register redirects based on the links
specified in the files.  Gum will look for files with both a `rel="canonical"`
and `rel="shortlink"` link.  If found, it will configure redirects from the
shortlink URL path to the canonical URL.  For example, assume an HTML file with
the contents:

    <html>
      <link rel="canonical" href="http://example.com/post/123">
      <link rel="shortlink" href="http://x.com/t123">
    </html>

Gum will configure a redirect from `/t123` to `http://example.com/post/123`.
Note that only the path of the shortlink is used for creating the redirect.

Static file redirects are configured with the `static_dir` flag, which
identifies the root directory containing the HTML files.  Gum will recursively
parse all files with a `.html` file extension, looking for the appropriate link
tags.  It will additionally watch the specified directory for any changes and
will automatically load new or updated files.

Note that when using gum with a static site generator, `static_dir` should
identify the folder containing the generated HTML files (for example, the
`_site` folder when using jekyll), not the source files.

#### Alternate Short URLs ####

An HTML file can also specify multiple alternate short URLs to register for a
given canonical URL.  This is useful, for example, if you have legacy short
URLs that you want to continue resolving.  Alternate short URLs are specified
on the `rel="shortlink"` link as a space separated list of URLs in the
`data-alt-href` attribute.  For example:

    <link rel="shortlink" href="http://x.com/t123"
      data-alt-href="http://x.com/b/123 http://x.com/b/456">

Gum will resolve all of the shortlinks `/t123`, `/b/123`, and `/b/456` to the
relevant canonical URL.

## License ##

Gum is copyright Google, but is not an official Google product.  It is
available under a [BSD License][].

[BSD License]: LICENSE
