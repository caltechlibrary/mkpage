// mkrss.go is a command line tool for generating an RSS file from a blog
// directory structure in the form of PATH_TO_BLOG/YYYY/MM/DD/BLOG_ARTICLES.html
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
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	// My packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/mkpage"
	"github.com/caltechlibrary/rss2"
)

var (
	// Usage and docs

	description = `
SYNOPSIS

%s walks the file system to generate a RSS2 file. It assumes 
that the directory for HTDOCS is is the base directory containing 
subdirectories in the form of /YYYY/MM/DD/ARTICLE_HTML where 
YYYY/MM/DD (Year, Month, Day) corresponds to the publication date 
of ARTICLE_HTML.
`

	examples = `
EXAMPLE

If our htdocs folder is our document root and out blog is
htdocs/myblog.

    %s -channel-title="This Great Beyond" \
        -channel-description="Blog to save the world" \
        -channel-link="http://blog.example.org" \
        htdocs htdocs/rss.xml

This would build an RSS 2 file in htdocs/rss.xml from the
articles in htdocs/myblog/YYYY/MM/DD.
`

	// Standard options
	showHelp         bool
	showLicense      bool
	showVersion      bool
	showExamples     bool
	inputFName       string
	outputFName      string
	quiet            bool
	generateMarkdown bool
	generateManPage  bool

	// App specific options
	excludeList        string
	articleLimit       int
	channelLanguage    string
	channelTitle       string
	channelDescription string
	channelLink        string
	channelGenerator   string
	channelPubDate     string
	channelBuildDate   string
	channelCopyright   string
	channelCategory    string
	bylineExp          string
	titleExp           string
	dateExp            string
)

func main() {

	app := cli.NewCli(mkpage.Version)
	appName := app.AppName()

	// App Parameters (non-options)
	app.SetParams(`HTDOCS`, `[RSS_FILENAME]`)

	// Add Help Docs
	app.AddHelp("license", []byte(fmt.Sprintf(mkpage.LicenseText, appName, mkpage.Version)))
	app.AddHelp("description", []byte(fmt.Sprintf(description, appName)))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName)))

	// Standard options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "examples", false, "display example(s)")
	app.StringVar(&inputFName, "i,input", "", "set input filename")
	app.StringVar(&outputFName, "o,output", "", "set output filename")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generate man page")

	// App specific options
	app.StringVar(&excludeList, "e", "", "A colon delimited list of path exclusions")
	app.IntVar(&articleLimit, "c", 0, "If non-zero, limit the number of articles in the RSS file")
	app.StringVar(&channelLanguage, "channel-language", "", "Language, e.g. en-ca")
	app.StringVar(&channelTitle, "channel-title", "", "Title of channel")
	app.StringVar(&channelDescription, "channel-description", "", "Description of channel")
	app.StringVar(&channelLink, "channel-link", "", "link to channel")
	app.StringVar(&channelGenerator, "channel-generator", "", "Name of RSS generator")
	app.StringVar(&channelPubDate, "channel-pubdate", "", "Pub Date for channel (e.g. 2006-01-02 15:04:05 -0700)")
	app.StringVar(&channelBuildDate, "channel-builddate", "", "Build Date for channel (e.g. 2006-01-02 15:04:05 -0700)")
	app.StringVar(&channelCopyright, "channel-copyright", "", "Copyright for channel")
	app.StringVar(&channelCategory, "channel-category", "", "category for channel")
	app.StringVar(&dateExp, "d,date-format", mkpage.DateExp, "set date regexp")
	app.StringVar(&titleExp, "t,title", mkpage.TitleExp, "set title regexp")
	app.StringVar(&bylineExp, "b,byline", mkpage.BylineExp, "set byline regexp")

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

	// Process options
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
		fmt.Fprintln(app.Out, app.License())
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintln(app.Out, app.Version())
		os.Exit(0)
	}

	if len(channelTitle) == 0 {
		channelTitle = `A website`
	}
	if len(channelDescription) == 0 {
		channelDescription = `A collection of web pages`
	}
	if len(channelLink) == 0 {
		channelLink = `http://localhost:8000`
	}

	// Setup the Channel metadata for feed.
	feed := new(rss2.RSS2)
	feed.Version = "2.0"
	feed.Title = channelTitle
	feed.Description = channelDescription
	feed.Link = channelLink
	if len(channelLanguage) > 0 {
		feed.Language = channelLanguage
	}
	if len(channelCopyright) > 0 {
		feed.Copyright = channelCopyright
	}
	if len(channelCategory) > 0 {
		feed.Category = channelCategory
	}
	if len(channelGenerator) == 0 {
		feed.Generator = app.Version()
	} else {
		feed.Generator = channelGenerator
	}
	now := time.Now()
	if len(channelPubDate) == 0 {
		// RSS spec shows RTF 1123 dates
		//feed.PubDate = now.Format(time.RFC822Z)
		feed.PubDate = now.Format(time.RFC1123)
	} else {
		dt, err := mkpage.NormalizeDate(channelPubDate)
		if err != nil {
			cli.ExitOnError(app.Eout, fmt.Errorf("Can't parse %q, %s\n", channelPubDate, err), quiet)
		}
		feed.PubDate = dt.Format(time.RFC1123)
	}
	if len(channelBuildDate) == 0 {
		// RSS spec shows RTF 1123 dates
		//feed.LastBuildDate = now.Format(time.RFC822Z)
		feed.LastBuildDate = now.Format(time.RFC1123)
	} else {
		dt, err := mkpage.NormalizeDate(channelBuildDate)
		if err != nil {
			cli.ExitOnError(app.Eout, fmt.Errorf("Can't parse %q, %s\n", channelBuildDate, err), quiet)
		}
		feed.LastBuildDate = dt.Format(time.RFC1123)
	}

	// Process command line parameters
	htdocs := "."
	rssPath := ""
	if len(args) > 0 {
		htdocs = args[0]
	}
	if len(args) > 1 {
		rssPath = args[1]
	}

	validBlogPath := regexp.MustCompile("/[0-9][0-9][0-9][0-9]/[0-9][0-9]/[0-9][0-9]/")
	err = mkpage.Walk(htdocs, func(p string, info os.FileInfo) bool {
		fname := path.Base(p)
		if validBlogPath.MatchString(p) == true &&
			strings.HasSuffix(fname, ".md") == true {
			// NOTE: We have a possible published markdown article.
			// Make sure we have a HTML version before adding it
			// to the feed.
			if _, err := os.Stat(path.Join(p, path.Base(fname)+".html")); os.IsNotExist(err) {
				return false
			}
			return true
		}
		return false
	}, func(p string, info os.FileInfo) error {
		// Read the article
		buf, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		// Calc URL path
		pname := strings.TrimPrefix(p, htdocs)
		if strings.HasPrefix(pname, "/") {
			pname = strings.TrimPrefix(pname, "/")
		}
		dname := path.Dir(pname)
		bname := strings.TrimSuffix(path.Base(pname), ".md") + ".html"
		articleURL := fmt.Sprintf("%s/%s", channelLink, path.Join(dname, bname))
		u, err := url.Parse(articleURL)
		if err != nil {
			return err
		}
		// Collect metadata
		src := fmt.Sprintf("%s", buf)
		title := strings.TrimPrefix(mkpage.Grep(titleExp, src), "# ")
		byline := mkpage.Grep(bylineExp, src)
		pubDate := mkpage.Grep(dateExp, byline)
		author := byline
		if len(byline) > 2 {
			author = strings.TrimSpace(strings.TrimSuffix(byline[2:], pubDate))
		}
		// Reformat pubDate to conform to RSS2 date formats
		dt, err := time.Parse(`2006-01-02`, pubDate)
		if err != nil {
			return err
		}
		pubDate = dt.Format(time.RFC1123)
		item := new(rss2.Item)
		item.Title = title
		item.Author = author
		item.PubDate = pubDate
		item.Link = u.String()
		feed.ItemList = append(feed.ItemList, *item)
		return nil
	})
	if err != nil {
		fmt.Fprintf(app.Eout, "%s\n", err)
		os.Exit(1)
	}

	// Marshal RSS2 and render output
	src, err := xml.MarshalIndent(feed, "", "    ")
	if err != nil {
		fmt.Fprintf(app.Eout, "%s\n", err)
		os.Exit(1)
	}
	txt := fmt.Sprintf(`<?xml version="1.0"?>
%s`, src)
	if len(rssPath) > 0 {
		err = ioutil.WriteFile(rssPath, []byte(txt), 0664)
		cli.ExitOnError(app.Eout, err, quiet)
		os.Exit(0)
	}
	fmt.Fprintf(app.Out, "%s\n", txt)
}
