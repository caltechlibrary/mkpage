
# USAGE

>	sitemapper \[OPTIONS\] HTDOCS_PATH MAP_FILENAME PUBLIC_BASE_URL

## DESCRIPTION

sitemapper generates a sitemap for the website.


## ENVIRONMENT

Environment variables can be overridden by corresponding options

```
    MKPAGE_DOCROOT   # set the document root, defaults to current working directory
    MKPAGE_SITEMAP   # set the sitemap filename and path
    MKPAGE_SITEURL   # set the site url
```

## OPTIONS

Below are a set of options available. Options will override any corresponding environment settings.

```
    -docs                        set the htdoc root
    -examples                    display example(s)
    -exclude                     A colon delimited list of path parts to exclude from sitemap
    -generate-markdown           generate markdown documentation
    -h, -help                    display help
    -l, -license                 display license
    -o, -output                  output filename (for logging)
    -quiet                       suppress error messages
    -sitemap                     set the sitemap filename and path
    -update, -update-frequency   Set the change frequencely value, e.g. daily, weekly, monthly
    -url                         set the site URL
    -v, -version                 display version
```


## EXAMPLES

```
    sitemapper htdocs htdocs/sitemap.xml http://eprints.example.edu
```

sitemapper v0.1.1

