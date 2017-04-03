//
// publisher.go - an experimental tool for maintaining a mkpage website hosted
// remotely.
//
// @Author R. S. Doiel, <rsdoiel@library.caltech.edu>
//
// Copyright (c) 2017, Caltech
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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/mkpage"
	"github.com/caltechlibrary/storage"
)

var (
	usage = `USAGE: %s [OPTIONS]`

	description = `
%s is a tool to interact with a mkpage project published on AWS S3.
It support the basic CRUD operations in your S3 bucket.
`

	examples = `
EXAMPLES

Examples assume you've previously setup your AWS access via environment 
variables on configuration files.

Create an index.html file in your AWS Bucket. 

	%s -create /index.html -local index.html

Read an item from the bucket (writes to stdout) 

	%s -read /index.html

Update (actually a delete, followed by create) a file

	%s -update /index.html -local index.html

Delete removes the copy in the bucket

    %s -delete /index.html
`

	// Standard Options
	showHelp    bool
	showLicense bool
	showVersion bool

	// Application Options
	createName string
	localName  string
	readName   string
	updateName string
	deleteName string
)

func init() {
	// Standard Option
	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showHelp, "help", false, "display help")
	flag.BoolVar(&showLicense, "l", false, "display license")
	flag.BoolVar(&showLicense, "license", false, "display license")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.BoolVar(&showVersion, "version", false, "display version")

	// Application Option
	flag.StringVar(&createName, "create", "", "create the named object in bucket")
	flag.StringVar(&readName, "read", "", "read the named object in bucket writing to stdout locally")
	flag.StringVar(&updateName, "update", "", "update (delete,create) an object in a bucket")
	flag.StringVar(&deleteName, "delete", "", "delete an item in the bucket")
	flag.StringVar(&localName, "local", "", "local name of object, needed by create and update options")
}

func main() {
	appName := path.Base(os.Args[0])
	flag.Parse()

	// Configuration and command line interation
	cfg := cli.New(appName, "MKPAGE", fmt.Sprintf(mkpage.LicenseText, appName, mkpage.Version), mkpage.Version)
	cfg.UsageText = fmt.Sprintf(usage, appName)
	cfg.DescriptionText = fmt.Sprintf(description, appName)
	cfg.ExampleText = fmt.Sprintf(examples, appName, appName, appName, appName)

	if showHelp == true {
		fmt.Println(cfg.Usage())
		os.Exit(0)
	}
	if showLicense == true {
		fmt.Println(cfg.License())
		os.Exit(0)
	}
	if showVersion == true {
		fmt.Println(cfg.Version())
		os.Exit(0)
	}

	if len(createName) > 0 && len(localName) == 0 {
		fmt.Fprintf(os.Stderr, "Missing local name for create option")
		os.Exit(0)
	}
	if len(updateName) > 0 && len(localName) == 0 {
		fmt.Fprintf(os.Stderr, "Missing local name for update option")
		os.Exit(0)
	}

	site, err := storage.Init(storage.S3)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	defer site.Close()

	switch {
	case len(createName) > 0 && len(localName) > 0:
		buf, err := ioutil.ReadFile(localName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		if err := site.Create(createName, buf); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	case len(readName) > 0:
		buf, err := site.Read(readName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "%s", buf)
	case len(updateName) > 0:
		buf, err := ioutil.ReadFile(localName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		if err := site.Update(createName, buf); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	case len(deleteName) > 0:
		if err := site.Delete(deleteName); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintln(os.Stderr, cfg.Usage())
		os.Exit(1)
	}
}
