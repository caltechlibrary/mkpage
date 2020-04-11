//
// Package mkpage is an experiment in a light weight template and markdown processor.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2020, Caltech
// All rights not granted herein are expressly reserved by Caltech
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
package mkpage

import (
	"bytes"
	"encoding/json"
	"fmt"
	//"go/ast"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"time"

	// 3rd Party Packages
	"github.com/BurntSushi/toml"
	"github.com/ghodss/yaml"

	// Gomarkdown implementation of markdown
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	// Mmarkdown implementation
	"github.com/mmarkdown/mmark/lang"
	"github.com/mmarkdown/mmark/mast"
	"github.com/mmarkdown/mmark/mparser"
	"github.com/mmarkdown/mmark/render/mhtml"

	// Fountain support for scripts, interviews and narration
	"github.com/rsdoiel/fountain"
)

const (
	// Version holds the semver assocaited with this version of mkpage.
	Version = `v0.0.32`

	// LicenseText provides a string template for rendering cli license info
	LicenseText = `
%s %s

Copyright (c) 2019, Caltech
All rights not granted herein are expressly reserved by Caltech.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
`

	// Prefix for explicit string types

	// JSONPrefix designates a string as JSON formatted content
	JSONPrefix = "json:"
	// MMarkPrefix designates a string as mmarkdown or MMark content
	MMarkPrefix = "mmark:"
	// MarkdownPrefix designates a string as Markdown (common mark) content
	// to be parsed by 'defaultProcessor' (i.e. gomarkdown)
	MarkdownPrefix = "markdown:"
	// GomarkdownPrefix designates a string to be parsed explicitly by
	// Gomarkdown
	GomarkdownPrefix = "gomarkdown:"
	// TextPrefix designates a string as text/plain not needed processing
	TextPrefix = "text:"
	// FountainPrefix designates a string as Fountain formatted content
	FountainPrefix = "fountain:"

	// SOMEDAY: should add XML, BibTeX, YaML support...

	// DateExp is the default format used by mkpage utilities for date exp
	DateExp = `[0-9][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9]`
	// BylineExp is the default format used by mkpage utilities
	BylineExp = (`^[B|b]y\s+(\w|\s|.)+` + DateExp + "$")
	// TitleExp is the default format used by mkpage utilities
	TitleExp = `^#\s+(\w|\s|.)+$`

	// Config types for Front Matter

	// ConfigIsUnknown means front matter and we can't parse it
	ConfigIsUnknown = iota
	// ConfigIsYAML means that YAML has been detected in the front matter (per Hugo style fencing)
	ConfigIsYAML
	// ConfigIsTOML means we have TOML Front Matter based on Hugo fencing or Mmarkdown fencing
	ConfigIsTOML
	// ConfigIsJSON means we have detected JSON front matter
	ConfigIsJSON
)

var (
	// DefaultTemplateSource is defined in init by Defaults["/templates/page.tmpl"]
	// Defaults is a map to assets defined in assets.go which is build with pkgasset and
	// the contents of the defaults folder in this repository.
	DefaultTemplateSource string

	// DefaultSlideTemplateSource provides the default HTML template for mkslides package,
	// you probably want to override this... is defined in init by Defaults["/templates/slides.tmpl"]
	// Defaults is a map to assets defined in assets.go which is build with pkgasset and
	// the contents of the defaults folder in this repository.
	DefaultSlideTemplateSource string

	// Config holds a global config.
	// Uses the same structure as Front Matter in that it is
	// the result of parsing TOML, YAML or JSON into a
	// map[string]interface{} tree
	Config map[string]interface{}
)

// normalizeEOL takes a []byte and normalizes the end of line
// to a `\n' from `\r\n` and `\r`
func normalizeEOL(input []byte) []byte {
	if bytes.Contains(input, []byte("\r\n")) {
		input = bytes.Replace(input, []byte("\r\n"), []byte("\n"), -1)
	}
	/*
		if bytes.Contains(input, []byte("\r")) {
			input = bytes.Replace(input, []byte("\r"), []byte("\n"), -1)
		}
	*/
	return input
}

// SplitFrontMatter takes a []byte input splits it into front matter
// source and Markdown source. If either is missing an empty []byte
// is returned for the missing element.
func SplitFrontMatter(input []byte) (int, []byte, []byte) {
	// YAML front matter uses ---, note this conflicts with Mmark practice, do I want to support YAML like this?
	if bytes.HasPrefix(input, []byte("---\n")) {
		parts := bytes.SplitN(bytes.TrimPrefix(input, []byte("---\n")), []byte("\n---\n"), 2)
		return ConfigIsYAML, parts[0], parts[1]
	}
	// TOML front matter as used in Hugo
	if bytes.HasPrefix(input, []byte("+++\n")) {
		parts := bytes.SplitN(bytes.TrimPrefix(input, []byte("+++\n")), []byte("\n+++\n"), 2)
		return ConfigIsTOML, parts[0], parts[1]
	}
	// TOML front matter identified in Mmark as three % or dashes,
	// We can support the %, dashes are taken by Hugo style, but
	// maybe I don't want to support that?
	if bytes.HasPrefix(input, []byte("%%%\n")) {
		parts := bytes.SplitN(bytes.TrimPrefix(input, []byte("%%%\n")), []byte("\n%%%\n"), 2)
		return ConfigIsTOML, parts[0], parts[1]
	}
	// JSON front matter, most Markdown processors.
	if bytes.HasPrefix(input, []byte("{\n")) {
		parts := bytes.SplitN(bytes.TrimPrefix(input, []byte("{\n")), []byte("\n}\n"), 2)
		src := []byte(fmt.Sprintf("{\n%s\n}\n", parts[0]))
		return ConfigIsJSON, src, parts[1]
	}
	// Handle case of no front matter
	return ConfigIsUnknown, []byte(""), input
}

// ProcessorConfig takes front matter and returns
// a map[string]interface{} containing configuration
func ProcessorConfig(configType int, frontMatterSrc []byte) (map[string]interface{}, error) {
	//FIXME: Need to merge with .Config and return the merged result.
	m := map[string]interface{}{}
	// Do nothing is we have zero front matter to process.
	if len(frontMatterSrc) == 0 {
		return m, nil
	}
	// Convert Front Matter to JSON
	switch configType {
	case ConfigIsYAML:
		// YAML Front Matter
		jsonSrc, err := yaml.YAMLToJSON(frontMatterSrc)
		if err != nil {
			return nil, fmt.Errorf("Can't parse YAML front matter, %s", err)
		}
		if err = json.Unmarshal(jsonSrc, &m); err != nil {
			return nil, err
		}
	case ConfigIsTOML:
		// TOML Front Matter
		if _, err := toml.Decode(fmt.Sprintf("%s", frontMatterSrc), &m); err != nil {
			return nil, err
		}
	case ConfigIsJSON:
		// JSON Front Matter
		if err := json.Unmarshal(frontMatterSrc, &m); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown front matter format")
	}
	return m, nil
}

// mmarkExtensions takes a config (map[string]interface{}) and
// returns the ORed exentions flags for map["mmark"]
func mmarkExtensions(config map[string]interface{}) parser.Extensions {
	ext := parser.NoExtensions
	if thing, ok := config["gomarkdown"]; ok == true {
		extensions := thing.(map[string]interface{})
		for k, v := range extensions {
			onoff := false
			option := k
			switch v.(type) {
			case int:
				onoff = v.(int) == 1
			case int64:
				onoff = v.(int64) == 1
			case bool:
				onoff = v.(bool)
			case string:
				onoff = strings.ToLower(v.(string)) == "true"
			}
			if onoff {
				switch option {
				case "NoIntraEmphasis":
					ext |= parser.NoIntraEmphasis
				case "Tables":
					ext |= parser.Tables
				case "FencedCode":
					ext |= parser.FencedCode
				case "Autolink":
					ext |= parser.Autolink
				case "Strikethrough":
					ext |= parser.Strikethrough
				case "LaxHTMLBlocks":
					ext |= parser.LaxHTMLBlocks
				case "SpaceHeadings":
					ext |= parser.SpaceHeadings
				case "HardLineBreak":
					ext |= parser.HardLineBreak
				case "TabSizeEight":
					ext |= parser.TabSizeEight
				case "Footnotes":
					ext |= parser.Footnotes
				case "NoEmptyLineBeforeBlock":
					ext |= parser.NoEmptyLineBeforeBlock
				case "HeadingIDs":
					ext |= parser.HeadingIDs
				case "Titleblock":
					ext |= parser.Titleblock
				case "AutoHeadingIDs":
					ext |= parser.AutoHeadingIDs
				case "BackslashLineBreak":
					ext |= parser.BackslashLineBreak
				case "DefinitionLists":
					ext |= parser.DefinitionLists
				case "MathJax":
					ext |= parser.MathJax
				case "OrderedListStart":
					ext |= parser.OrderedListStart
				case "Attributes":
					ext |= parser.Attributes
				case "SuperSubscript":
					ext |= parser.SuperSubscript
				case "CommonExtensions":
					ext |= parser.CommonExtensions
				}
			}
		}
	}
	return ext
}

// mmarkRenderOptions config (map[string]interface{}) and
// returns the ORed extenions flags for map["mmark"].
func mmarkRenderOptions(config map[string]interface{}) html.Flags {
	flag := html.FlagsNone
	if thing, ok := config["mmark"]; ok == true {
		flags := thing.(map[string]interface{})
		for k, v := range flags {
			option := k
			onoff := false
			switch v.(type) {
			case int:
				onoff = v.(int) == 1
			case int64:
				onoff = v.(int64) == 1
			case bool:
				onoff = v.(bool)
			case string:
				onoff = strings.ToLower(v.(string)) == "true"
			}
			if onoff {
				switch option {
				case "SkipHTML":
					flag |= html.SkipHTML
				case "SkipImages":
					flag |= html.SkipImages
				case "SkipLinks":
					flag |= html.SkipLinks
				case "Safelink":
					flag |= html.Safelink
				case "NofollowLinks":
					flag |= html.NofollowLinks
				case "NoreferrerLinks":
					flag |= html.NoreferrerLinks
				case "HrefTargetBlank":
					flag |= html.HrefTargetBlank
				case "CompletePage":
					flag |= html.CompletePage
				case "UseXHTML":
					flag |= html.UseXHTML
				case "FootnoteReturnLinks":
					flag |= html.FootnoteReturnLinks
				case "FootnoteNoHRTag":
					flag |= html.FootnoteNoHRTag
				case "Smartypants":
					flag |= html.Smartypants
				case "SmartypantsFractions":
					flag |= html.SmartypantsFractions
				case "SmartypantsDashes":
					flag |= html.SmartypantsDashes
				case "SmartypantsLatexDashes":
					flag |= html.SmartypantsLatexDashes
				case "SmartypantsAngledQuotes":
					flag |= html.SmartypantsAngledQuotes
				case "SmartypantsQuotesNBSP":
					flag |= html.SmartypantsQuotesNBSP
				case "TOC":
					flag |= html.TOC
				case "CommonFlags":
					flag |= html.CommonFlags
				}
			}
		}
	}
	return flag
}

// ConfigMmark tames a map[string]interface{} and updates
// the configuration returning parser.Extensions and html.Flag
// based on the processed map. NOTE: Settings for markdown are
// objects under m["mmark"]. This lets you pass your whole
// app config map and still have control of the markdown bit
// independently.
func ConfigMmark(config map[string]interface{}) (parser.Extensions, html.Flags, error) {
	ext := mmarkExtensions(config)
	// Set markdown to Mmark
	ext |= parser.Mmark
	htmlFlags := mmarkRenderOptions(config)
	return ext, htmlFlags, nil
}

// mmarkProcessor runs gomarkdown engine using the Mmark extentions and an HTML renderer setup
func mmarkProcessor(fName string, input []byte) ([]byte, error) {
	configType, frontMatterSrc, mmarkSrc := SplitFrontMatter(input)
	//NOTE: should inspect front matter for Markdown parsing configuration
	config, err := ProcessorConfig(configType, frontMatterSrc)
	if err != nil {
		return nil, err
	}
	ext, htmlFlags, err := ConfigMmark(config)
	// Used when processing XML versus man or HTML output.
	parserFlags := parser.FlagsNone
	p := parser.NewWithExtensions(mparser.Extensions | ext)
	init := mparser.NewInitial(fName)
	documentTitle := "" // hack to get document title from TOML title block and then set it here.
	p.Opts = parser.Options{
		ParserHook: func(data []byte) (ast.Node, []byte, int) {
			node, data, consumed := mparser.Hook(data)
			if t, ok := node.(*mast.Title); ok {
				if !t.IsTriggerDash() {
					documentTitle = t.TitleData.Title
				}
			}
			return node, data, consumed
		},
		ReadIncludeFn: init.ReadInclude,
		Flags:         parserFlags,
	}

	doc := markdown.Parse(mmarkSrc, p)
	mparser.AddBibliography(doc)
	mparser.AddIndex(doc)

	documentLanguage := ""
	mhtmlOpts := mhtml.RendererOptions{
		Language: lang.New(documentLanguage),
	}

	opts := html.RendererOptions{
		Comments:       [][]byte{[]byte("//"), []byte("#")}, // used for callouts.
		RenderNodeHook: mhtmlOpts.RenderHook,
		Flags:          htmlFlags, //html.CommonFlags | html.FootnoteNoHRTag | html.FootnoteReturnLinks | html.CompletePage,
		Generator:      `  <meta name="GENERATOR" content="github.com/caltechylibrary/mkpage using Mmark/gomarkdown processor`,
	}
	opts.Title = documentTitle // hack to add-in discovered title

	renderer := html.NewRenderer(opts)
	out := markdown.Render(doc, renderer)
	return out, nil
}

// gomarkdownExtensions takes a config (map[string]interface{}) and
// returns the ORed extenions flags for map["gomarkdown"].
func gomarkdownExtensions(config map[string]interface{}) parser.Extensions {
	ext := parser.NoExtensions
	if thing, ok := config["gomarkdown"]; ok == true {
		extensions := thing.(map[string]interface{})
		for k, v := range extensions {
			onoff := false
			option := k
			switch v.(type) {
			case int:
				onoff = v.(int) == 1
			case int64:
				onoff = v.(int64) == 1
			case bool:
				onoff = v.(bool)
			case string:
				onoff = strings.ToLower(v.(string)) == "true"
			}
			if onoff {
				switch option {
				case "NoIntraEmphasis":
					ext |= parser.NoIntraEmphasis
				case "Tables":
					ext |= parser.Tables
				case "FencedCode":
					ext |= parser.FencedCode
				case "Autolink":
					ext |= parser.Autolink
				case "Strikethrough":
					ext |= parser.Strikethrough
				case "LaxHTMLBlocks":
					ext |= parser.LaxHTMLBlocks
				case "SpaceHeadings":
					ext |= parser.SpaceHeadings
				case "HardLineBreak":
					ext |= parser.HardLineBreak
				case "TabSizeEight":
					ext |= parser.TabSizeEight
				case "Footnotes":
					ext |= parser.Footnotes
				case "NoEmptyLineBeforeBlock":
					ext |= parser.NoEmptyLineBeforeBlock
				case "HeadingIDs":
					ext |= parser.HeadingIDs
				case "Titleblock":
					ext |= parser.Titleblock
				case "AutoHeadingIDs":
					ext |= parser.AutoHeadingIDs
				case "BackslashLineBreak":
					ext |= parser.BackslashLineBreak
				case "DefinitionLists":
					ext |= parser.DefinitionLists
				case "MathJax":
					ext |= parser.MathJax
				case "OrderedListStart":
					ext |= parser.OrderedListStart
				case "Attributes":
					ext |= parser.Attributes
				case "SuperSubscript":
					ext |= parser.SuperSubscript
				case "CommonExtensions":
					ext |= parser.CommonExtensions
				}
			}
		}
	}
	return ext
}

// gomarkdownRenderOptions config (map[string]interface{}) and
// returns the ORed extenions flags for map["gomarkdown"].
func gomarkdownRenderOptions(config map[string]interface{}) html.Flags {
	flag := html.FlagsNone
	if thing, ok := config["gomarkdown"]; ok == true {
		flags := thing.(map[string]interface{})
		for k, v := range flags {
			option := k
			onoff := false
			switch v.(type) {
			case int:
				onoff = v.(int) == 1
			case int64:
				onoff = v.(int64) == 1
			case bool:
				onoff = v.(bool)
			case string:
				onoff = strings.ToLower(v.(string)) == "true"
			}
			if onoff {
				switch option {
				case "SkipHTML":
					flag |= html.SkipHTML
				case "SkipImages":
					flag |= html.SkipImages
				case "SkipLinks":
					flag |= html.SkipLinks
				case "Safelink":
					flag |= html.Safelink
				case "NofollowLinks":
					flag |= html.NofollowLinks
				case "NoreferrerLinks":
					flag |= html.NoreferrerLinks
				case "HrefTargetBlank":
					flag |= html.HrefTargetBlank
				case "CompletePage":
					flag |= html.CompletePage
				case "UseXHTML":
					flag |= html.UseXHTML
				case "FootnoteReturnLinks":
					flag |= html.FootnoteReturnLinks
				case "FootnoteNoHRTag":
					flag |= html.FootnoteNoHRTag
				case "Smartypants":
					flag |= html.Smartypants
				case "SmartypantsFractions":
					flag |= html.SmartypantsFractions
				case "SmartypantsDashes":
					flag |= html.SmartypantsDashes
				case "SmartypantsLatexDashes":
					flag |= html.SmartypantsLatexDashes
				case "SmartypantsAngledQuotes":
					flag |= html.SmartypantsAngledQuotes
				case "SmartypantsQuotesNBSP":
					flag |= html.SmartypantsQuotesNBSP
				case "TOC":
					flag |= html.TOC
				case "CommonFlags":
					flag |= html.CommonFlags
				}
			}
		}
	}
	return flag
}

// ConfigMarkdown takes a map[string]interface{} and updates
// the configuration returning parser.Extensions and html.Flag
// based on the processed map. NOTE: Settings for markdown are
// objects under m["gomarkdown"]. This lets you pass your whole
// app config map and still have control of the markdown bit
// independently.
func ConfigMarkdown(config map[string]interface{}) (parser.Extensions, html.Flags, error) {
	ext := gomarkdownExtensions(config)
	htmlFlags := gomarkdownRenderOptions(config)
	return ext, htmlFlags, nil
}

// markdownProcessor applies Markdown processing with overrides
// for gomarkdown, mmark and fountain in the
// document's frontmatter.  You need to supply a
// markdown processor as a func to envoke the preferred default.
func markdownProcessor(input []byte, defaultProcessor func([]byte) ([]byte, error)) ([]byte, error) {
	input = normalizeEOL(input)
	configType, frontMatterSrc, mdSrc := SplitFrontMatter(input)
	//NOTE: should inspect front matter for Markdown parsing configuration
	config, err := ProcessorConfig(configType, frontMatterSrc)
	if err != nil {
		return nil, err
	}

	if thing, ok := config["markup"]; ok == true {
		markup := thing.(string)
		switch markup {
		case "mmark":
			ext, htmlFlags, err := ConfigMmark(config)
			if err != nil {
				return nil, err
			}
			p := parser.NewWithExtensions(mparser.Extensions | ext)
			parserFlags := parser.FlagsNone
			// We're working as if we're getting data from stdin here ...
			init := mparser.NewInitial("")
			documentTitle := "" // hack to get document title from TOML title block and then set it here.
			p.Opts = parser.Options{
				ParserHook: func(data []byte) (ast.Node, []byte, int) {
					node, data, consumed := mparser.Hook(data)
					if t, ok := node.(*mast.Title); ok {
						if !t.IsTriggerDash() {
							documentTitle = t.TitleData.Title
						}
					}
					return node, data, consumed
				},
				ReadIncludeFn: init.ReadInclude,
				Flags:         parserFlags,
			}
			doc := markdown.Parse(mdSrc, p)
			mparser.AddBibliography(doc)
			mparser.AddIndex(doc)

			documentLanguage := ""
			mhtmlOpts := mhtml.RendererOptions{
				Language: lang.New(documentLanguage),
			}

			opts := html.RendererOptions{
				Comments:       [][]byte{[]byte("//"), []byte("#")}, // used for callouts.
				RenderNodeHook: mhtmlOpts.RenderHook,
				Flags:          htmlFlags,
				Generator:      `  <meta name="GENERATOR" content="github.com/mmarkdown/mmark Mmark Markdown Processor - mmark.nl`,
			}
			opts.Title = documentTitle // hack to add-in discovered title
			renderer := html.NewRenderer(opts)
			return markdown.Render(doc, renderer), nil
		case "gomarkdown":
			ext, htmlFlags, err := ConfigMarkdown(config)
			if err != nil {
				return nil, err
			}

			p := parser.NewWithExtensions(ext)
			opts := html.RendererOptions{Flags: htmlFlags}
			r := html.NewRenderer(opts)
			return markdown.ToHTML(mdSrc, p, r), nil
		case "fountain":
			if err := ConfigFountain(config); err != nil {
				return nil, err
			}
			src, err := fountain.Run(mdSrc)
			if err != nil {
				return nil, err
			}
			return src, nil
		default:
			return nil, fmt.Errorf("unknown markup engine")
		}
	}
	return defaultProcessor(mdSrc)
}

// gomarkdownProcessor wraps gomarkdown with overrides for
// mmark, and fountain handling the splitting off the front matter if
// present and configuration via front matter.
func gomarkdownProcessor(input []byte) ([]byte, error) {
	input = normalizeEOL(input)
	return markdownProcessor(input, func(input []byte) ([]byte, error) {
		// Default to gomarkdown markdown processor
		// with CommonExtensions and CommonHTMLFlags
		p := parser.NewWithExtensions(parser.CommonExtensions)
		opts := html.RendererOptions{Flags: html.CommonFlags}
		r := html.NewRenderer(opts)
		return markdown.ToHTML(input, p, r), nil
	})
}

// ConfigFountain sets the fountain defaults then applies
// the map[string]interface{} overwriting the defaults
// returns error necessary.
func ConfigFountain(config map[string]interface{}) error {
	if thing, ok := config["fountain"]; ok == true {
		cfg := thing.(map[string]interface{})
		for k, v := range cfg {
			switch v.(type) {
			case bool:
				onoff := v.(bool)
				switch k {
				case "AsHTMLPage":
					fountain.AsHTMLPage = onoff
				case "InlineCSS":
					fountain.InlineCSS = onoff
				case "LinkCSS":
					fountain.LinkCSS = onoff
				}
			case string:
				if k == "IncludeCSS" {
					fountain.CSS = v.(string)
				}
			default:
				return fmt.Errorf("Unknown fountain option %q", k)
			}
		}
	}
	return nil
}

// fountainProcessor wraps fountain.Run() splitting off the front
// matter if present.
func fountainProcessor(input []byte) ([]byte, error) {
	var err error

	configType, frontMatterSrc, fountainSrc := SplitFrontMatter(input)
	config, err := ProcessorConfig(configType, frontMatterSrc)
	if err != nil {
		return nil, err
	}
	if err := ConfigFountain(config); err != nil {
		return nil, err
	}
	src, err := fountain.Run(fountainSrc)
	if err != nil {
		return nil, err
	}
	return src, nil
}

// ResolveData takes a data map and reads in the files and URL sources
// as needed turning the data into strings to be applied to the template.
func ResolveData(data map[string]string) (map[string]interface{}, error) {
	var (
		out map[string]interface{}
	)

	isContentType := func(vals []string, target string) bool {
		for _, h := range vals {
			if strings.Contains(h, target) == true {
				return true
			}
		}
		return false
	}

	out = make(map[string]interface{})
	for key, val := range data {
		switch {
		case strings.HasPrefix(val, TextPrefix) == true:
			out[key] = strings.TrimPrefix(val, TextPrefix)
		case strings.HasPrefix(val, GomarkdownPrefix) == true:
			src, err := gomarkdownProcessor([]byte(strings.TrimPrefix(val, GomarkdownPrefix)))
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, MarkdownPrefix) == true:
			//NOTE: We're using Gomarkdown as our default processor
			src, err := gomarkdownProcessor([]byte(strings.TrimPrefix(val, MarkdownPrefix)))
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, FountainPrefix) == true:
			src, err := fountainProcessor([]byte(strings.TrimPrefix(val, FountainPrefix)))
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, JSONPrefix) == true:
			var o interface{}
			err := json.Unmarshal(bytes.TrimPrefix([]byte(val), []byte(JSONPrefix)), &o)
			if err != nil {
				return out, fmt.Errorf("Can't JSON decode (%s) %s, %s", key, val, err)
			}
			out[key] = o
		case strings.HasPrefix(val, "http://") == true || strings.HasPrefix(val, "https://") == true:
			resp, err := http.Get(val)
			if err != nil {
				return out, fmt.Errorf("Error from (%s) %s, %s", key, val, err)
			}
			defer resp.Body.Close()
			if resp.StatusCode == 200 {
				buf, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return out, err
				}
				if contentTypes, ok := resp.Header["Content-Type"]; ok == true {
					switch {
					case isContentType(contentTypes, "application/json") == true:
						var o interface{}
						err := json.Unmarshal(buf, &o)
						if err != nil {
							return out, fmt.Errorf("Can't JSON decode (%s) %s, %s", key, val, err)
						}
						out[key] = o
					case isContentType(contentTypes, "text/markdown") == true:
						src, err := gomarkdownProcessor(buf)
						if err != nil {
							return nil, err
						}
						out[key] = fmt.Sprintf("%s", src)
					case isContentType(contentTypes, "text/fountain") == true:
						src, err := fountainProcessor(buf)
						if err != nil {
							return nil, err
						}
						out[key] = fmt.Sprintf("%s", src)
					default:
						out[key] = string(buf)
					}
				} else {
					out[key] = string(buf)
				}
			}
		default:
			buf, err := ioutil.ReadFile(val)
			if err != nil {
				return out, fmt.Errorf("Can't read (%s) %q, %s", key, val, err)
			}
			ext := path.Ext(val)
			switch {
			case strings.Compare(ext, ".fountain") == 0 ||
				strings.Compare(ext, ".spmd") == 0:
				src, err := fountainProcessor(buf)
				if err != nil {
					return nil, err
				}
				out[key] = fmt.Sprintf("%s", src)
			case strings.Compare(ext, ".md") == 0:
				src, err := gomarkdownProcessor(buf)
				if err != nil {
					return nil, err
				}
				out[key] = fmt.Sprintf("%s", src)
			case strings.Compare(ext, ".mmark") == 0:
				src, err := mmarkProcessor(val, buf)
				if err != nil {
					return nil, err
				}
				out[key] = fmt.Sprintf("%s", src)
			case strings.Compare(ext, ".json") == 0:
				var o interface{}
				err := json.Unmarshal(buf, &o)
				if err != nil {
					return out, fmt.Errorf("Can't JSON decode (%s) %s, %s", key, val, err)
				}
				out[key] = o
			default:
				out[key] = string(buf)
			}
		}
	}
	return out, nil
}

// MakePage applies the key/value map to the named template in tmpl and renders to writer and returns an error if something goes wrong
func MakePage(wr io.Writer, templateName string, tmpl *template.Template, keyValues map[string]string) error {
	data, err := ResolveData(keyValues)
	if err != nil {
		return fmt.Errorf("Can't resolve data source %s", err)
	}
	return tmpl.ExecuteTemplate(wr, templateName, data)
}

// MakePageString applies the key/value map to the named template tmpl and renders the results to a string and error if something goes wrong
func MakePageString(templateName string, tmpl *template.Template, keyValues map[string]string) (string, error) {
	var buf bytes.Buffer
	wr := io.Writer(&buf)
	err := MakePage(wr, templateName, tmpl, keyValues)
	return buf.String(), err
}

// RelativeDocPath calculate the relative path from source to target based on
// implied common base.
//
// Example:
//
//     docPath := "docs/chapter-01/lesson-02.html"
//     cssPath := "css/site.css"
//     fmt.Printf("<link href=%q>\n", MakeRelativePath(docPath, cssPath))
//
// Output:
//
//     <link href="../../css/site.css">
//
func RelativeDocPath(source, target string) string {
	var result []string

	sep := string(os.PathSeparator)
	dname, _ := path.Split(source)
	for i := 0; i < strings.Count(dname, sep); i++ {
		result = append(result, "..")
	}
	result = append(result, target)
	p := strings.Join(result, sep)
	if strings.HasSuffix(p, "/.") {
		return strings.TrimSuffix(p, ".")
	}
	return p
}

//
// Below is addition code to support mkslides
//

// Slide is the metadata about a slide to be generated.
type Slide struct {
	CurNo   int    `json:"cur_no,omitemtpy"`
	PrevNo  int    `json:"prev_no,omitempty"`
	NextNo  int    `json:"next_no,omitempty"`
	FirstNo int    `json:"first_no,omitempty"`
	LastNo  int    `json:"last_no,omitempty"`
	FName   string `json:"filename,omitempty"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	CSSPath string `json:"csspath,omitempty"`
	JSPath  string `json:"jspath,omitempty"`
	CSS     string `json:"css,omitempty"`
	Header  string `json:"header,omitempty"`
	Footer  string `json:"footer,omitempty"`
	Nav     string `json:"nav,omitempty"`
}

// MarkdownToSlides turns a markdown file into one or more Slide structs
// Which populate predefined key/value pairs for later rendering in Markdown
func MarkdownToSlides(fname string, mdSource []byte) ([]*Slide, error) {
	var slides []*Slide

	// Note: handle legacy CR/LF endings as well as normal LF line endings
	if bytes.Contains(mdSource, []byte("\r\n")) {
		mdSource = bytes.Replace(mdSource, []byte("\r\n"), []byte("\n"), -1)
	}
	// Note: We're only spliting on whole line that contain "---", not lines ending with where text might end with "---"
	mdSlides := bytes.Split(mdSource, []byte("\n---\n"))

	lastSlide := len(mdSlides) - 1
	for i, s := range mdSlides {
		src, err := gomarkdownProcessor(s)
		if err != nil {
			return nil, fmt.Errorf("%s slide %d error, %s", fname, i+1, err)
		}
		slides = append(slides, &Slide{
			FName:   strings.TrimSuffix(path.Base(fname), path.Ext(fname)),
			CurNo:   i,
			PrevNo:  (i - 1),
			NextNo:  (i + 1),
			FirstNo: 0,
			LastNo:  lastSlide,
			Content: fmt.Sprintf("%s", src),
		})
	}
	return slides, nil
}

// MakeSlide this takes a io.Writer, a template, key/value map pairs and Slide struct.
// It resolves the data int key/value pairs, merges the prefined mapping from Slide struct
// then executes the template.
func MakeSlide(wr io.Writer, templateName string, tmpl *template.Template, keyValues map[string]string, slide *Slide) error {
	data, err := ResolveData(keyValues)
	if err != nil {
		return fmt.Errorf("Can't resolve data source %s", err)
	}
	// Merge the slide metadata into data pairs for template
	data["filename"] = slide.FName
	data["cur_no"] = slide.CurNo
	data["prev_no"] = slide.PrevNo
	data["next_no"] = slide.NextNo
	data["first_no"] = slide.FirstNo
	data["last_no"] = slide.LastNo
	data["content"] = slide.Content
	data["header"] = slide.Header
	data["footer"] = slide.Header
	data["nav"] = slide.Nav
	return tmpl.ExecuteTemplate(wr, templateName, data)
}

// MakeSlideFile this takes a template and slide and renders the results to a file.
func MakeSlideFile(templateName string, tmpl *template.Template, keyValues map[string]string, slide *Slide) error {
	sname := fmt.Sprintf(`%02d-%s.html`, slide.CurNo, strings.TrimSuffix(path.Base(slide.FName), path.Ext(slide.FName)))
	fp, err := os.Create(sname)
	if err != nil {
		return fmt.Errorf("%s %s", sname, err)
	}
	defer fp.Close()
	err = MakeSlide(fp, templateName, tmpl, keyValues, slide)
	if err != nil {
		return fmt.Errorf("%s %s", sname, err)
	}
	return nil
}

// MakeSlideString this takes a template and slide and renders the results to a string
func MakeSlideString(templateName string, tmpl *template.Template, keyValues map[string]string, slide *Slide) (string, error) {
	var buf bytes.Buffer
	wr := io.Writer(&buf)
	err := MakeSlide(wr, templateName, tmpl, keyValues, slide)
	return buf.String(), err
}

// NormalizeDate takes a MySQL like date string and returns a time.Time or error
func NormalizeDate(s string) (time.Time, error) {
	switch len(s) {
	case len(`2006-01-02 15:04:05 -0700`):
		dt, err := time.Parse(`2006-01-02 15:04:05 -0700`, s)
		return dt, err
	case len(`2006-01-02 15:04:05`):
		dt, err := time.Parse(`2006-01-02 15:04:05`, s)
		return dt, err
	case len(`2006-01-02`):
		dt, err := time.Parse(`2006-01-02`, s)
		return dt, err
	default:
		return time.Time{}, fmt.Errorf("Can't format %s, expected format like 2006-01-02 15:04:05 -0700", s)
	}
}

// Walk takes a start path and walks the file system to process Markdown files for useful elements.
func Walk(startPath string, filterFn func(p string, info os.FileInfo) bool, outputFn func(s string, info os.FileInfo) error) error {
	err := filepath.Walk(startPath, func(p string, info os.FileInfo, err error) error {
		// Are we interested in this path?
		if filterFn(p, info) == true {
			// Yes, so send to output function.
			if err := outputFn(p, info); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// Grep looks for the first line matching the expression
// in src.
func Grep(exp string, src string) string {
	re, err := regexp.Compile(exp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%q is not a valid, %s\n", exp, err)
		return ""
	}
	lines := strings.Split(src, "\n")
	for _, line := range lines {
		s := re.FindString(line)
		if len(s) > 0 {
			return s
		}
	}
	return ""
}

func init() {
	if bString, ok := Defaults["/templates/page.tmpl"]; ok == true {
		DefaultTemplateSource = fmt.Sprintf("%s", bString)
	}
	if bString, ok := Defaults["/templates/slides.tmpl"]; ok == true {
		DefaultSlideTemplateSource = fmt.Sprintf("%s", bString)
	}
}
