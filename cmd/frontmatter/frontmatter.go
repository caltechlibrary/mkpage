//
// frontmatter.go - is a command line tool that reads a Markdown file
// and returns the frontmatter portion.
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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	// 3rd Party packages
	"github.com/BurntSushi/toml"
	//"gopkg.in/yaml.v2"
	// ghodss Yaml implements the yaml 2 json func, wraps go-yaml
	"github.com/ghodss/yaml"

	// My packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/mkpage"
)

var (
	description = `
%s extracts a front matter from a Markdown file. If no front matter is present then an empty file is returned. Note %s doesn’t process the data extracted. It returns it unprocessed. Other tools can be used to process the front matter appropriately. By default %s reads from standard in and writes to standard out. This makes it very suitable for pipeline processing or for passing JSON formatted front matter back to mkpage for integration into the templates processed.
`

	examples = `
Extract a front matter from article.md.

    cat article.md | %s

This will display the front matter if found in article.md.

    %s -i article.md

Will also do the same.
`

	// Standard Options
	showHelp         bool
	showLicense      bool
	showVersion      bool
	showExamples     bool
	inputFName       string
	outputFName      string
	quiet            bool
	generateMarkdown bool
	generateManPage  bool

	// App Options
	jsonFormat bool
)

func main() {
	app := cli.NewCli(mkpage.Version)
	appName := app.AppName()

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "examples", false, "display example(s)")
	app.StringVar(&inputFName, "i,input", "", "input filename")
	app.StringVar(&outputFName, "o,output", "", "output filename")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate Markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generate man page")

	// Configuration and command line interation
	app.AddHelp("license", []byte(fmt.Sprintf(mkpage.LicenseText, appName, mkpage.Version)))
	app.AddHelp("description", []byte(fmt.Sprintf(description, appName, appName, appName)))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName, appName)))

	// App options
	app.BoolVar(&jsonFormat, "j,json", false, "output as JSON")

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

	// Handle Options
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
		} else if showExamples {
			fmt.Fprintln(app.Out, app.Help("examples"))
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

	//NOTE: read input and pass front matter to output.
	buf, err := ioutil.ReadAll(app.In)
	if err != nil {
		fmt.Fprintf(app.Eout, "%s", err)
		os.Exit(1)
	}
	_, frontMatterSrc, _ := mkpage.SplitFrontMatter(buf)
	if len(frontMatterSrc) > 0 {
		if jsonFormat {
			obj := make(map[string]interface{})
			switch {
			case bytes.HasPrefix(buf, []byte("+++\n")):
				// Make sure we have valid Toml
				if err := toml.Unmarshal(frontMatterSrc, &obj); err != nil {
					fmt.Fprintf(app.Eout, "Toml error: %s", err)
					os.Exit(1)
				}
			case bytes.HasPrefix(buf, []byte("%%%\n")):
				// Make sure we have valid Toml
				if err := toml.Unmarshal(frontMatterSrc, &obj); err != nil {
					fmt.Fprintf(app.Eout, "Toml error: %s", err)
					os.Exit(1)
				}
			case bytes.HasPrefix(buf, []byte("---\n")):
				if src, err := yaml.YAMLToJSON(frontMatterSrc); err != nil {
					fmt.Fprintf(app.Eout, "Yaml to JSON error: %s", err)
					os.Exit(1)
				} else {
					fmt.Fprintf(app.Out, "%s", src)
					os.Exit(0)
				}
			default:
				// Make sure we have valid JSON
				if err := json.Unmarshal(frontMatterSrc, &obj); err != nil {
					fmt.Fprintf(app.Eout, "JSON error: %s", err)
					os.Exit(1)
				}
			}
			if src, err := json.MarshalIndent(obj, "", "    "); err != nil {
				fmt.Fprintf(app.Eout, "%+v\n", obj)
				fmt.Fprintf(app.Eout, "%s\n", src)
				fmt.Fprintf(app.Eout, "JSON marshal error: %s", err)
				os.Exit(0)
			} else {
				fmt.Fprintf(app.Out, "%s", src)
				os.Exit(0)
			}
		}
		fmt.Fprintf(app.Out, "%s", frontMatterSrc)
	}
	os.Exit(0)
}
