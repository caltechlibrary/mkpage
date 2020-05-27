package mkpage

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/caltechlibrary/rss2"
)

// Generate a Feed from walking the BlogMeta structure
func BlogMetaToRSS(blog *BlogMeta, articleLimit int, feed *rss2.RSS2) error {
	if len(blog.Name) > 0 {
		feed.Title = blog.Name
	}
	if len(blog.BaseURL) > 0 {
		feed.Link = blog.BaseURL
	}
	if len(blog.Quip) > 0 {
		feed.Description = "> " + blog.Quip + "\n\n"
	}
	if len(blog.Description) > 0 {
		feed.Description += blog.Description
	}
	if len(blog.Updated) > 0 {
		dt, err := time.Parse("2006-01-02", blog.Updated)
		if err != nil {
			return err
		}
		feed.PubDate = dt.Format(time.RFC1123)
		feed.LastBuildDate = dt.Format(time.RFC1123)
	}
	if len(blog.Language) > 0 {
		feed.Language = blog.Language
	}
	if len(blog.Copyright) > 0 {
		feed.Copyright = blog.Copyright
	}
	if feed.ItemList == nil {
		feed.ItemList = []rss2.Item{}
	}
	//FIXME: Need to iterate over years, months, days and build our
	// blog items.
	for _, years := range blog.Years {
		yr := years.Year
		for _, months := range years.Months {
			mn := months.Month
			for _, days := range months.Days {
				dy := days.Day
				for _, post := range days.Posts {
					pubDate, err := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", yr, mn, dy))
					if err != nil {
						return err
					}
					item := new(rss2.Item)
					item.Title = post.Title
					item.Link = strings.Join([]string{blog.BaseURL, post.Document}, "/")
					item.GUID = item.Link
					item.PubDate = pubDate.Format(time.RFC1123)
					if post.Abstract != "" {
						item.Description = post.Abstract
					} else if post.Description != "" {
						item.Description = post.Description
					}
					feed.ItemList = append(feed.ItemList, *item)
					fmt.Printf("DEBUG feed.ItemList (%d) -> %+v\n", len(feed.ItemList), feed.ItemList)
				}
			}
		}
	}
	return nil
}

// Generate an Feed by walking the file system.
func WalkRSS(feed *rss2.RSS2, htdocs string, excludeList string, articleLimit int, titleExp string, bylineExp string, dateExp string) error {
	// Required
	channelLink := feed.Link

	validBlogPath := regexp.MustCompile("/[0-9][0-9][0-9][0-9]/[0-9][0-9]/[0-9][0-9]/")
	err := Walk(htdocs, func(p string, info os.FileInfo) bool {
		fname := path.Base(p)
		if validBlogPath.MatchString(p) == true &&
			strings.HasSuffix(fname, ".md") == true {
			// NOTE: We have a possible published markdown article.
			// Make sure we have a HTML version before adding it
			// to the feed.
			if _, err := os.Stat(path.Join(p, path.Base(fname)+".html")); os.IsNotExist(err) {
				return false
			}
			return true
		}
		return false
	}, func(p string, info os.FileInfo) error {
		// Read the article
		buf, err := ioutil.ReadFile(p)
		if err != nil {
			return err
		}
		fMatter := map[string]interface{}{}
		fType, fSrc, tSrc := SplitFrontMatter(buf)
		if len(fSrc) > 0 {
			if err := UnmarshalFrontMatter(fType, fSrc, &fMatter); err != nil {
				fMatter = map[string]interface{}{}
			}
		}

		// Calc URL path
		pname := strings.TrimPrefix(p, htdocs)
		if strings.HasPrefix(pname, "/") {
			pname = strings.TrimPrefix(pname, "/")
		}
		dname := path.Dir(pname)
		bname := strings.TrimSuffix(path.Base(pname), ".md") + ".html"
		articleURL := fmt.Sprintf("%s/%s", channelLink, path.Join(dname, bname))
		u, err := url.Parse(articleURL)
		if err != nil {
			return err
		}
		// Collect metadata
		//NOTE: Use front matter if available otherwise
		var (
			title, byline, author, description, pubDate string
		)
		src := fmt.Sprintf("%s", buf)
		if val, ok := fMatter["title"]; ok {
			title = val.(string)
		} else {
			title = strings.TrimPrefix(Grep(titleExp, src), "# ")
		}
		if val, ok := fMatter["byline"]; ok {
			byline = val.(string)
		} else {
			byline = Grep(bylineExp, src)
		}
		if val, ok := fMatter["pubDate"]; ok {
			pubDate = val.(string)
		} else {
			pubDate = Grep(dateExp, byline)
		}
		if val, ok := fMatter["description"]; ok {
			description = val.(string)
		} else {
			description = OpeningParagraphs(fmt.Sprintf("%s", tSrc), 5, "\n\n")
			if len(description) < len(tSrc) {
				description += " ..."
			}
			description = PandocBlock(description, "markdown", "html")
			description = PandocBlock(description, "html", "xml")
		}
		if val, ok := fMatter["creator"]; ok {
			author = val.(string)
		} else if val, ok = fMatter["author"]; ok {
			author = val.(string)
		} else {
			author = byline
			if len(byline) > 2 {
				author = strings.TrimSpace(strings.TrimSuffix(byline[2:], pubDate))
			}
		}
		// Reformat pubDate to conform to RSS2 date formats
		var (
			dt time.Time
		)
		if pubDate == "" {
			dt = time.Now()
		} else {
			dt, err = time.Parse(`2006-01-02`, pubDate)
			if err != nil {
				return err
			}
		}
		pubDate = dt.Format(time.RFC1123)
		item := new(rss2.Item)
		item.GUID = articleURL
		item.Title = title
		item.Author = author
		item.PubDate = pubDate
		item.Link = u.String()
		item.Description = description
		feed.ItemList = append(feed.ItemList, *item)
		return nil
	})
	return err
}
