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
	"log"
	"os"
	"os/exec"
)

// MakePandoc resolves key/value map rendering metadata suitable for processing with pandoc along with template information
// rendering and returns an error if something goes wrong
func MakePandoc(wr io.Writer, templateName string, keyValues map[string]string) error {
	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		log.Fatalf("Pandoc (see https://pandoc.org): %q", err)
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
		log.Fatalf("Marshal error, %q", err)
	}
	metadata, err := ioutil.TempFile(".", "pandoc.*.json")
	if err != nil {
		log.Fatalf("Cannot create temp metadata file, %s", err)
	}
	if _, err := metadata.Write(src); err != nil {
		log.Fatalf("Write error, %q", err)
	}
	// Check if document has front matter, split and write to temp files.
	defer os.Remove(metadata.Name())
	//FIXME: take temp metadata file, pass to pandoc command and render output.
	// pandoc -s -f markdown -t html  --metadata-file=/var/folders/14/d4r81n210gq5p00l3_byft400000gn/T/pandoc.690211801.json
	cmd := exec.Command(pandoc, "-s", "-f", "markdown", "-t", "html", "--metadata-file", metadata.Name())
	var (
		out bytes.Buffer
	)
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Fatalf("%s, %s", pandoc, err)
	}
	wr.Write(out.Bytes())
	return nil
}
