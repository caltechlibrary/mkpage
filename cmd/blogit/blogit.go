//
// mkpage is a thought experiment in a light weight template and
// markup (markdown, fountain) processing.
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
	"path"
	"strings"

	// Caltech Library packages
	"github.com/caltechlibrary/cli"
)

var (
	description = `
SYNOPSIS

Blogit provides a quick tool to add or replace blog content
organized around a date oriented file path. In addition to
placing documents it also will generate simple markdown documents
for inclusion in navigation.
`

	examples = `

EXAMPLE

I have a Markdown file called, "my-vacation-day.md". I want to
add it to my blog for the date July 1, 2020.  I've written
"my-vacation-day.md" in my home "Documents" folder and my blog
repository is in my "Sites" folder under "Sites/me.example.org".
Adding "my-vacation-day.md" to the blog me.example.org would
use the following command.

   cd $HOME/Sites/me.example.org
   %s $HOME/my-vacation-day.md 2020-07-01

The *%s* command will copy "my-vacation-day.md", creating any
necessary file directories to "$HOME/Sites/me.example.org/2020/06/01".
It will also update article lists (index.md) at the year level, 
month, and day level and month level of the directory tree and
and generate/update a posts.json in the "$HOME/Sites/my.example.org"
that can be used in your home page template for listing recent
posts.

*%s% includes an option to set the prefix path to
the blog posting.  In this way you could have separate blogs 
structures for things like podcasts or videocasts.

    # Add a landing page for the podcast
    %s -prefix=podcast my-vacation.md 2020-07-01
    # Add an audio file containing the podcast
    %s -prefix=podcast my-vacation.wav 2020-07-01

   -p, -prefix    Set the prefix path before the YYYY/MM/DD path.


`

	// Standard Options
	showHelp         bool
	showVersion      bool
	showVerbose bool
	showLicense      bool
	showExamples     bool
	quiet            bool
	generateMarkdown bool
	generateManPage  bool

	// Application Options
	prefixPath string

)

func main() {
	app := cli.NewCli(mkpage.Version)
	appName := app.AppName()

	// Document expected parameters
	app.SetParams(`DOCUMENT`, `[DATE]`)

	// Add Help docs
	app.AddHelp("license", []byte(fmt.Sprintf(mkpage.LicenseText, appName, mkpage.Version)))
	app.AddHelp("description", []byte(fmt.Sprintf(description)))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName, appName, appName, appName)))

	// Setup Environment variables
	app.EnvStringVar(&templateFNames, "MKPAGE_TEMPLATES", "", "set the default template path")

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showVerbose, "V,verbose", false, "verbose output")
	app.BoolVar(&showExamples, "e,examples", false, "verbose output")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generate man page")

	// Application specific options
	app.StringVar(&prefixPath, "p,prefix", "", "Set the prefix path before YYYY/MM/DD.")

	app.Parse()
	args := app.Args()

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

	if showVerbose {
		quiet = false
	}

	// Setup IO
	var err error

	app.Eout = os.Stderr


	// Process flags and update the environment as needed.
	if generateMarkdown {
		app.GenerateMarkdown(app.Out)
		os.Exit(0)
	}
	if generateManPage {
		app.GenerateManPage(app.Out)
		os.Exit(0)
	}
	err = mkpage.BlogIt(prefixPath, docName, dateString)
	cli.ExitOnError(app.Eout, err, quiet)
}
