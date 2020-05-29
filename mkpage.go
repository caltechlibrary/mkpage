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
	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v3"

	// Fountain support for scripts, interviews and narration
	"github.com/rsdoiel/fountain"
)

const (
	// Prefix for explicit string types

	// JSONPrefix designates a string as JSON formatted content
	JSONPrefix = "json:"
	// MarkdownPrefix designates a string as Markdown (common mark) content
	// to be parsed by pandoc
	MarkdownPrefix = "markdown:"
	// MMarkPrefix designates MMark format, for now this will just be passed to pandoc.
	MMarkPrefix = "mmark:"
	// TextPrefix designates a string as text/plain not needed processing
	TextPrefix = "text:"
	// FountainPrefix designates a string as Fountain formatted content
	FountainPrefix = "fountain:"
	// TextilePrefix designates source as Textile for processing by pandoc.
	TextilePrefix = "textile:"
	// ReStructureText designates source as ReStructureText for processing by pandoc
	ReStructureTextPrefix = "rst:"
	// JiraPrefix markup designates source as Jire text for processing by pandoc
	JiraPrefix = "jira:"

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

func yamlToJson(src []byte) ([]byte, error) {
	m1 := make(map[interface{}]interface{})
	err := yaml.Unmarshal(src, &m1)
	if err != nil {
		return nil, err
	}
	m2 := make(map[string]interface{})
	for key, value := range m1 {
		switch key.(type) {
		case string:
			m2[key.(string)] = value
		case int:
			m2[fmt.Sprintf("%d", key)] = value
		case float64:
			m2[fmt.Sprintf("%f", key)] = value
		default:
			return nil, fmt.Errorf("JSON conversion failed, can't convert %T, %+v to string", key, key)
		}
	}
	return json.MarshalIndent(m2, "", "    ")
}

// SplitFrontMatter takes a []byte input splits it into front matter type,
// front matter source and Markdown source. If either is missing an
// empty []byte is returned for the missing element.
func SplitFrontMatter(input []byte) (int, []byte, []byte) {
	// YAML front matter uses ---, note this conflicts with Mmark practice, do I want to support YAML like this?
	if bytes.HasPrefix(input, []byte("---\n")) {
		parts := bytes.SplitN(bytes.TrimPrefix(input, []byte("---\n")), []byte("\n---\n"), 2)
		if len(parts) > 1 {
			return ConfigIsYAML, parts[0], parts[1]
		}
		if len(parts) > 0 {
			return ConfigIsYAML, parts[0], []byte("")
		}
		return ConfigIsYAML, []byte(""), []byte("")
	}
	// TOML front matter as used in Hugo
	if bytes.HasPrefix(input, []byte("+++\n")) {
		parts := bytes.SplitN(bytes.TrimPrefix(input, []byte("+++\n")), []byte("\n+++\n"), 2)
		if len(parts) > 1 {
			return ConfigIsTOML, parts[0], parts[1]
		}
		if len(parts) > 0 {
			return ConfigIsTOML, parts[0], []byte("")
		}
		return ConfigIsTOML, []byte(""), []byte("")
	}
	// TOML front matter identified in Mmark as three % or dashes,
	// We can support the %, dashes are taken by Hugo style, but
	// maybe I don't want to support that?
	if bytes.HasPrefix(input, []byte("%%%\n")) {
		parts := bytes.SplitN(bytes.TrimPrefix(input, []byte("%%%\n")), []byte("\n%%%\n"), 2)
		if len(parts) > 1 {
			return ConfigIsTOML, parts[0], parts[1]
		}
		if len(parts) > 0 {
			return ConfigIsTOML, parts[0], []byte("")
		}
		return ConfigIsTOML, []byte(""), []byte("")
	}
	// JSON front matter, most Markdown processors.
	if bytes.HasPrefix(input, []byte("{\n")) {
		parts := bytes.SplitN(bytes.TrimPrefix(input, []byte("{\n")), []byte("\n}\n"), 2)
		src := []byte(fmt.Sprintf("{\n%s\n}\n", parts[0]))
		if len(parts) > 1 {
			return ConfigIsJSON, src, parts[1]
		}
		return ConfigIsJSON, src, []byte("")
	}
	// Handle case of no front matter
	return ConfigIsUnknown, []byte(""), input
}

// UnmarshalFrontMatter takes a []byte of front matter source
// and unmarshalls using either JSON, TOML and YAML unmarshalling
// methods.
func UnmarshalFrontMatter(srcType int, src []byte, obj *map[string]interface{}) error {
	switch srcType {
	case ConfigIsTOML:
		// Make sure we have valid Toml
		if err := toml.Unmarshal(src, obj); err != nil {
			return err
		}
	case ConfigIsYAML:
		// With YAML we go through two step conversion
		// YAML to JSON then Unmarshal JSON into our
		// map.
		if jsonSrc, err := yamlToJson(src); err != nil {
			return err
		} else {
			if err := json.Unmarshal(jsonSrc, &obj); err != nil {
				return err
			}
		}
	default:
		// Make sure we have valid JSON
		if err := json.Unmarshal(src, &obj); err != nil {
			return err
		}
	}
	return nil
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
		jsonSrc, err := yamlToJson(frontMatterSrc)
		if err != nil {
			return nil, fmt.Errorf("Can't parse YAML front matter, %s", err)
		}
		if err = json.Unmarshal(jsonSrc, &m); err != nil {
			return nil, err
		}
	case ConfigIsTOML:
		// TOML Front Matter
		if err := toml.Unmarshal(frontMatterSrc, &m); err != nil {
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
		case strings.HasPrefix(val, MMarkPrefix) == true:
			//NOTE: We're using pandoc as our default processor
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, MMarkPrefix)), "markdown", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, MarkdownPrefix) == true:
			//NOTE: We're using pandoc as our default processor
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, MarkdownPrefix)), "markdown", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, JiraPrefix) == true:
			//NOTE: We're using pandoc as our default processor
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, JiraPrefix)), "jira", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, TextilePrefix) == true:
			//NOTE: We're using pandoc as our default processor
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, TextilePrefix)), "textile", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, ReStructureTextPrefix) == true:
			//NOTE: We're using pandoc as our default processor
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, ReStructureTextPrefix)), "rst", "html")
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
				fmType, fmSrc, docSrc := SplitFrontMatter(buf)
				if len(fmSrc) > 0 {
					buf = docSrc
					fmData := map[string]interface{}{}
					if err := UnmarshalFrontMatter(fmType, fmSrc, &fmData); err != nil {
						return out, fmt.Errorf("Can't process front matter (%s), %q, %q", key, val, err)
					}
					// Update, Overwrite `out` with front matter values
					for k, v := range fmData {
						out[k] = v
					}
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
						src, err := pandocProcessor(buf, "markdown", "html")
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
			fmType, fmSrc, docSrc := SplitFrontMatter(buf)
			if len(fmSrc) > 0 {
				buf = docSrc
				fmData := map[string]interface{}{}
				if err := UnmarshalFrontMatter(fmType, fmSrc, &fmData); err != nil {
					return out, fmt.Errorf("Can't process front matter (%s), %q, %q", key, val, err)
				}
				// Update, Overwrite `out` with front matter values
				for k, v := range fmData {
					out[k] = v
				}
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
				src, err := pandocProcessor(buf, "markdown", "html")
				if err != nil {
					return nil, err
				}
				out[key] = fmt.Sprintf("%s", src)
			case strings.Compare(ext, ".rst") == 0:
				src, err := pandocProcessor(buf, "rst", "html")
				if err != nil {
					return nil, err
				}
				out[key] = fmt.Sprintf("%s", src)
			case strings.Compare(ext, ".textile") == 0:
				src, err := pandocProcessor(buf, "textile", "html")
				if err != nil {
					return nil, err
				}
				out[key] = fmt.Sprintf("%s", src)
			case strings.Compare(ext, ".jira") == 0:
				src, err := pandocProcessor(buf, "jira", "html")
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

//
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
