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
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"
)

func TestPrivateFuncs(t *testing.T) {
	rsdoiel := "R. S. Doiel"
	obj := map[string]string{}
	obj["name"] = rsdoiel
	objs := append([]interface{}{}, obj)
	creators := unpackCreators(objs)
	if len(creators) != 1 {
		t.Errorf("expected one creator got %d", len(creators))
	} else {
		if fmt.Sprintf("%T", creators[0].Name) != "string" {
			t.Errorf("Expected crestors[0].Name to ba a string, got %T", creators[0].Name)
		}
		if strings.Compare(creators[0].Name, rsdoiel) != 0 {
			t.Errorf("expected creators[0].Name = %q, got %q", rsdoiel, creators[0].Name)
		}
	}
	ymd, err := calcYMD("2003-01-02")
	if err != nil {
		t.Errorf("Unexpected err %s, calcYMD()", err)
		t.FailNow()
	}
	if len(ymd) != 3 {
		t.Errorf("expected len(ymd) = 3, got %+v", ymd)
	} else {
		for i, val := range []string{"2003", "01", "02"} {
			if ymd[i] != val {
				t.Errorf("expected %q, got %q", val, ymd[i])
			}
		}
		expectedS := path.Join("blog", "2003", "01", "02")
		resultS, err := calcPath("blog", ymd)
		if err != nil {
			t.Errorf("Unexpected err %s, calcPath()", err)
		}

		if expectedS != resultS {
			t.Errorf("expected %q, got %q", expectedS, resultS)
		}
	}
}

func TestExportedFuncs(t *testing.T) {
	var (
		pName      string
		prefix     string
		blogPrefix string
		blogJSON   string
		src        []byte
		blogMeta   *BlogMeta
	)
	prefix = "test"
	blogPrefix = path.Join(prefix, "blog")
	blogJSON = path.Join(blogPrefix, "blog.json")

	// Start with an empty blog ...
	os.RemoveAll(path.Join("test", "blog"))
	blogMeta = new(BlogMeta)

	// Generate and write test data for BlogIt()
	for i := 1; i <= 10; i = i + 1 {
		src = []byte(fmt.Sprintf(`{
	"title": "Hello No. %d",
	"subtitle": "This is the %d(th) test blog post",
	"date": "2020-05-%02d",
	"keywords": [ "test" ],
	"creators": [ "R. S. Doiel" ],
	"byline": "By R. S. Doiel"
}


# Hello World!

Test Blog post.
`, i, i, i))
		pName = path.Join(prefix, fmt.Sprintf("hello_%d.md", i))
		if err := ioutil.WriteFile(pName, src, 0666); err != nil {
			t.Errorf("Can't created %q, %s", pName, err)
			t.FailNow()
		}

		dateString := fmt.Sprintf("2020-05-%02d", i)
		if err := blogMeta.BlogIt(blogPrefix, pName, dateString); err != nil {
			t.Errorf("BlogIt(%q, %q, %q) failed, %s", blogPrefix, pName, dateString, err)
			t.FailNow()
		}
		if err := blogMeta.Save(blogJSON); err != nil {
			t.Errorf("Failed to write %q, %s", blogJSON, err)
			t.FailNow()
		}
	}
}
