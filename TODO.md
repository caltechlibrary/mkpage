
# Action items

## Bugs

+ [ ] ws, `*.mjs` should be served as "text/javascript" by default, `.mjs` is the JS file extension for JS Modules some sites use
+ [ ] **sitemapper** needs to respect the 50K/50MB url and size limits per spec, see https://www.sitemaps.org/protocol.html

## Next (road to v1.0.0)

+ [x] Figure out how to support all pandoc conversion formats
+ [ ] Better error messaging (preferrably passing through Pandoc's own error message) for non-0 exit values from Pandoc
+ [x] Remove support for Go templates, remove pkgassets dependency
+ [ ] Decide what do about `sitemappper`, depreciate in favor of other tools or improve to support large sites (e.g. feeds.library.caltech.edu)
    + Current sitemap cli is too naive for sites more than a couple dozen pages
    + Need to support possibly nested sitemap XML references
    + Need some sort of front matter to identify where/if content would show up in sitemap
    + If `blog.json` available take advantage of that metadata in `sitemappe`

## Someday, Maybe

+ [ ] Figure out how to co-mingle Markdown, Fountain, safely 
+ [ ] Consider extending mkpage's parameter language to include filters on JSON data, e.g. `myvar=json(.[0:1]):[ "one", "two", "three" "four"]` would return the first (zeroth) element for two elements.
+ [ ] consider creating a `wikit` cli for wiki like static sites
    + a `wiki.json` could be used to generate a sitemap file
+ [ ] Pandoc template examples organized as themes targetting research community
+ [ ] Carpentry style tutorials for developing sites with MkPage Project

