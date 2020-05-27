package mkpage

import (
	"os"
	"path"
	"testing"
	"time"

	// Caltech Library Packages
	"github.com/caltechlibrary/rss2"
)

var (
	today,
	prefix,
	blogJSON,
	fName string
)

func TestBlogMetaToRSS(t *testing.T) {
	meta := new(BlogMeta)
	if err := meta.BlogIt(prefix, fName, today); err != nil {
		t.Errorf("Could not BlotIt(%q, %q, %q), %s", prefix, fName, today, err)
		t.FailNow()
	}
	if err := meta.Save(blogJSON); err != nil {
		t.Errorf("Could not Save(%q), %s", blogJSON, err)
		t.FailNow()
	}
	feed := new(rss2.RSS2)
	if err := BlogMetaToRSS(meta, feed); err != nil {
		t.Errorf("Count not BlogMetaToRSS(meta, feed), %s", err)
		t.FailNow()
	}
	if len(feed.ItemList) != 1 {
		t.Errorf("Expected a single entry in feed.ItemList, got %+v\n", feed.ItemList)
		t.FailNow()
	}
	if len(feed.ItemList[0].Description) == 0 {
		t.Errorf("Expected a description feed.ItemList[0].Description, got %+v\n", feed.ItemList[0])
		t.FailNow()
	}
}

func TestWalkRSS(t *testing.T) {
	feed := new(rss2.RSS2)
	if err := WalkRSS(feed, prefix, "", "", "", ""); err != nil {
		t.Errorf(`Expected WalkRSS(feed, %q, "", "", "", "") err to be nil, got %s`, prefix, err)
		t.FailNow()
	}
}

func TestMain(m *testing.M) {
	today = time.Now().Format("2006-01-02")
	prefix = path.Join("test", "blog")
	os.MkdirAll(prefix, 0777)
	blogJSON = path.Join(prefix, "blog.json")
	fName = "README.md"
}
