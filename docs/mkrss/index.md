
USAGE
=====

	mkrss [OPTIONS] HTDOCS [RSS_FILENAME]

DESCRIPTION
-----------


mkrss walks the file system to generate a RSS2 file. It assumes 
that the directory for HTDOCS is is the base directory containing 
subdirectories in the form of /YYYY/MM/DD/ARTICLE_HTML where 
YYYY/MM/DD (Year, Month, Day) corresponds to the publication date 
of ARTICLE_HTML.


OPTIONS
-------

Below are a set of options available.

```
    -b, -byline            set byline regexp
    -channel-builddate     Build Date for channel (e.g. 2006-01-02 15:04:05 -0700)
    -channel-category      category for channel
    -channel-copyright     Copyright for channel
    -channel-description   Description of channel
    -channel-generator     Name of RSS generator
    -channel-language      Language, e.g. en-ca
    -channel-link          link to channel
    -channel-pubdate       Pub Date for channel (e.g. 2006-01-02 15:04:05 -0700)
    -channel-title         Title of channel
    -d, -date-format       set date regexp
    -e                     A colon delimited list of path exclusions
    -examples              display example(s)
    -generate-markdown     generate markdown documentation
    -h, -help              display help
    -i, -input             set input filename
    -l, -license           display license
    -o, -output            set output filename
    -quiet                 suppress error messages
    -t, -title             set title regexp
    -v, -version           display version
```


EXAMPLES
--------


If our htdocs folder is our document root and out blog is
htdocs/myblog.

    mkrss -channel-title="This Great Beyond" \
        -channel-description="Blog to save the world" \
        -channel-link="http://blog.example.org" \
        htdocs htdocs/rss.xml

This would build an RSS 2 file in htdocs/rss.xml from the
articles in htdocs/myblog/YYYY/MM/DD.


mkrss v0.2.4
