//
// urlencode.go is a simple command line utility to encode a string in a URL friendly way.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2020, Caltech
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
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	// CaltechLibrary Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/mkpage"
)

var (
	description = `

SYNOPSIS

%s is a simple command line utility to URL encode content. By default
it reads from standard input and writes to standard out.  You can
also specifty the string to encode as a command line parameter.

`

	examples = `

EXAMPLES

    echo "This is the string to encode & nothing else!" | %s

would yield

    This%%20is%%20the%%20string%%20to%%20encode%%20&%%20nothing%%20else%%0A

`

	// Standard Options
	showHelp         bool
	showLicense      bool
	showVersion      bool
	showExamples     bool
	inputFName       string
	outputFName      string
	newLine          bool
	generateMarkdown bool
	generateManPage  bool
	quiet            bool

	// App Options
	useQueryEscape bool
)

func main() {
	app := cli.NewCli(mkpage.Version)
	appName := app.AppName()

	// Document non-option parameters
	app.SetParams(`[STRING_TO_ENCODE]`)

	// Add Help Docs
	app.AddHelp("license", []byte(fmt.Sprintf(mkpage.LicenseText, appName, mkpage.Version)))
	app.AddHelp("description", []byte(fmt.Sprintf(description, appName)))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName)))

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "examples", false, "display example(s)")
	app.StringVar(&inputFName, "i,input", "", "set input filename")
	app.StringVar(&outputFName, "o,output", "", "set output filename")
	app.BoolVar(&newLine, "nl,newline", false, "add a trailing newline to output")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generate man page")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")

	// App Options
	app.BoolVar(&useQueryEscape, "q,query", false, "use query escape (pluses for spaces)")

	app.Parse()
	args := app.Args()

	// Setup IO
	var err error

	app.Eout = os.Stderr

	app.In, err = cli.Open(inputFName, os.Stdin)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(inputFName, app.In)
	app.Out, err = cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(outputFName, app.Out)

	// Handle the default options
	if generateMarkdown {
		app.GenerateMarkdown(app.Out)
		os.Exit(0)
	}
	if generateManPage {
		app.GenerateManPage(app.Out)
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
	if showVersion {
		fmt.Println(app.Version())
		os.Exit(0)
	}
	if showLicense {
		fmt.Println(app.License())
		os.Exit(0)
	}

	nl := "\n"
	if newLine == false {
		nl = ""
	}

	var src string

	if len(args) > 0 {
		src = strings.Join(args, " ")
	} else {
		buf, err := ioutil.ReadAll(app.In)
		cli.ExitOnError(app.Eout, err, quiet)
		src = fmt.Sprintf("%s", buf)
	}
	if useQueryEscape == true {
		fmt.Fprintf(app.Out, "%s%s", url.QueryEscape(src), nl)
	} else {
		fmt.Fprintf(app.Out, "%s%s", url.PathEscape(src), nl)
	}
}
