
# One Element

One element features a theme where the template, [page.tmpl](page.tmpl), has
one replacable element named "content". Included is a python script for 
assembling markdown files into HTML pages. The python
script is called [mk_website.py](mk-website.py).

To test this theme do the following run the following commands in this directory.

```shell
    ./mk_website.py
    ws
```

Then point your webbrowser at http://localhost:8000 and view this page.

## Limitations

1. This theme assumes this directory is the root HTML directory
2. No unified navigation beyond what you provide in your Markdown files is available.


