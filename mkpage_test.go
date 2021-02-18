//
// mkpage is a thought experiment in a light weight template and markdown processor.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2021, Caltech
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
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"testing"
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
	buf, err := pandocProcessor(src, "markdown", "html")
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
	buf, err = pandocProcessor(src, "markdown", "html")
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
Hello ${hello}

Nav: ${nav}

Content: ${content}

Weather: ${weather.data.text}
`

	keyValues := map[string]string{
		"hello":   "text:Hi there!",
		"nav":     path.Join("testdata", "nav.md"),
		"content": path.Join("testdata", "content.md"),
		"weather": "http://forecast.weather.gov/MapClick.php?lat=13.4712&lon=144.7496&FcstType=json",
	}

	out, err := MakePandocString(src, keyValues)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	checkForString(out, "Hi there!")
	checkForString(out, "<ul>")
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
	src1, err := pandocProcessor(srcLF, "markdown", "html")
	if err != nil {
		t.Errorf("pandocProcessor(srcLF, 'markdown', 'html') error %s", err)
	}
	// Render HTML using normalize old DOS eol
	src2, err := pandocProcessor(srcCRLF, "markdown", "html")
	if err != nil {
		t.Errorf("pandocProcessor(srcCRLF, 'markdown', 'html') error %s", err)
	}
	if bytes.Compare(src1, src2) != 0 {
		t.Errorf("expected (eol normalized) ->\n%s\ngot ->\n%s\n",
			src1, src2)
	}
}

func TestScanArgs2(t *testing.T) {
	expectedGenerator := `one`
	expectedParams := []string{`two`, `three`}
	src := `one two three`
	generator, params := scanArgs(src)

	if generator != expectedGenerator {
		t.Errorf("expected %q, got %q from %q", expectedGenerator, generator, src)
	}
	if len(params) != len(expectedParams) {
		t.Errorf("expected %d, got %d from %+v", len(expectedParams), len(params), params)
		t.FailNow()
	}
	for i, val := range expectedParams {
		if val != params[i] {
			t.Errorf("expected param(%d) %q, got %q from %+v", i, val, params[i], params)
		}
	}

	expectedGenerator = `python3`
	expectedParams = []string{`eprints_html_view.py`, `3`, `et el.`}
	src = `python3 eprints_html_view.py 3 "et el."`
	generator, params = scanArgs(src)
	if generator != expectedGenerator {
		t.Errorf("expected %q, got %q from %q", expectedGenerator, generator, src)
	}
	if len(params) != len(expectedParams) {
		t.Errorf("expected %d, got %d from %+v", len(expectedParams), len(params), params)
		t.FailNow()
	}
	for i, val := range expectedParams {
		if val != params[i] {
			t.Errorf("expected param(%d) %q, got %q from %+v", i, val, params[i], params)
		}
	}
	t.Errorf("DEBUG param[2] -> %q", params[2])
}
