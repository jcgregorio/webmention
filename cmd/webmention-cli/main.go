// Copyright 2014 Google Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file or at
// https://developers.google.com/open-source/licenses/bsd

// The webmention binary is a command line utiltiy for sending webmentions to
// the URLs linked to by a given webpage.
package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	"github.com/wsxiaoys/terminal/color"
	"willnorris.com/go/webmention"
)

const usageText = `webmention is a tool for sending webmentions.

Usage:
	webmention [flags] <url>

Flags:
`

var (
	client *webmention.Client
	input  string

	selector = flag.String("selector", ".h-entry", "CSS Selector limiting where to look for links")
)

func main() {
	flag.Parse()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageText)
		flag.PrintDefaults()
	}

	client = webmention.New(nil)
	input = flag.Arg(0)
	if input == "" {
		flag.Usage()
		return
	}

	if u, err := url.Parse(input); err != nil {
		fatalf("Not a valid URL: %q", input)
	} else if !u.IsAbs() {
		fatalf("URL %q is not an absolute URL", input)
	}

	fmt.Printf("Searching for links from %q to send webmentions to...\n\n", input)
	dl, err := client.DiscoverLinks(input, *selector)
	if err != nil {
		fatalf("error discovering links for %q: %v", input, err)
	}
	var links []link
	for _, l := range dl {
		links = append(links, link{url: l})
	}

	sendWebmentions(links)
}

type link struct {
	url string
}

func sendWebmentions(links []link) {
	fmt.Println("Sending webmentions...")
	for _, l := range links {
		fmt.Printf("  %v ... ", l.url)
		endpoint, err := client.DiscoverEndpoint(l.url)
		if err != nil {
			errorf("%v", err)
			continue
		} else if endpoint == "" {
			color.Println("@{!r}no webmention support@|")
			continue
		}

		_, err = client.SendWebmention(endpoint, input, l.url)
		if err != nil {
			errorf("%v", err)
			continue
		}
		color.Println("@gsent@|")
	}
}

func fatalf(format string, args ...interface{}) {
	errorf(format, args...)
	os.Exit(1)
}

func errorf(format string, args ...interface{}) {
	color.Fprintf(os.Stderr, "@{!r}ERROR:@| ")
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}
