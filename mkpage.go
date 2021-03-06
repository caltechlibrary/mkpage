//
// Package mkpage is an experiment in a light weight template and markdown processor.
//
// @author R. S. Doiel, <rsdoiel@caltech.edu>
//
// Copyright (c) 2021, Caltech
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
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	// 3rd Part support (e.g. YAML)
	"gopkg.in/yaml.v3"

	// Fountain support for scripts, interviews and narration
	"github.com/rsdoiel/fountain"
)

const (
	// Prefix for explicit string types

	// JSONPrefix designates a string as JSON formatted content
	JSONPrefix = "json:"
	// CommonMarkPrefix designates a string as Common Mark
	// (a rich markdown dialect) content
	CommonMarkPrefix = "commonmark:"
	// MarkdownPrefix designates a string as Markdown (pandoc's dialect)
	// content
	MarkdownPrefix = "markdown:"
	// MarkdownStrict designates a strnig as John Gruber's Markdown content
	MarkdownStrictPrefix = "markdown_strict:"
	// GfmMarkdownPrefix designates a string as GitHub Flavored Markdown
	GfmMarkdownPrefix = "gfm:"
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
	// JSONGeneratorPrefix evaluates the value as a command line that
	// returns JSON.
	JSONGeneratorPrefix = "json-generator:"

	// DateExp is the default format used by mkpage utilities for date exp
	DateExp = `[0-9][0-9][0-9][0-9]-[0-1][0-9]-[0-3][0-9]`
	// BylineExp is the default format used by mkpage utilities
	BylineExp = (`^[B|b]y\s+(\w|\s|.)+` + DateExp + "$")
	// TitleExp is the default format used by mkpage utilities
	TitleExp = `^#\s+(\w|\s|.)+$`

	//
	// Supported types for Front Matter
	//

	// FrontMatterIsUnknown means front matter and we can't parse it
	FrontMatterIsUnknown = iota
	// FrontMatterIsJSON means we have detected JSON front matter
	FrontMatterIsJSON
	// FrontMatterIsPandocMetadata means we have detected a Pandoc
	// style metadata block, e.g. opening lines start with
	// '%' attribute name followed by value(s)
	// E.g.
	//      % title
	//      % author(s)
	//      % date
	FrontMatterIsPandocMetadata
	// FrontMatterIsYAML means we have detected a Pandoc YAML
	// front matter block.
	FrontMatterIsYAML
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

// SplitFrontMatter takes a []byte input splits it into front matter type,
// front matter source and Markdown source. If either is missing an
// empty []byte is returned for the missing element.
// NOTE: removed yaml, toml support as of v0.2.4
// NOTE: Added support for Pandoc title blocks v0.2.5
func SplitFrontMatter(input []byte) (int, []byte, []byte) {
	// JSON front matter, most Markdown processors.
	if bytes.HasPrefix(input, []byte("{\n")) {
		parts := bytes.SplitN(bytes.TrimPrefix(input, []byte("{\n")), []byte("\n}\n"), 2)
		src := []byte(fmt.Sprintf("{\n%s\n}\n", parts[0]))
		if len(parts) > 1 {
			return FrontMatterIsJSON, src, parts[1]
		}
		return FrontMatterIsJSON, src, []byte("")
	}
	if bytes.HasPrefix(input, []byte("---\n")) {
		parts := bytes.SplitN(bytes.TrimPrefix(input, []byte("---\n")), []byte("\n---\n"), 2)
		src := []byte(fmt.Sprintf("---\n%s\n---\n", parts[0]))
		if len(parts) > 1 {
			return FrontMatterIsYAML, src, parts[1]
		}
		return FrontMatterIsYAML, src, []byte("")
	}
	if bytes.HasPrefix(input, []byte("% ")) {
		lines := bytes.Split(input, []byte("\n"))
		i := 0
		fieldCnt := 0
		src := []byte{}
		for ; (i < len(lines)) && (fieldCnt < 3); i++ {
			if bytes.HasPrefix(lines[i], []byte("% ")) {
				fieldCnt += 1
				src = append(append(src, lines[i]...), []byte("\n")...)
			} else if fieldCnt < 3 {
				//NOTE: Dates can only one line, so we stop extra
				// line consumption with authors.
				src = append(append(src, lines[i]...), []byte("\n")...)
			}
		}
		if fieldCnt == 3 {
			return FrontMatterIsPandocMetadata, src, input[len(src):]
		}
	}
	// Handle case of no front matter
	return FrontMatterIsUnknown, []byte(""), input
}

// UnmarshalFrontMatter takes a []byte of front matter source
// and unmarshalls using only JSON frontmatter
// NOTE: removed yaml, toml support as of v0.2.4
// NOTE: Added support for Pandoc title blocks as of v0.2.5
func UnmarshalFrontMatter(configType int, src []byte, obj *map[string]interface{}) error {
	var (
		txt []byte
		err error
	)
	switch configType {
	case FrontMatterIsPandocMetadata:
		block := MetadataBlock{}
		if err = block.Unmarshal(txt); err != nil {
			return err
		}
		if txt, err = block.Marshal(); err != nil {
			return nil
		}
		if err = json.Unmarshal(txt, &obj); err != nil {
			return err
		}
	case FrontMatterIsJSON:
		// Make sure we have valid JSON
		if err = json.Unmarshal(src, &obj); err != nil {
			return err
		}
	case FrontMatterIsYAML:
		if err = yaml.Unmarshal(src, &obj); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Unsupported Front matter format")
	}
	return nil
}

// ProcessorConfig takes front matter and returns
// a map[string]interface{} containing configuration
// NOTE: removed yaml, toml support as of v0.2.4
// NOTE: added Pandoc Metadata block as of v0.2.5
func ProcessorConfig(configType int, frontMatterSrc []byte) (map[string]interface{}, error) {
	//FIXME: Need to merge with .Config and return the merged result.
	m := map[string]interface{}{}
	// Do nothing is we have zero front matter to process.
	if len(frontMatterSrc) == 0 {
		return m, nil
	}
	// Convert Front Matter to JSON
	switch configType {
	case FrontMatterIsPandocMetadata:
		block := MetadataBlock{}
		if err := block.Unmarshal(frontMatterSrc); err != nil {
			return nil, err
		}
		m["title"] = block.Title
		m["authors"] = block.Authors
		m["date"] = block.Date
	case FrontMatterIsJSON:
		// JSON Front Matter
		if err := json.Unmarshal(frontMatterSrc, &m); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown supported front matter format")
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
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, MMarkPrefix)), "markdown_mmd", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, CommonMarkPrefix) == true:
			//NOTE: We're using pandoc as our default processor
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, CommonMarkPrefix)), "commonmark_x", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, MarkdownPrefix) == true:
			//NOTE: We're using pandoc's flavor Markdown as our processor
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, MarkdownPrefix)), "markdown", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, MarkdownStrictPrefix) == true:
			//NOTE: We're using origanal John Gruber Markdown
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, MarkdownStrictPrefix)), "markdown_strict", "html")
			if err != nil {
				return out, err
			}
			out[key] = fmt.Sprintf("%s", src)
		case strings.HasPrefix(val, GfmMarkdownPrefix) == true:
			//NOTE: We're using pandoc as our default processor
			src, err := pandocProcessor([]byte(strings.TrimPrefix(val, GfmMarkdownPrefix)), "gfm", "html")
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
		case strings.HasPrefix(val, JSONGeneratorPrefix) == true:
			//NOTE: JSONGenerator expects a command line that results
			// in JSON written to stdout. It then passes this back to
			// be processed by pandoc in the metadata file.
			var o interface{}
			cmd := strings.TrimPrefix(val, JSONGeneratorPrefix)
			err := JSONGenerator(cmd, &o)
			if err != nil {
				return out, fmt.Errorf("(key: %q) %q failed, %s", key, cmd, err)
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
						src, err := pandocProcessor(buf, "", "html")
						if err != nil {
							return nil, err
						}
						out[key] = fmt.Sprintf("%s", src)
					case isContentType(contentTypes, "text/commonmark") == true:
						src, err := pandocProcessor(buf, "commonmark_x", "html")
						if err != nil {
							return nil, err
						}
						out[key] = fmt.Sprintf("%s", src)
					case isContentType(contentTypes, "text/mmark") == true:
						src, err := pandocProcessor(buf, "mmark", "html")
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
			ext := path.Ext(val)
			buf, err := ioutil.ReadFile(val)
			if err != nil {
				return out, fmt.Errorf("Can't read (%s) %q, %s", key, val, err)
			}
			//NOTE: We only split front matter for supported markup
			// formats, e.g. MultiMarkdown, CommonMark, Markdown, Textile,
			// ReStructureText, JiraText, Fountain
			if strings.Compare(ext, ".json") != 0 {
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
			}
			switch {
			case strings.Compare(ext, ".fountain") == 0 ||
				strings.Compare(ext, ".spmd") == 0:
				src, err := fountainProcessor(buf)
				if err != nil {
					return nil, err
				}
				out[key] = fmt.Sprintf("%s", src)
			case strings.Compare(ext, ".md") == 0:
				src, err := pandocProcessor(buf, "", "html")
				if err != nil {
					return nil, err
				}
				out[key] = fmt.Sprintf("%s", src)
			case strings.Compare(ext, ".mmd") == 0:
				src, err := pandocProcessor(buf, "markdown_mmd", "html")
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
