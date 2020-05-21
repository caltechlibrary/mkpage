
# Action items

## Bugs

+ [ ] **sitemapper** needs to respect the 50K/50MB url and size limits per spec, see https://www.sitemaps.org/protocol.html

## Next (road to v0.1.0)

+ [x] Switch to Pandoc for conversion for markup to HTML as well as for template support
+ [ ] Figure out how to co-mingle Markdown, Fountain, safely 
+ [x] **byline** should pickup a by line from front matter OR the regexp
+ [x] **titleline** should pickup a title from front matter OR the regexp
+ [x] **mkslides** should be depreciated, pandoc provides slide generation
+ [ ] **sitemapper** should consider front matter in deciding the structure of sitemap.xml, also should allow for more than once sitemap.xml to be generated (E.g. a blog might have its own sitemap, see https://www.sitemaps.org/protocol.html
+ [ ] Document mkpage front matter metadata practices based on codemeta.json and relavant Scheme.org scheme, Pandoc can take advantage of most of this.
    + [ ] `.doi` the DOI associated with a page
    + [ ] `.creator` should be an array of creator info (e.g. ORCID, given_name, family_name)
    + [ ] `.title`
    + [ ] `.date`
    + [ ] `.publishDate`
    + [ ] `.lastmod`
    + [ ] `.description`
    + [ ] `.draft` (bool)
    + [ ] `.keywords`
    + [ ] `.linkTitle`
    + [ ] `.markup` (e.g. markdown, fountain, maybe remarkjs)
    + [ ] `.series`
    + [ ] `.issue`
    + [ ] `.volume`
    + [ ] `.no`
    + [ ] `.slug`
    + [ ] `.type` (e.g. post, article, homepage)
    + [ ] `.permalink`  (e.g. resolver URL)
    + [ ] `.language`
+ [ ] mkpage Sitemap support
    + Current sitemap cli is too naive for sites more than a couple dozen pages
    + Need to support possibly nested sitemap XML references
    + Need some sort of front matter to identify where/if content would show up in sitemap

## Someday, Maybe

+ [ ] Consider extending mkpage's parameter language to include filters on JSON data, e.g. `myvar=json(.[0:1]):[ "one", "two", "three" "four"]` would return the first (zeroth) element for two elements.

