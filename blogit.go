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

const (
	DateFmt = "2006-01-02"
)

type CreatorObj struct {
	ORCID string `json:"orcid,omitempty"`
	Name  string `json:"name,omitempty"`
}

type PostObj struct {
	Slug        string       `json:"slug"`
	Document    string       `json:"document"`
	Title       string       `json:"title,omitempty"`
	SubTitle    string       `json:"subtitle,omitempty"`
	Byline      string       `json:"byline,omitempty"`
	Series      string       `json:"series,omitempty"`
	Number      string       `json:"number,omitempty"`
	Subject     string       `json:"subject,omitempty"`
	Keywords    []string     `json:"keywords,omitempty"`
	Abstract    string       `json:"abstract,omitempty"`
	Description string       `json:"description,omitempty"`
	Category    string       `json:"category,omitempty"`
	Lang        string       `json:"lang,omitempty"`
	Direction   string       `json:"direction,omitempty"`
	Draft       bool         `json:"draft,omitempty"`
	Creators    []CreatorObj `json:"creators,omitempty"`
	Created     string       `json:"date,omitempty"`
	Updated     string       `json:"updated,omitempty"`
}

type DayObj struct {
	Day   string     `json:"day"`
	Posts []*PostObj `json:"posts"`
}

type MonthObj struct {
	Month string    `json:"month"`
	Days  []*DayObj `json:"days"`
}

type YearObj struct {
	Year   string      `json:"year"`
	Months []*MonthObj `json:"months"`
}

type BlogMeta struct {
	Name    string     `json:"name,omitempty"`
	Quip    string     `json:"quip,omitempty"`
	BaseUrl string     `json:"url,omitempty"`
	Updated string     `json:"updated,omitempty"`
	Years   []*YearObj `json:"years"`
}

//
// Support funcs
//

func calcYMD(dateString string) ([]string, error) {
	ymd := strings.SplitN(dateString, "-", 3)
	return ymd, nil
}

func calcPath(prefix string, ymd []string) (string, error) {
	if len(ymd) != 3 {
		return "", fmt.Errorf("Invalid Year, Month and Date")
	}
	dPath := path.Join(prefix, ymd[0], ymd[1], ymd[2])
	return dPath, nil
}

func unpackCreators(objects []interface{}) []CreatorObj {
	creators := []CreatorObj{}
	for _, obj := range objects {
		switch obj.(type) {
		case string:
			creator := CreatorObj{}
			creator.Name = obj.(string)
			creators = append(creators, creator)
		case map[string]string:
			m := obj.(map[string]string)
			creator := CreatorObj{}
			if name, ok := m["name"]; ok {
				creator.Name = name
			}
			if orcid, ok := m["orcid"]; ok {
				creator.ORCID = orcid
			}
			creators = append(creators, creator)
		}
	}
	return creators
}

//
// Exported funcs
//

// postIndex looks through the day's posts and find
// the position that matches the slug or returns -1
func (d DayObj) postIndex(slug string) int {
	postIndex := -1
	for i, obj := range d.Posts {
		if obj.Slug == slug {
			postIndex = i
			break
		}
	}
	return postIndex
}

// dayIndex looks through the months days and returns
// index position found or -1 if not found.
func (m MonthObj) dayIndex(day string) int {
	dayIndex := -1
	for i, obj := range m.Days {
		if obj.Day == day {
			dayIndex = i
			break
		}
	}
	return dayIndex
}

// monthIndex looks through a years' months and returns
// the index position found or -1 if not found.
func (y YearObj) monthIndex(month string) int {
	monthIndex := -1
	for i, obj := range y.Months {
		if obj.Month == month {
			monthIndex = i
			break
		}
	}
	return monthIndex
}

// yearIndex
func (meta BlogMeta) yearIndex(year string) int {
	yearIndex := -1
	for i, obj := range meta.Years {
		if obj.Year == year {
			yearIndex = i
			break
		}
	}
	return yearIndex
}

// updatePosts will create a new post if necessary and insert in to the
// post list.
func (dy *DayObj) updatePosts(ymd []string, targetName string) error {

	// Read in front matter from targetName
	src, err := ioutil.ReadFile(targetName)
	if err != nil {
		return fmt.Errorf("Failed to read %q, %s", targetName, err)
	}
	obj := map[string]interface{}{}
	fmType, src, _ := SplitFrontMatter(src)
	if len(src) > 0 {
		if err := UnmarshalFrontMatter(fmType, src, &obj); err != nil {
			return fmt.Errorf("Failed to unmarshal front matter %q, %s", targetName, err)
		}
	}
	// Create a new PostObj
	today := time.Now().Format(DateFmt)
	created := strings.Join(ymd, "-")
	post := new(PostObj)
	post.Document = targetName
	post.Updated = today
	post.Created = created
	post.Slug = strings.TrimSuffix(path.Base(targetName), filepath.Ext(targetName))
	if title, ok := obj["title"]; ok {
		post.Title = title.(string)
	}
	if subtitle, ok := obj["subtitle"]; ok {
		post.SubTitle = subtitle.(string)
	}
	if byline, ok := obj["byline"]; ok {
		post.Byline = byline.(string)
	}
	if series, ok := obj["series"]; ok {
		post.Series = series.(string)
	}
	if number, ok := obj["number"]; ok {
		post.Number = number.(string)
	}
	if subject, ok := obj["subject"]; ok {
		post.Subject = subject.(string)
	}
	if objects, ok := obj["keywords"]; ok {
		switch objects.(type) {
		case []string:
			keywords := objects.([]string)
			for _, kw := range keywords {
				post.Keywords = append(post.Keywords, kw)
			}
		}
	}
	if abstract, ok := obj["abstract"]; ok {
		post.Abstract = abstract.(string)
	}
	if description, ok := obj["description"]; ok {
		post.Description = description.(string)
	}
	if category, ok := obj["category"]; ok {
		post.Category = category.(string)
	}
	if lang, ok := obj["lang"]; ok {
		post.Lang = lang.(string)
	}
	if direction, ok := obj["direction"]; ok {
		post.Direction = direction.(string)
	}
	if draft, ok := obj["draft"]; ok {
		post.Draft = draft.(bool)
	} else {
		post.Draft = false
	}
	if creators, ok := obj["creators"]; ok {
		post.Creators = unpackCreators(creators.([]interface{}))
	}
	if dt, ok := obj["date"]; ok {
		post.Created = dt.(string)
	}
	if updated, ok := obj["updated"]; ok {
		post.Updated = updated.(string)
	}

	i := dy.postIndex(post.Slug)
	if i < 0 {
		// Add a post
		post.Created = today
		posts := dy.Posts[0:]
		dy.Posts = append([]*PostObj{post}, posts...)
	} else {
		// Update a post
		dy.Posts[i] = post
	}
	return nil
}

// updateDays will create a new day and insert in order
// before passing the post data to UpdatePost()
func (mn *MonthObj) updateDays(ymd []string, targetName string) error {
	dy := new(DayObj)
	dy.Day = ymd[2]
	i := mn.dayIndex(dy.Day)
	if i < 0 {
		if len(mn.Days) == 0 {
			mn.Days = append(mn.Days, dy)
			i = 0
		} else {
			for j, obj := range mn.Days {
				if dy.Day > obj.Day {
					i = j
					if i == 0 {
						mn.Days = append([]*DayObj{dy}, mn.Days...)
					} else {
						days := mn.Days[0 : j-1]
						days = append(days, dy)
						mn.Days = append(days, mn.Days[j:]...)
					}
					break
				}
			}
			if i < 0 {
				mn.Days = append(mn.Days, dy)
				i = len(mn.Days) - 1
			}
		}
	}
	return mn.Days[i].updatePosts(ymd, targetName)
}

// updateMonths will create/update month
// before passing the post data to UpdateDays()
func (yr *YearObj) updateMonths(ymd []string, targetName string) error {
	mn := new(MonthObj)
	mn.Month = ymd[1]
	i := yr.monthIndex(mn.Month)
	if i < 0 {
		if len(yr.Months) == 0 {
			yr.Months = append(yr.Months, mn)
			i = 0
		} else {
			for j, obj := range yr.Months {
				if mn.Month > obj.Month {
					i = j
					if i == 0 {
						yr.Months = append([]*MonthObj{mn}, yr.Months...)
					} else {
						months := yr.Months[0 : j-1]
						months = append(months, mn)
						yr.Months = append(months, yr.Months[j:]...)
					}
					break
				}
			}
			if i < 0 {
				yr.Months = append(yr.Months, mn)
				i = len(yr.Months) - 1
			}
		}
	}
	return yr.Months[i].updateDays(ymd, targetName)
}

// updateYears will create/update year in `meta.Years`
// before passing the post data to UpdateMonths()
func (meta *BlogMeta) updateYears(ymd []string, targetName string) error {
	yr := new(YearObj)
	yr.Year = ymd[0]
	i := meta.yearIndex(yr.Year)
	if i < 0 {
		if len(meta.Years) == 0 {
			meta.Years = append(meta.Years, yr)
			i = 0
		} else {
			for j, obj := range meta.Years {
				if yr.Year > obj.Year {
					i = j
					if i == 0 {
						meta.Years = append([]*YearObj{yr}, meta.Years...)
					} else {
						years := meta.Years[0 : j-1]
						years = append(years, yr)
						meta.Years = append(years, meta.Years[j:]...)
					}
					break
				}
			}
			if i < 0 {
				// We need to append the year, it's earlier than
				// known years.
				meta.Years = append(meta.Years, yr)
				i = len(meta.Years) - 1
			}
		}
	}
	return meta.Years[i].updateMonths(ymd, targetName)
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
func (meta *BlogMeta) BlogIt(prefix string, fName string, dateString string) error {
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
	dPath, err := calcPath(prefix, ymd)
	if err != nil {
		return err
	}
	// copy fName to target path.
	var (
		in, out *os.File
	)
	in, err = os.Open(fName)
	if err != nil {
		return err
	} else {
		os.MkdirAll(dPath, 0777)
		targetName = path.Join(dPath, path.Base(fName))
		out, err = os.Create(targetName)
		if err != nil {
			return fmt.Errorf("Creating %q, %s", targetName, err)
		}
		if _, err := io.Copy(out, in); err != nil {
			return err
		}
		in.Close()
		out.Close()
	}
	// NOTE: Updated is always today.
	meta.Updated = time.Now().Format(DateFmt)
	return meta.updateYears(ymd, targetName)
}

// Save writes a JSON blog meta document
func (meta *BlogMeta) Save(fName string) error {
	src, err := json.MarshalIndent(meta, "", "    ")
	if err != nil {
		return fmt.Errorf("Marshaling %q, %s", fName, err)
	}
	err = ioutil.WriteFile(fName, src, 0666)
	if err != nil {
		return fmt.Errorf("Writing %q, %s", fName, err)
	}
	return nil
}

// Reads a JSON blog meta document and popualtes a blog meta structure
func LoadBlogMeta(fName string, meta *BlogMeta) error {
	src, err := ioutil.ReadFile(fName)
	if err != nil {
		return fmt.Errorf("Reading %q, %s", fName, err)
	}
	if len(src) > 0 {
		if err := json.Unmarshal(src, meta); err != nil {
			return fmt.Errorf("Unmarshing %q, %s", fName, err)
		}
	}
	return nil
}

// hasExt checks if ext is in list of target ext.
func hasExt(ext string, targetExts []string) bool {
	for _, e := range targetExts {
		if strings.Compare(ext, e) == 0 {
			return true
		}
	}
	return false
}

// RefreshFromPath crawls the dircetory tree and rebuilds
// the `blog.json` file based on what is found. It takes a
// File extension to target (e.g. .md for Markdown) and
// analyzes the path for YYYY/MM/DD and transforms the
// information found into an entry in `blog.json`.
func (meta *BlogMeta) RefreshFromPath(prefix string, year string) error {
	var (
		ymd []string
	)
	targetExts := []string{
		".md",
		".rst",
		".textile",
		".jira",
		".txt",
	}
	months := map[string]int{
		"01": 31, "02": 29, "03": 31, "04": 30,
		"05": 31, "06": 30, "07": 31, "08": 31,
		"09": 30, "10": 31, "11": 30, "12": 31,
	}
	ymd = append(ymd, year, "", "")
	for month, cnt := range months {
		ymd[1] = month
		for day := 1; day <= cnt; day++ {
			ymd[2] = fmt.Sprintf("%02d", day)
			// CalcPath and find files.
			folder := path.Join(prefix, ymd[0], ymd[1], ymd[2])
			// Scan the fold for files ending in ext,
			files, err := ioutil.ReadDir(folder)
			if err == nil {
				// for each file with matching extension run updateYear(ymd, targetName)
				for _, file := range files {
					targetName := path.Join(prefix, ymd[0], ymd[1], ymd[2], file.Name())
					ext := filepath.Ext(targetName)
					if hasExt(ext, targetExts) {
						if err := meta.updateYears(ymd, targetName); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return nil
}
