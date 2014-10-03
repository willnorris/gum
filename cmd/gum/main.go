// The gum binary starts up a gum server, configured for willnorris.com.
// Eventually, this should read from a configuration file instead of handlers
// being hardcoded here.  In the meantime, this should serve as a starting
// point for others.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"willnorris.com/go/gum"
	"willnorris.com/go/gum/jekyll"
)

// goxc values
var (
	// VERSION is the version string for gum.
	VERSION = "HEAD"

	// BUILD_DATE is the timestamp of when gum was built.
	BUILD_DATE string
)

var addr = flag.String("addr", "localhost:8080", "TCP address to listen on")
var version = flag.Bool("version", false, "print version information")

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("%v\nBuild: %v\n", VERSION, BUILD_DATE)
		return
	}
	if port := os.Getenv("PORT"); port != "" {
		a := "localhost:" + port
		addr = &a
	}

	g := gum.NewServer()

	// add handlers
	r, err := gum.NewRedirectHandler("w", "/wiki/")
	if err != nil {
		log.Fatal("error adding redirect handler: ", err)
	}
	g.AddHandler(r)

	r, err = gum.NewRedirectHandler("+", "https://plus.google.com/+willnorris/")
	if err != nil {
		log.Fatal("error adding redirect handler: ", err)
	}
	g.AddHandler(r)

	j, err := jekyll.NewHandler("/var/www/willnorris.com")
	if err != nil {
		log.Fatal("error adding jekyll handler: ", err)
	}
	g.AddHandler(j)

	server := &http.Server{
		Addr:    *addr,
		Handler: g,
	}

	fmt.Printf("gum (version %v) listening on %s\n", VERSION, server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
