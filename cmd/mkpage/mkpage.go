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
	"github.com/caltechlibrary/mkpage"
	"github.com/caltechlibrary/tmplfn"
)

var (
	description = `
SYNOPSIS

Using the key/value pairs populate the template(s) and render to stdout.
`

	examples = `

EXAMPLE

Template (named "examples/weather.tmpl")
    
    Date: $now$
    
    Hello $name$,
        
    The weather forcast is
    
    $if(weather.data.weather)$
      $weather.data.weather[; ]$
    $endif$
    
    Thank you
    
    $signature$

Render the template above (i.e. examples/weather.tmpl) would be 
accomplished from the following data sources--

 + "now" and "name" are strings
 + "weather" is JSON data retrieved from a URL
 	+ ".data.weather" is a data path inside the JSON document
	+ "index" let's us pull our the "0"-th element (i.e. the initial element of the array)
 + "signature" comes from a file in our local disc (i.e. examples/signature.txt)

That would be expressed on the command line as follows

    %s "now=text:$(date)" "name=text:Little Frieda" \
        "weatherForecast=http://forecast.weather.gov/MapClick.php?lat=13.47190933300044&lon=144.74977715100056&FcstType=json" \
        signature=examples/signature.txt \
        examples/weather.tmpl     

`

	// Standard Options
	showHelp         bool
	showVersion      bool
	showLicense      bool
	showExamples     bool
	inputFName       string
	outputFName      string
	quiet            bool
	generateMarkdown bool
	generateManPage  bool

	// Application Options
	templateFNames string
	showTemplate   bool
	codesnip       bool
	codeType       string
	useGoTemplates bool
)

func main() {
	app := cli.NewCli(mkpage.Version)
	appName := app.AppName()

	// Document expected parameters
	app.SetParams(`[KEY/VALUE DATA PAIRS]`, `[TEMPLATE_FILENAMES]`)

	// Add Help docs
	app.AddHelp("license", []byte(fmt.Sprintf(mkpage.LicenseText, appName, mkpage.Version)))
	app.AddHelp("description", []byte(fmt.Sprintf(description)))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName)))

	// Setup Environment variables
	app.EnvStringVar(&templateFNames, "MKPAGE_TEMPLATES", "", "set the default template path")

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "examples", false, "display example(s)")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.StringVar(&inputFName, "i,input", "", "input filename")
	app.StringVar(&outputFName, "o,output", "", "output filename")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generate man page")

	// Application specific options
	app.BoolVar(&showTemplate, "s,show-template", false, "display source for a default page template")
	app.StringVar(&templateFNames, "t,templates", "", "colon delimited list of templates to use")
	app.BoolVar(&codesnip, "codesnip", false, "output just the code bocks, reads from standard input")
	app.StringVar(&codeType, "code", "", "outout just code blocks for specific language, e.g. shell or json, reads from standard input")
	app.BoolVar(&useGoTemplates, "gt,go-templates", false, "use Go's template engine instead of Pandoc's template engine")

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

	if showTemplate {
		fmt.Fprintf(app.Out, "%s\n", mkpage.Defaults["/pandoc/page.tmpl"])
		os.Exit(0)
	}

	// Default template name is page.tmpl
	templateName := ""
	templateSources := []string{}

	// Make sure we have a configured command to run
	if len(templateFNames) > 0 {
		for _, fname := range strings.Split(templateFNames, ":") {
			templateSources = append(templateSources, fname)
		}
	}

	// Setup IO
	var err error

	app.Eout = os.Stderr
	app.In, err = cli.Open(inputFName, os.Stdin)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(inputFName, app.In)

	app.Out, err = cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(outputFName, app.Out)

	// Process flags and update the environment as needed.
	if generateMarkdown {
		app.GenerateMarkdown(app.Out)
		os.Exit(0)
	}
	if generateManPage {
		app.GenerateManPage(app.Out)
		os.Exit(0)
	}

	if codesnip || codeType != "" {
		err = mkpage.Codesnip(app.In, app.Out, codeType)
		cli.ExitOnError(app.Eout, err, quiet)
		os.Exit(0)
	}

	data := map[string]string{}
	for i, arg := range args {
		if strings.Contains(arg, "=") == true {
			// Update data map
			pair := strings.SplitN(arg, "=", 2)
			if len(pair) != 2 {
				fmt.Fprintf(app.Eout, "Can't read pair (%d) %s\n", i+1, arg)
				os.Exit(1)
			}
			data[pair[0]] = pair[1]
		} else {
			// Must be the template source
			templateSources = append(templateSources, arg)
		}
	}

	// Make the page with pandoc, go templates and Go Markdown
	switch {
	case useGoTemplates:
		// DEPRECIATED: This is maintained for backard compatibility
		// for now. It should be removed when the transition is
		// completed.

		// Create our Tmpl struct with our function map
		tmpl := tmplfn.New(tmplfn.AllFuncs())

		// Load any user supplied templates
		if len(templateSources) > 0 {
			err = tmpl.ReadFiles(templateSources...)
			cli.ExitOnError(app.Eout, err, quiet)
			templateName = path.Base(templateSources[0])
		} else {
			// Load our default template maps
			err = tmpl.Add(templateName, mkpage.Defaults["/templates/page.tmpl"])
			cli.ExitOnError(app.Eout, err, quiet)
		}
		// Build a template and send to MakePage
		t, err := tmpl.Assemble()
		cli.ExitOnError(app.Eout, err, quiet)
		err = mkpage.MakePage(app.Out, templateName, t, data)
	default:
		if len(templateSources) > 0 {
			templateName = templateSources[0]
		}
		err = mkpage.MakePandoc(app.Out, templateName, data)
	}
	cli.ExitOnError(app.Eout, err, quiet)
}
