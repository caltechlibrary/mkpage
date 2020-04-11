//
// reldocpath.go takes a source document path and a target document path with same base path
// returning a relative path to the target file.
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
	"os"

	// My packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/mkpage"
)

var (
	description = `
SYNOPSIS

Given a source document path, a target document path calculate and
the implied common base path calculate the relative path for target.
`

	examples = `
EXAMPLE

Given

    %s chapter-01/lesson-03.html css/site.css

would output

    .../css/site.css
`

	// Standard options
	showHelp         bool
	showVersion      bool
	showLicense      bool
	showExamples     bool
	generateMarkdown bool
	generateManPage  bool
	quiet            bool
)

func main() {
	app := cli.NewCli(mkpage.Version)
	appName := app.AppName()

	// Define the command line parameters (non-options)
	app.SetParams(`SOURCE_DOC_PATH`, `TARGET_DOC_PATH`)

	// Configuration and command line interation
	app.AddHelp("license", []byte(fmt.Sprintf(mkpage.LicenseText, appName, mkpage.Version)))
	app.AddHelp("description", []byte(description))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName)))

	// Standard options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "examples", false, "display example(s)")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generate man page")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")

	app.Parse()
	args := app.Args()

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
	if showLicense {
		fmt.Println(app.License())
		os.Exit(0)
	}
	if showVersion {
		fmt.Println(app.Version())
		os.Exit(0)
	}

	if len(args) != 2 {
		cli.ExitOnError(app.Eout, fmt.Errorf("Expected a source and target file path\n For help try: %s -help", appName), quiet)
	}
	source, target := args[0], args[1]
	fmt.Fprintf(app.Out, `%s`, mkpage.RelativeDocPath(source, target))
}
