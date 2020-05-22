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
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type CreatorObj struct {
	ORCID    string `json:"orcid,omitempty"`
	Name     string `json:"name,omitempty"`
	SortName string `json:"sort_name,omitempty"`
}

type PostObj struct {
	Slug        string       `json:"slug,omitempty"`
	Title       string       `json:"title,omitempty"`
	SubTitle    string       `json:"subtitle,omitempty"`
	Byline      string       `json:"byline,omitempty"`
	Series      string       `json:"series,omitempty"`
	No          string       `json:"number,omitempty"`
	Subject     string       `json:"subject,omitempty"`
	Keywords    []string     `json:"keywords,omitempty"`
	Abstract    string       `json:"abstract,omitempty"`
	Description string       `json:"description,omitempty"`
	Category    string       `json:"category,omitempty"`
	Lang        string       `json:"lang,omitempty"`
	Dir         string       `json:"dir,omitempty"`
	Draft       bool         `json:"draft,omitempty"`
	Creators    []CreatorObj `json:"creators,omitempty"`
	Created     *time.Time   `json:"date,omitempty"`
	Updated     *time.Time   `json:"updated,omitempty"`
}

type DayObj struct {
	Day   string    `json:"Day,omitempty"`
	Posts []PostObj `json:"posts,omitempty"`
	Count int       `json:"count,omitempty"`
}

type MonthObj struct {
	Month  string     `json:"Month,omitempty"`
	Days   []DayObj   `json:"days,omitempty"`
	Months []MonthObj `json:"months,omitempty"`
	Count  int        `json:"count,omitempty"`
}

type YearObj struct {
	Year   string     `json:"Year,omitempty"`
	Months []MonthObj `json:"months,omitempty"`
	Count  int        `json:"count,omitempty"`
}

type BlogMeta struct {
	Name    string    `json:"name,omitempty"`
	Quote   string    `json:"quote,omitemtpy"`
	BaseUrl string    `json:"url,omitempty"`
	Status  string    `json:"status,omitempty"`
	Updated string    `json:"updated,omitempty"`
	Years   []YearObj `json:"years,omitempty"`
}

func (meta *BlogMeta) Update(isNew bool, ymd []string, targetName string) {
	// Create PostObj
	post := new(PostObj)
	bName := path.Base(targetName)
	post.Slug = strings.TrimSuffix(bName, filepath.Ext(bName))

	// Get metadata from targetName and map into post
	// Check to see if posts already exists, if so replace it otherwise, insert it

	meta.Status = "active"
	meta.Updated = fmt.Sprintf("%s-%s-%s", ymd[0], ymd[1], ymd[2])
}

func calcYMD(dateString string) ([]string, error) {
	dt := new(time.Time)
	if err := dt.UnmarshalText([]byte(dateString)); err != nil {
		return nil, err
	}
	ymd := []string{}
	ymd = append(ymd, dt.Format("2015"), dt.Format("05"), dt.Format("07"))
	fmt.Printf("DEBUG ymd -> %+v\n", ymd)
	return ymd, nil
}

func calcPath(prefix string, ymd []string) (string, error) {
	if len(ymd) != 3 {
		return "", fmt.Errorf("Invalid Year, Month and Date")
	}
	dPath := path.Join(prefix, ymd[0], ymd[1], ymd[2])
	fmt.Printf("DEBUG dPath -> %q", dPath)
	return dPath, nil
}

// BlogIt is a tool for posting and updating a blog directory
// structure on local disc.  It includes maintaining additional
// metadata resources to make it easy to script blogsites and
// podcast sites.
// @param prefix - a prefix path for the blog (relative to working dir)
// @param fName - the name of the file to publish to generated blog path
// @param dateString - A date string used to calculate the blog path, e.g.
//                  YYYY-MM-DD maps to YYYY/MM/DD.
// @returns an error type
func BlogIt(prefix string, fName string, dateString string) error {
	var (
		targetName string
	)
	// Check to see if dateStr is empty, if so default to today.
	if dateString == "" {
		dateString = time.Now().Format("2015-05-07")
	}

	// Check to see if path.join(prefix, datePath(dateStr)) exists
	// and create it if needed.
	ymd, err := calcYMD(dateString)
	if err != nil {
		return err
	}
	dPath, err := calcPath(prefix, ymd...)
	if err != nil {
		return err
	}
	blogitJSON := path.Join(prefix, "blogit.json")
	if _, err := os.Stat(dPath); os.IsNotExist(err) {
		fmt.Printf("Creating %s\n", dPath)
		if err := os.MkdirAll(dPath, 0666); err != nil {
			return err
		}
	}
	blogitMeta := new(BlogMeta)
	if _, err := os.Stat(blogitJSON); err != nil {
		if src, err := ioutil.ReadFile(blogitJSON); err != nil {
			fmt.Printf("Reading %s\n", blogitJSON)
			if err := json.Unmarshal(src, &blogitMeta); err != nil {
				return err
			}
		}
	}

	// copy fName to target path.
	var (
		in, out *os.File
		isNew   bool
	)
	in, err := os.Open(fName)
	if err != nil {
		return err
	}
	defer in.Close()
	isNew = false
	targetName = path.Join(dPath, path.Base(fName))
	if _, err := os.Stat(targetName); os.IsNotExist(err) {
		isNew = true
	}
	out, err = os.Create(targetName)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	// Update targetName in blogit.json and write updated blogit.json
	if err = blogitMeta.Update(isNew, ymd, targetName); err != nil {
		return err
	}
	src, err := json.MarshalIndent(blogitMeta, "", "    ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(blogitJSON, src, 0666)
	if err != nil {
		return err
	}
	return nil
}
