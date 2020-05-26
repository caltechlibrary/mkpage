{
    "has_code": true
}


# One Element

One element features a theme with one template, [page.tmpl](page.tmpl), which has
one replacable element named "Content". 

```template
    <!DOCTYPE html>
    <html>
    <head>
        <title>One Element</title>
        <link rel="stylesheet" href="/css/site.css">
    </head>
    <body>
        <header>
            <h1>One Element<h1>
        </header>
        <nav>
            <ul>
                <li><a href="/">Home</a></li>
                <li><a href="../">Up</a></li>
            </ul>
        </nav>
        {{with .Content}}<section>{{- . -}}</section>{{- end}}
        <footer>This template features a single replacable element</footer>
    </body>
    </html>
```

To build this one template site we can use a Bash script.
This example will assembling markdown files into HTML pages. The Bash
script is called [mk-website.bash](mk-website.bash).

```shell
    #!/bin/bash

    START="$(pwd)"
    cd "$(dirname "$0")"

    function softwareCheck() {
    	for NAME in "$@"; do
    		APP_NAME="$(which "$NAME")"
    		if [ "$APP_NAME" = "" ] && [ ! -f "./bin/$NAME" ]; then
    			echo "Missing $NAME"
    			exit 1
    		fi
    	done
    }

    echo "Checking necessary software is installed"
    softwareCheck mkpage ws
    echo "Converting Markdown files to HTML"
    for MARKDOWN_FILE in $(find . -type f | grep -E "\.md$"); do
    	# Calculate the HTML filename
    	HTML_FILE="$(dirname "$MARKDOWN_FILE")/$(basename "$MARKDOWN_FILE" .md).html"
    	mkpage \
    		"Content=$MARKDOWN_FILE" \
    		page.tmpl >"$HTML_FILE"

    done

    cd "$START"
```

To test this theme do the following run the following commands in this directory.

```shell
    ./mk-website.bash
    ws
```

Then point your webbrowser at http://localhost:8000 and view this page.

## Limitations

1. This theme assumes this directory is the root HTML directory
2. No unified navigation beyond what you provide in your Markdown files is available.


