
# Release Notes

## v0.2.1

+ Compiled now depends on v0.0.5 of wsfm package, this is used in the ws command

## v0.2.0

+ mkpage assumes pandoc v2.10.x or better
+ mkpage nolonger support Golang templates
+ `-t, -to` now set the "to" value passed to pandoc v2.10.x or better
+ `-f, -from` now sets the "from" value passed to pandoc v2.10.x or better
+ removed mkslides command line program, Pandoc replaces that functionality

## v0.1.2

+ Switched default pandoc parsing mode to commonmark to limit the false positive attempts to process TeX in markdown documents
+ Added json-generator type to allow processing of JSON content by external scripts (e.g. Python), script should return JSON via stdout
+ Added support to include more explicit pandoc markdown interpretation, e.g. commonmark, markdown_strict, gfm

## v0.1.1

+ To use Go language templates you must use the `-t` or `-templates` option
+ Fixed bug where JSON documents were parsed as frontmatter
+ Added pandoc info to `-version` documents
+ Removed support for environment variable in tools
+ Fully dropped let's encrypt support (available from other none project tools)


## v0.1.0a

+ Documentation cleanup, code cleanup, some bug fixes
+ If v0.1.0a proves stable consider v1.0.0-rc1
+ Go templates depreciated, a option is required to access go template processing require

## v0.0.33

+ Finish Pandoc integration
+ Added `blogit` tool for creating blog styles paths and maintaining a `blog.json` file describing blog
+ Removed Pongo2 integration
+ Go templates are now depreciated, they can still be used using `-pandoc=false`

## v0.0.32i

+ mkslides has be removed
+ default templates are Pandoc templates
+ Pandoc is **required** as it provides the markup convertion to HTML and the default template engine

## v0.0.26

+ Compiled with Go's 1.12 templates supporting variable creation and substution
+ *frontmatter* command line tool for extracting the Hugo/Rmarkdown front matter from a Markdown file so you can process it separately, by default frontmatter reads from standard in and writes to stand out so you can use it as a datasource in mkpage
+ *mkpage* now skips over front matter like that used in Hugo and Rmarkdown 
+ Let's Encrypt support removed


## v0.0.18

+ Templates are now all assumed to start with a define with the master template listed first and matching its basename
    + Affects _mkpage_, _mkslides_
+ Various bug fixes
    + Fixed some CORS handling in _ws_
+ Added ACME cert support for https in _ws_
    + Any http request will automatically redirect to https when ACME cert support is enabled
+ mkpage, mkslide accept templates via stdin, normalize key/value data pairs between them, updated mkslide docs
+ Bug fixes
    + Improved documentation
    + Fixed sitemapper bug when mapping from the current working directory
    + Update copyright link for Caltech Library

