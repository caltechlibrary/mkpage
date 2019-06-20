---
{
    "has_code": true
}
---

# Simple Theme

This theme demonstates the replacement of three content elements in the
template. Two are explicit text lines and one like the one element theme
is a Markdown file.

This theme supports using a common Title element and CSSPath element across
all the pages in the website. The [mk-website.bash](mk-website.bash) will 
traverse all the Markdown files and render corresponding HTML pages.

This theme relies on three _mkpage_ project commands - _mkpage_, 
_reldocpath_ and _ws_ (for testing the website and viewing from your web 
browser over http://localhost:8000)


To test this theme do the following run the following commands in this 
directory.

```shell
    export WEBSITE_TITLE="Simple Theme Demo"
    ./mk-website.bash
    ws
```

Point your webbrowser at http://localhost:8000 and view this page.

### Template example

```template
    <!DOCTYPE html>
    <html>
    <head>
        {{with .title -}}<title>{{- . -}}</title>{{- end }}
        {{with .csspath -}}<link rel="stylesheet" href="{{- . -}}">{{- end}}
        {{with .css -}}<style>{{- . -}}</style>{{- end }} 
    </head>
    <body>
        <header>
            {{with .title -}}<h1>{{- . -}}</h1>{{- end}}
        </header>
        <nav>
            <ul>
                <li><a href="/">Home</a></li>
                <li><a href="../">Up</a></li>
            </ul>
        </nav>
        {{with .content }}<section>{{- . -}}</section>{{- end}}
        <footer>Simple is a theme that works with  three elements Title, CSSPath, and Content</footer>
    </body>
    </html>
```

### Bash script

```shell
    #!/bin/bash

    START="$(pwd)"
    cd "$(dirname "$0")"

    function SoftwareCheck() {
    	for NAME in "$@"; do
    		APP_NAME="$(which "$NAME")"
    		if [ "$APP_NAME" = "" ] && [ ! -f "./bin/$NAME" ]; then
    			echo "Missing $NAME"
    			exit 1
    		fi
    	done
    }

    echo "Checking necessary software is installed"
    SoftwareCheck mkpage reldocpath ws
    if [ "$WEBSITE_TITLE" = "" ]; then
    	WEBSITE_TITLE="Simple Theme Demo"
    fi

    echo "Converting Markdown files to HTML supporting a relative document path to the CSS file"
    for MARKDOWN_FILE in $(find . -type f | grep -E "\.md$"); do
    	# Caltechlate DOCPath
    	DOCPath="$(dirname "$MARKDOWN_FILE")"
    	# Calculate the HTML filename
    	HTML_FILE="$DOCPath/$(basename "$MARKDOWN_FILE" .md).html"
    	CSSPath="$(reldocpath "$DOCPath" css)"
    	mkpage \
    		"title=text:${WEBSITE_TITLE}" \
    		"csspath=text:${CSSPath}/site.css" \
    		"content=${MARKDOWN_FILE}" \
    		page.tmpl >"${HTML_FILE}"
    done

    cd "$START"
```


## Improvements over one-element

The "title" value can be set for the whole site by modifying by setting an
environment variable WEBSITE_TITLE.

The "csspath" (CSS file path) is calculate with _reldocpath_. This means that you could
place content rendered with this them in a subdirectory of a larger website 
and still use the CSS that comes with this theme.

## Limitations

1. This theme assumes this directory is the root HTML directory
2. No unified navigation beyond what you provide in your Markdown files is available.



