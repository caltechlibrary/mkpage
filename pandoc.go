// Package mkpage is an experimental static site generator
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
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

var (
	// Set a default -f (from) value used by Pandoc
	PandocFrom string
	// Set a default -t (to) value used by Pandoc
	PandocTo string
)

// Return the Pandoc version that will be used when calling Pandoc.
func GetPandocVersion() (string, error) {
	var (
		out, eOut bytes.Buffer
	)
	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		return "", err
	}
	cmd := exec.Command(pandoc, "--version")
	cmd.Stdout = &out
	cmd.Stderr = &eOut
	err = cmd.Run()
	if err != nil {
		if eOut.Len() > 0 {
			err = fmt.Errorf("%q says, %s\n%s", pandoc, eOut.String(), err)
		} else {
			err = fmt.Errorf("%q exit error, %s", pandoc, err)
		}
		return "", err
	}
	if eOut.Len() > 0 {
		fmt.Fprintf(os.Stderr, "%q warns, %s", pandoc, eOut.String())
	}
	return out.String(), err
}

// pandocProcessor accepts an array of bytes as input and returns
// a `pandoc -f {From} -t html` output of an array if
// bytes and error.
func pandocProcessor(input []byte, from string, to string) ([]byte, error) {
	var (
		out, eOut bytes.Buffer
	)

	if from == "" {
		from = PandocFrom
	}
	if to == "" {
		to = PandocTo
	}
	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		return nil, err
	}
	options := []string{}
	if from != "" {
		options = append(options, "-f", from)
	}
	if to != "" {
		options = append(options, "-t", to)
	}
	cmd := exec.Command(pandoc, options...)
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stdout = &out
	cmd.Stderr = &eOut
	err = cmd.Run()
	if err != nil {
		if eOut.Len() > 0 {
			err = fmt.Errorf("%q says, %s\n%s", pandoc, eOut.String(), err)
		} else {
			err = fmt.Errorf("%q exit error, %s", pandoc, err)
		}
		return nil, err
	}
	if eOut.Len() > 0 {
		fmt.Fprintf(os.Stderr, "%q warns, %s", pandoc, eOut.String())
	}
	return out.Bytes(), err
}

// MakePandoc resolves key/value map rendering metadata suitable for processing with pandoc along with template information
// rendering and returns an error if something goes wrong
func MakePandoc(wr io.Writer, templateName string, keyValues map[string]string) error {
	var (
		out, eOut bytes.Buffer
		options   []string
	)

	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		return fmt.Errorf("Pandoc (see https://pandoc.org): %q", err)
	}
	data, err := ResolveData(keyValues)
	if err != nil {
		return fmt.Errorf("Data resolution error: %s", err)
	}
	// NOTE: If a template is not provided (empty string) then
	// see is one is specified in the metadata
	if templateName == "" {
		if val, ok := data["template"]; ok == true {
			templateName = val.(string)
		}
	}
	// NOTE: Pandocs default template expects content to be called $body$.
	// we need to remap from data["content"] to data["body"] otherwise
	// we need to look in data to see if a template was specified.
	if templateName == "" {
		if val, ok := data["content"]; ok == true {
			delete(data, "content")
			data["body"] = val
		}
	}
	// NOTE: when using a template, title metadata is required.
	if _, ok := data["title"]; !ok {
		// Insert a title to prevent warning.
		data["title"] = "..."
	}

	src, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Marshal error, %q", err)
	}
	metadata, err := ioutil.TempFile(".", "pandoc.*.json")
	if err != nil {
		return fmt.Errorf("Cannot create temp metadata file, %s", err)
	}
	if _, err := metadata.Write(src); err != nil {
		return fmt.Errorf("Write error, %q", err)
	}
	defer os.Remove(metadata.Name())
	// Check if document has front matter, split and write to temp files.
	options = []string{}
	if PandocFrom != "" {
		options = append(options, "-f", PandocFrom)
	}
	if PandocTo != "" {
		options = append(options, "-t", PandocTo)
	}
	options = append(options, "--metadata-file", metadata.Name())
	if templateName != "" {
		options = append(options, []string{"--template", templateName}...)
	} else {
		options = append(options, "--standalone")
	}
	cmd := exec.Command(pandoc, options...)
	cmd.Stdout = &out
	cmd.Stderr = &eOut
	err = cmd.Run()
	if err != nil {
		if eOut.Len() > 0 {
			err = fmt.Errorf("%q says, %s\n%s", pandoc, eOut.String(), err)
		} else {
			err = fmt.Errorf("%q exit error, %s", pandoc, err)
		}
		return err
	}
	if eOut.Len() > 0 {
		fmt.Fprintf(os.Stderr, "%q warns, %s", pandoc, eOut.String())
	}
	wr.Write(out.Bytes())
	return err
}

// MakePandocString resolves key/value map rendering metadata suitable for processing with pandoc along with template information
// rendering and returns an error if something goes wrong
func MakePandocString(tmplSrc string, keyValues map[string]string) (string, error) {
	var (
		out, eOut bytes.Buffer
		options   []string
	)

	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		return "", fmt.Errorf("Pandoc (see https://pandoc.org): %q", err)
	}
	data, err := ResolveData(keyValues)
	if err != nil {
		return "", fmt.Errorf("Data resolution error: %s", err)
	}

	src, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("Marshal error, %q", err)
	}
	metadata, err := ioutil.TempFile(".", "pandoc.*.json")
	if err != nil {
		return "", fmt.Errorf("Cannot create temp metadata file, %s", err)
	}
	if _, err := metadata.Write(src); err != nil {
		return "", fmt.Errorf("Write error, %q", err)
	}
	defer os.Remove(metadata.Name())

	options = []string{}
	if PandocFrom != "" {
		options = append(options, "-f", PandocFrom)
	}
	if PandocTo != "" {
		options = append(options, "-t", PandocTo)
	}
	options = append(options, "--metadata-file", metadata.Name())
	if tmplSrc != "" {
		// Pandoc expects to read the template from disc so write
		// out to a temp file.
		// Check if document has front matter, split and write to temp files.
		template, err := ioutil.TempFile(".", "pandoc.*.tmpl")
		if err != nil {
			return "", fmt.Errorf("Cannot create temp template file, %s", err)
		}
		if _, err := template.Write([]byte(tmplSrc)); err != nil {
			return "", fmt.Errorf("Write error, %q", err)
		}
		defer os.Remove(template.Name())
		options = append(options, []string{"--template", template.Name()}...)
	} else {
		options = append(options, "--standalone")
	}
	cmd := exec.Command(pandoc, options...)
	cmd.Stdout = &out
	cmd.Stderr = &eOut
	err = cmd.Run()
	if err != nil {
		if eOut.Len() > 0 {
			err = fmt.Errorf("%q says, %s\n%s", pandoc, eOut.String(), err)
		} else {
			err = fmt.Errorf("%q exit error, %s", pandoc, err)
		}
		return "", err
	}
	if eOut.Len() > 0 {
		return "", fmt.Errorf("%q warns, %s", pandoc, eOut.String())
	}
	return fmt.Sprintf("%s", out.Bytes()), nil
}

// PandocBlock will attempt to convert a src byte array from
// one to another formats. If conversion is successful it returns
// a converted string, if unsuccessful it returns the original
// byte array as a string.
func PandocBlock(src string, from string, to string) string {
	if buf, err := pandocProcessor([]byte(src), from, to); err == nil {
		return fmt.Sprintf("%s", buf)
	}
	return src
}
