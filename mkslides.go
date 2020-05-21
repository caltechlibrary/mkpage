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
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	// Caltech Library packages
	"github.com/caltechlibrary/tmplfn"
)

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

const (
	// Note: We're only spliting on whole line that contain "--", not lines ending with where text might end with "--"
	slideDelimiter = "\n--\n"
)

var (
	// DefaultSlideTemplateSource provides the default HTML template for mkslides package,
	// you probably want to override this... is defined in init by Defaults["/templates/slides.tmpl"]
	// Defaults is a map to assets defined in assets.go which is build with pkgasset and
	// the contents of the defaults folder in this repository.
	DefaultSlideTemplateSource string
)

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

//
// Below is addition code to support mkslides
//

// MarkdownToSlides turns a markdown file into one or more Slide structs
// Which populate predefined key/value pairs for later rendering in Markdown
func MarkdownToSlides(fname string, mdSource []byte) ([]*Slide, error) {
	var slides []*Slide

	// Note: handle legacy CR/LF endings as well as normal LF line endings
	if bytes.Contains(mdSource, []byte("\r\n")) {
		mdSource = bytes.Replace(mdSource, []byte("\r\n"), []byte("\n"), -1)
	}
	mdSlides := bytes.Split(mdSource, []byte(slideDelimiter))

	lastSlide := len(mdSlides) - 1
	for i, s := range mdSlides {
		src, err := pandocProcessor(s, "markdown", "html")
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

// MakeSlideSet is the complete mkslides engine based on the gomarkdown
// Markdown engine and Go templates.
func MakeSlideSet(inputFName string, templateSources []string, defaultTemplate []byte, keyValues map[string]string) error {
	var (
		templateName string
		err          error
	)
	src, err := ioutil.ReadFile(inputFName)
	if err != nil {
		return err
	}

	// Build the slides
	slides, err := MarkdownToSlides(inputFName, src)
	if err != nil {
		return err
	}

	// Default Template Name is slides.tmpl
	if len(templateSources) == 0 {
		templateSources = []string{"slides.tmpl"}
	}

	// Create our Tmpl with its function map
	tmpl := tmplfn.New(tmplfn.AllFuncs())

	// Load any user supplied templates
	if len(templateSources) > 0 {
		err = tmpl.ReadFiles(templateSources...)
		if err != nil {
			return err
		}
		templateName = templateSources[0]
	} else {
		// Load our default template maps
		err = tmpl.Add(templateName, defaultTemplate)
		if err != nil {
			return err
		}
	}

	// Assemble our templates
	t, err := tmpl.Assemble()
	if err != nil {
		return err
	}

	// Render the slides
	for _, slide := range slides {
		err = MakeSlideFile(templateName, t, keyValues, slide)
		if err == nil {
			// Note: Give some feed back when slide written successful
			fmt.Printf("Wrote %02d-%s.html\n", slide.CurNo, strings.TrimSuffix(path.Base(slide.FName), path.Ext(slide.FName)))
		} else {
			// Note: return an error if we have a problem
			return err
		}
	}
	return nil
}
