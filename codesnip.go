// Package mkpage codesnip.go - extracts code snippets from markdown
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
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Codesnip works on io.Reader and ioWriter and string used
// to identified the langauge. It reads a Markdown block looking
// for code fenses and snips out the code to write to a file.
func Codesnip(in io.Reader, out io.Writer, language string) error {
	var (
		inCodeBlock bool
	)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				inCodeBlock = false
			} else if strings.HasPrefix(line, "```"+language) {
				inCodeBlock = true
				continue
			}
		}
		if inCodeBlock {
			switch {
			case strings.HasPrefix(line, "    "):
				fmt.Fprintln(out, strings.TrimPrefix(line, "    "))
			case strings.HasPrefix(line, "\t"):
				fmt.Fprintln(out, strings.TrimPrefix(line, "\t"))
			default:
				fmt.Fprintln(out, line)
			}
		}
	}
	return scanner.Err()
}
