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

// pandocProcessor accepts an array of bytes as input and returns
// a `pandoc -f markdown -t html` output of an array if
// bytes and error.
func pandocProcessor(input []byte) ([]byte, error) {
	var out bytes.Buffer

	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(pandoc, "-f", "markdown", "-t", "html")
	cmd.Stdin = bytes.NewReader(input)
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// MakePandoc resolves key/value map rendering metadata suitable for processing with pandoc along with template information
// rendering and returns an error if something goes wrong
func MakePandoc(wr io.Writer, templateName string, keyValues map[string]string) error {
	var (
		out     bytes.Buffer
		options []string
	)

	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		return fmt.Errorf("Pandoc (see https://pandoc.org): %q", err)
	}
	data, err := ResolveData(keyValues)
	if err != nil {
		return fmt.Errorf("Can't resolve data source %s", err)
	}
	// NOTE: Pandocs default template expects content to be called $body$.
	// we need to remap from data["content"] to data["body"]
	if val, ok := data["content"]; ok == true {
		delete(data, "content")
		data["body"] = val
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
	// Check if document has front matter, split and write to temp files.
	defer os.Remove(metadata.Name())
	options = []string{
		"-s", "-t", "html",
		"--metadata-file", metadata.Name(),
	}
	if templateName != "" {
		options = append(options, []string{"--template", templateName}...)
	}
	cmd := exec.Command(pandoc, options...)
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s, %s", pandoc, err)
	}
	wr.Write(out.Bytes())
	return nil
}
