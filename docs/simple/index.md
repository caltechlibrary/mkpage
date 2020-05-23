
# Simple Theme

This theme demonstates the replacement of three content elements in the
template. Two are explicit text lines and one like the one element theme
is a Markdown file.

This theme supports using a common Title element and CSSPath element across
all the pages in the website. The [mk_website.py](mk_website.py) will 
traverse all the Markdown files and render corresponding HTML pages.

This theme relies on three _mkpage_ project commands - _mkpage_, 
_reldocpath_ and _ws_ (for testing the website and viewing from your web 
browser over http://localhost:8000)


To test this theme do the following run the following commands in this 
directory.

```shell
    export WEBSITE_TITLE="Simple Theme Demo"
    ./mk_website.py
    ws
```

Point your webbrowser at http://localhost:8000 and view this page.

## Improvements over one-element

The Title value can be set for the whole site by modifying by setting an
environment variable WEBSITE_TITLE.

The CSS file path is calculate with _reldocpath_. This means that you could
place content rendered with this them in a subdirectory of a larger website 
and still use the CSS that comes with this theme.

## Limitations

1. This theme assumes this directory is the root HTML directory
2. No unified navigation beyond what you provide in your Markdown files is available.



