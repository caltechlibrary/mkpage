//
// mkpage is a thought experiment in a light weight template and markdown processor.
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
package mkpage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"testing"
	"text/template"
)

func TestResolveData(t *testing.T) {
	checkMap := func(ky string, expected string, m map[string]interface{}) error {
		if val, ok := m[ky]; ok == true {
			switch vv := val.(type) {
			case string:
				s := fmt.Sprintf("%s", val)
				if strings.Compare(expected, s) == 0 {
					return nil
				}
				return fmt.Errorf("expected %q, found %q, %d", expected, s, strings.Compare(expected, s))
			default:
				return fmt.Errorf("expected %s, found type %b, %s", expected, vv, val)
			}
		} else {
			return fmt.Errorf("expected %s, missing %s", expected, ky)
		}
	}

	keyValues := map[string]string{
		"Hello":   "text:Hi there!",
		"Hi":      "markdown:*Hi there!*",
		"Nav":     path.Join("testdata", "nav.md"),
		"Content": path.Join("testdata", "content.md"),
		"Weather": "http://forecast.weather.gov/MapClick.php?lat=13.4712&lon=144.7496&FcstType=json",
	}
	data, err := ResolveData(keyValues)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := checkMap("Hello", "Hi there!", data); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := checkMap("Hi", "<p><em>Hi there!</em></p>\n", data); err != nil {
		t.Error(err)
		t.FailNow()
	}

	src, err := ioutil.ReadFile(path.Join("testdata", "nav.md"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	buf, err := gomarkdownProcessor(src)
	expected := string(buf)

	if err := checkMap("Nav", expected, data); err != nil {
		t.Error(err)
		t.FailNow()
	}

	src, err = ioutil.ReadFile(path.Join("testdata", "content.md"))
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	buf, err = gomarkdownProcessor(src)
	expected = string(buf)

	if err := checkMap("Content", expected, data); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if _, ok := data["Weather"]; ok == false {
		t.Error("Expected a JSON blob for weather")
		t.FailNow()
	}
}

func TestMakePage(t *testing.T) {
	checkForString := func(src, target string) bool {
		if strings.Contains(src, target) == false {
			t.Errorf("expected %q in %s", target, src)
			return false
		}
		return true
	}

	src := `
{{define "Hello"}}
Hello {{.hello}}

Nav: {{.nav}}

Content: {{.content}}

Weather: {{.weather.data.text}}
{{end}}
`

	keyValues := map[string]string{
		"hello":   "text:Hi there!",
		"nav":     path.Join("testdata", "nav.md"),
		"content": path.Join("testdata", "content.md"),
		"weather": "http://forecast.weather.gov/MapClick.php?lat=13.4712&lon=144.7496&FcstType=json",
	}

	tmpl := template.Must(template.New("test.tmpl").Parse(src))
	out, err := MakePageString("Hello", tmpl, keyValues)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	checkForString(out, "Hi there!")
	checkForString(out, "<ul>")
}

func TestBasic(t *testing.T) {
	src := `

# Opening Slide

## This is an introduction

+ me, <me@example.org>
+ [mygitrepo.org](http://mygitrepo.org)
+ February 31, 2016, Somebody's University, Someplace

---

# Slide Two

| col A | col B | Col C |
|-------|-------|-------|
| A     | One   | 1     |

---

## Slide Three

This is slide three, just a random paragraph of text. Blah, blah, blah, blah, blah, blah, boink!


`

	tmpl, err := template.New("slides.tmpl").Parse(DefaultSlideTemplateSource)
	if err != nil {
		t.Errorf("Can't parse DefaultTemplateSource templates %s", err)
		t.FailNow()
	}

	titles := []string{
		"<h1>Opening Slide</h1>",
		"<h1>Slide Two</h1>",
		"<h2>Slide Three</h2>",
	}

	slides, err := MarkdownToSlides("test.html", []byte(src))
	if err != nil {
		t.Errorf("MarkdownToSlides(%q, ...) error %s", "test.html", err)
		t.FailNow()
	}
	if len(slides) != 3 {
		src, _ := json.Marshal(slides)
		t.Errorf("Was expected three slides %+v\n%s", slides, src)
	}

	keyVals := map[string]string{}
	for i, slide := range slides {
		keyVals["Title"] = "text:" + titles[i]
		s, err := MakeSlideString("slides.tmpl", tmpl, keyVals, slide)
		if err != nil {
			t.Errorf("MakeSlideString() failed %d - %s", i, err)
		}
		if strings.Contains(s, titles[i]) == false {
			t.Errorf("Expected %q in slide %d -> %s", titles[i], i, s)
		}
	}
}

func TestGrep(t *testing.T) {
	src := `
# This is some article

by Jane Doe 3001-01-01

It was some New Years day...
`
	expected := `by Jane Doe 3001-01-01`
	result := Grep(BylineExp, src)
	if expected != result {
		t.Errorf("expected %q, got %q for Grep() for byline", expected, result)
	}
	expected = `# This is some article`
	result = Grep(TitleExp, src)
	if expected != result {
		t.Errorf("expected %q, got %q for Grep() for title", expected, result)
	}
}

func TestCRLFHandling(t *testing.T) {
	srcRaw := []byte(`

# Title 
 
*Italics* 
 
This is text **in bold** in text. 
 
## List: 
 
-   Item one. 
-   Item two. 
-   Item three. 
 
## Topics: 
 
- Item one. 
- Item two. 
- Item three. 
 
## Spacing: 
 
-Item one. 
-Item two. 
-Item three. 
 
## Plus: 
 
+ Item one. 
+ Item two. 
+ Item three. 
 
## Star: 
 
* Item one. 
* Item two. 
* Item three. 
 
## Indent: 
 
* Item one. 
    * Item two. 
* Item three. 
 
## Numbered 
 
1. One 
1. Two 
1. Three 
 
`)
	srcCRLF := bytes.Replace(srcRaw, []byte("\n"), []byte("\r\n"), -1)
	srcLF := srcRaw
	// Render HTML using normalize Unix eol
	src1, err := gomarkdownProcessor(srcLF)
	if err != nil {
		t.Errorf("gomarkdownProcessor(srcCRLF) error %s", err)
	}
	// Render HTML using normalize old DOS eol
	src2, err := gomarkdownProcessor(srcCRLF)
	if err != nil {
		t.Errorf("gomarkdownProcessor(srcLF) error %s", err)
	}
	if bytes.Compare(src1, src2) != 0 {
		t.Errorf("expected (eol normalized) ->\n%s\ngot ->\n%s\n",
			src1, src2)
	}
}
