{
    "has_code": true
}


# mkrss

## USAGE

    mkrss [OPTION] HTDOCS [RSS_FILENAME]

## SYNOPSIS

mkrss walks the file system to generate a RSS2 file. It assumes 
that the directory for HTDOCS is is the base directory containing 
subdirectories in the form of /YYYY/MM/DD/ARTICLE_HTML where 
YYYY/MM/DD (Year, Month, Day) corresponds to the publication date 
of ARTICLE_HTML.

## OPTIONS

```
	-b	set byline regexp
	-byline	set byline regexp
	-c	If non-zero, limit the number of articles in the RSS file
	-channel-builddate	Build Date for channel (e.g. 2006-01-02 15:04:05 -0700)
	-channel-category	category for channel
	-channel-copyright	Copyright for channel
	-channel-description	Description of channel
	-channel-generator	Name of RSS generator
	-channel-language	Language, e.g. en-ca
	-channel-link	link to channel
	-channel-pubdate	Pub Date for channel (e.g. 2006-01-02 15:04:05 -0700)
	-channel-title	Title of channel
	-d	set date regexp
	-date-format	set date regexp
	-e	A colon delimited list of path exclusions
	-h	display help
	-help	display help
	-i	set input filename
	-input	set input filename
	-l	display license
	-license	display license
	-o	set output filename
	-output	set output filename
	-t	set title regexp
	-title	set title regexp
	-v	display version
	-version	display version
```

## EXAMPLE

If our htdocs folder is our document root and out blog is
htdocs/myblog.

```shell
    mkrss -channel-title="This Great Beyond" \
        -channel-description="Blog to save the world" \
        -channel-link="http://blog.example.org" \
        htdocs htdocs/rss.xml
```

This would build an RSS 2 file in htdocs/rss.xml from the
articles in htdocs/myblog/YYYY/MM/DD.

