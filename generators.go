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
	"os"
	"os/exec"
	"strings"
)

// JSONGenerator accepts  command line string and executes it.
// It take command's output, validates that it is JSON and returns it.
func JSONGenerator(cmdExpr string, obj interface{}) error {
	var (
		out, eOut bytes.Buffer
		generator string
		params    []string
		err       error
	)
	line := strings.Split(cmdExpr, " ")
	switch len(line) {
	case 0:
		err = fmt.Errorf("Missing generator command")
		return err
	case 1:
		generator = cmdExpr
	default:
		generator, params = line[0], line[1:]
	}

	cmd := exec.Command(generator, params...)
	cmd.Stdout = &out
	cmd.Stderr = &eOut
	err = cmd.Run()
	if err != nil {
		if eOut.Len() > 0 {
			err = fmt.Errorf("%q says, %s\n%s", cmdExpr, eOut.String(), err)
		} else {
			err = fmt.Errorf("%q exit error, %s", cmdExpr, err)
		}
		return err
	}
	if eOut.Len() > 0 {
		fmt.Fprintf(os.Stderr, "%q warns, %s", cmdExpr, eOut.String())
		return err
	}
	src := out.Bytes()
	//NOTE: Validate our JSON by trying to unmarshaling it
	err = json.Unmarshal(src, &obj)
	if err != nil {
		err = fmt.Errorf("Invalid JSON from %q exit error, %s", cmdExpr, err)
	}
	return err
}
