//
// ws.go - A simple web server for static files.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2021, Caltech
// All rights not granted herein are expressly reserved by Caltech.
//
// Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/mkpage"
	"github.com/caltechlibrary/wsfn"
)

// Flag options
var (
	description = `

	a nimble web server

%s is a command line utility for developing and testing static websites.
It uses Go's standard http libraries and can supports both http 1 and 2
out of the box.  It is intended as a minimal wrapper for Go's standard
http libraries supporting http/https versions 1 and 2 out of the box.
`

	examples = `
Run web server using the content in the current directory

   %s

Run web server using a specified directory

   %s /www/htdocs
`

	// Standard options
	showHelp         bool
	showVersion      bool
	showLicense      bool
	showExamples     bool
	outputFName      string
	generateMarkdown bool
	quiet            bool

	// local app options
	uri          string
	docRoot      string
	sslKey       string
	sslCert      string
	CORSOrigin   string
	redirectsCSV string
)

func logRequest(r *http.Request) {
	log.Printf("Request: %s Path: %s RemoteAddr: %s UserAgent: %s\n", r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		next.ServeHTTP(w, r)
	})
}

func main() {
	app := cli.NewCli(mkpage.Version)
	appName := app.AppName()

	// Document non-option parameters
	app.SetParams(`[DOCROOT]`)

	// Add Help Docs
	app.AddHelp("license", []byte(fmt.Sprintf(mkpage.LicenseText, appName, mkpage.Version)))
	app.AddHelp("description", []byte(fmt.Sprintf(description, appName)))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName, appName)))

	defaultDocRoot := "."
	defaultURL := "http://localhost:8000"

	// Standard Options
	app.BoolVar(&showHelp, "h", false, "display help")
	app.BoolVar(&showHelp, "help", false, "display help")
	app.BoolVar(&showLicense, "l", false, "display license")
	app.BoolVar(&showLicense, "license", false, "display license")
	app.BoolVar(&showVersion, "v", false, "display version")
	app.BoolVar(&showVersion, "version", false, "display version")
	app.BoolVar(&showExamples, "example", false, "display example(s)")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate markdown documentation")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")

	// Application Options
	app.StringVar(&docRoot, "d,docs", defaultDocRoot, "Set the htdocs path")
	app.StringVar(&uri, "u,url", defaultURL, "The protocol and hostname listen for as a URL")
	app.StringVar(&sslKey, "k,key", "", "Set the path for the SSL Key")
	app.StringVar(&sslCert, "c,cert", "", "Set the path for the SSL Cert")
	app.StringVar(&CORSOrigin, "cors-origin", "*", "Set the CORS Origin Policy to a specific host or *")
	app.StringVar(&redirectsCSV, "redirects-csv", "", "Use target,destination replacement paths defined in CSV file")

	app.Parse()
	args := app.Args()

	// Setup IO
	var err error

	app.Eout = os.Stderr

	app.Out, err = cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(outputFName, app.Out)

	// Process flags and update the environment as needed.
	if generateMarkdown {
		app.GenerateMarkdown(app.Out)
		os.Exit(0)
	}
	if showHelp || showExamples {
		if len(args) > 0 {
			fmt.Fprintln(app.Out, app.Help(args...))
		} else {
			app.Usage(app.Out)
		}
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintln(app.Out, app.License())
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintln(app.Out, app.Version())
		os.Exit(0)
	}

	// setup from command line
	if len(args) > 0 {
		docRoot = args[0]
	}

	log.Printf("DocRoot %s", docRoot)

	u, err := url.Parse(uri)
	if err != nil {
		cli.ExitOnError(app.Eout, err, quiet)
	}

	if u.Scheme == "https" {
		log.Printf("SSL Key %s", sslKey)
		log.Printf("SSL Cert %s", sslCert)
	}
	log.Printf("Listening for %s", uri)
	cors := wsfn.CORSPolicy{
		Origin: CORSOrigin,
	}
	// Setup redirects defined the redirects CSV
	var rService *wsfn.RedirectService
	if redirectsCSV != "" {
		src, err := ioutil.ReadFile(redirectsCSV)
		if err != nil {
			log.Fatalf("Can't read %s, %s", redirectsCSV, err)
		}
		r := csv.NewReader(bytes.NewReader(src))
		// Allow support for comment rows
		r.Comment = '#'
		// Make a redirect map[string]string
		rmap := map[string]string{}
		for {
			row, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Can't read %s, %s", redirectsCSV, err)
			}
			if len(row) == 2 {
				// Define direct here.
				target := ""
				destination := ""
				if strings.HasPrefix(row[0], "/") {
					target = row[0]
				} else {
					target = "/" + row[0]
				}
				if strings.HasPrefix(row[1], "/") {
					destination = row[1]
				} else {
					destination = "/" + row[1]
				}
				rmap[target] = destination
			}
		}
		rService, err = wsfn.MakeRedirectService(rmap)
		if err != nil {
			log.Fatalf("Can't make redirect service, %s", err)
		}
	}
	http.Handle("/", cors.Handler(http.FileServer(http.Dir(docRoot))))

	if u.Scheme == "https" {
		if rService != nil {
			err = http.ListenAndServeTLS(u.Host, sslCert, sslKey, wsfn.RequestLogger(wsfn.StaticRouter(rService.RedirectRouter(http.DefaultServeMux))))
		} else {
			err = http.ListenAndServeTLS(u.Host, sslCert, sslKey, wsfn.RequestLogger(wsfn.StaticRouter(http.DefaultServeMux)))
		}
		cli.ExitOnError(app.Eout, err, quiet)
	} else {
		if rService != nil {
			err = http.ListenAndServe(u.Host, wsfn.RequestLogger(wsfn.StaticRouter(rService.RedirectRouter(http.DefaultServeMux))))
		} else {
			err = http.ListenAndServe(u.Host, wsfn.RequestLogger(wsfn.StaticRouter(http.DefaultServeMux)))
		}
		cli.ExitOnError(app.Eout, err, quiet)
	}
}
