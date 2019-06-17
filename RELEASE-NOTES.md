
# Release Notes

## v0.0.26

+ Compiled with Go's 1.12 templates supporting variable creation and substution
+ *frontmatter* command line tool for extracting the Hugo/Rmarkdown front matter from a Markdown file so you can process it separately, by default frontmatter reads from standard in and writes to stand out so you can use it as a datasource in mkpage
+ *mkpage* now skips over front matter like that used in Hugo and Rmarkdown 


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

