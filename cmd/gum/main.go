// The gum binary starts up a gum server.  See the usage text for configuration
// options.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"willnorris.com/go/gum"
)

// goxc values
var (
	// VERSION is the version string for gum.
	VERSION = "HEAD"

	// BUILD_DATE is the timestamp of when gum was built.
	BUILD_DATE string
)

// Flags
var (
	addr      = flag.String("addr", "localhost:8002", "TCP address to listen on")
	version   = flag.Bool("version", false, "print version information")
	staticDir = flag.String("static_dir", "", "directory of static site to setup redirects for")
	redirects redirectSlice
)

func init() {
	flag.Var(&redirects, "redirect", "redirect handler definition of the form 'prefix=destination'")
}

type redirect struct {
	Prefix, Destination string
}

type redirectSlice []redirect

func (r *redirectSlice) String() string {
	return fmt.Sprintf("%v", *r)
}

func (r *redirectSlice) Set(value string) error {
	parts := strings.SplitN(value, "=", 2)
	if len(parts) != 2 {
		return errors.New("redirect flag value should be of the form 'prefix=dest'")
	}
	if _, err := url.Parse(parts[1]); err != nil {
		return fmt.Errorf("Destination %q is not a valid URL: %v", parts[1], err)
	}
	*r = append(*r, redirect{Prefix: parts[0], Destination: parts[1]})
	return nil
}

func usage() {
	fmt.Print(`gum is a personal short URL resolver.
Usage:
  gum [-redirect=<redirect>] [-static_dir=<static_dir>]

Gum supports two styles of handlers, which are configured with command line flags:

Redirect Handlers are configured by providing a mapping of the form
"prefix=dest" using the -redirect flag, which will redirect all URLs matching a
specified prefix to a given destination URL.  For example, you could call gum
with the following flags:

  gum -redirect x=http://example.com/

This would redirect all requests whose path begins with "/x" to the
corresponding URL on "http://example.com/".  For example, the request
"/x/a/b?c=d" would be redirected to "http://example.com/a/b?c=d".  Destination
URLs can be absolute or relative.


A Static Site Handler is configured by passing the root directory of a static
website using the -static_dir flag.  This handler will parse all HTML files
under that directory looking for files with both a rel="canonical" and
rel="shortlink" link.  If found, it will configure redirects from the shortlink
URL path to the canonical URL.  In addition, the rel="shortlink" link can
provide a space-separate list of additional shortlinks in the 'data-alt-href'
attribute, which will also be redirected.  This is mostly useful for legacy
shortlinks that should be redirected, but that you don't standard parsers to
see.  For example, given an HTML file with the contents:

  <html>
    <link rel="canonical" href="http://example.com/post/12345678">
    <link rel="shortlink" href="http://x.com/t123"
      data-alt-href="http://x.com/b/123 http://x.com/b/456">
  </html>

Requests whose path matched any of "/t123", "/b/123", or "/b/456" would be
redirected to "http://example.com/post/12345678".  The static site handler
currently ignores the hostname of the shortlink URL when setting up redirects.

Flags:
`)
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Printf("%v\nBuild: %v\n", VERSION, BUILD_DATE)
		return
	}
	if port := os.Getenv("PORT"); *addr == "" && port != "" {
		a := "localhost:" + port
		addr = &a
	}

	g := gum.NewServer()

	for _, r := range redirects {
		h, err := gum.NewRedirectHandler(r.Prefix, r.Destination)
		if err != nil {
			log.Fatal("error adding redirect handler: ", err)
		}
		g.AddHandler(h)
	}

	if *staticDir != "" {
		h, err := gum.NewStaticHandler(*staticDir)
		if err != nil {
			log.Fatal("error adding static handler: ", err)
		}
		g.AddHandler(h)
	}

	server := &http.Server{
		Addr:    *addr,
		Handler: g,
	}

	fmt.Printf("gum (version %v) listening on %s\n", VERSION, server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
