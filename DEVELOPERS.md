
# Notes on compile _mkpage_

## Requirements

+ Pandoc needs to be installed, it is used to convert markup formats to HTMLand as the template engine for HTML output.
+ go 1.14 or better
+ _make_ if you want to the _Makefile_ to build the project.
+ [pkgassets](https://github.com/caltechlibrary/pkgassets) for generating a new _assets.go_
+ Caltech Library Go Packages
    + github.com/caltechlibrary/cli
    + github.com/caltechlibrary/tmplfn
    + github.com/caltechlibrary/rss2

## Compiling from source

Using _go get_

```shell
    go get -u github.com/caltechlibrary/pkgassets/...
    go get -u github.com/caltechlibrary/mkpage/...
```

Manual using only the go command

```shell
    for PNAME in byline mkpage mkrss mkslides reldocpath sitemapper titleline urldecode urlencode ws; do
        go build -o "bin/${PNAME}" "cmds/${PNAME}/${PNAME}.go"
    done
```

### regenerating assets.go

_assets.go_ holds the go source code for a map containing the contents of the _defaults_ directory (e.g.
templates/page.tmpl and templates/slides.tmpl). If you modify those files you'll need to recreate
_assets.go_. You can do so with the [pkgassets](https://github.com/caltechlibrary/pkgassets) tool.

```shell
    pkgassets -o assets.go -p mkpage Defaults defaults
```

If you're not modifying the contents of the defaults directory you do not need to regenerate _assets.go_.

