{
    "has_code": true
}


# sitemapper

## USAGE

    sitemapper [OPTIONS] HTDOCS_PATH MAP_FILENAME PUBLIC_BASE_URL

## OVERVIEW

sitemapper generates a sitemap for the website.

## OPTIONS

```
	-e	A colon delimited list of path parts to exclude from sitemap
	-exclude	A colon delimited list of path parts to exclude from sitemap
	-h	display help
	-help	display help
	-l	display license
	-license	display license
	-u	Set the change frequencely value, e.g. daily, weekly, monthly
	-update-frequency	Set the change frequencely value, e.g. daily, weekly, monthly
	-v	display version
	-version	display version
```

## EXAMPLE

```
    sitemapper htdocs htdocs/sitemap.xml http://eprints.example.edu
```

Renders out a new sitemap at htdocs/sitemap.xml

