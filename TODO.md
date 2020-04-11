
# Action items

## Bugs

+ [ ] **sitemapper** needs to respect the 50K/50MB url and size limits per spec, see https://www.sitemaps.org/protocol.html

## Next (road to v1.0.0)

+ [ ] **byline** should pickup a by line from front matter OR the regexp
+ [ ] **titleline** should pickup a title from front matter OR the regexp
+ [ ] **mkslides** should be depreciated in favor of **mkpage** using front matter to indicate an output format of slides.
+ [ ] **sitemapper** should consider front matter in deciding the structure of sitemap.xml, also should allow for more than once sitemap.xml to be generated (E.g. a blog might have its own sitemap, see https://www.sitemaps.org/protocol.html
+ [ ] Read in mkpage.toml, mkpage.json or mkpage.yaml for mkpage config
+ [ ] Add support for rendering remarkjs content
+ [ ] Add support for passing configuration to markup engine from front matter
+ [ ] Figure out how to co-mingle Markdown, Fountain, remarkjs safely 
+ [ ] mkpage front matter based on library metadata practices, codemeta.json and relavant Scheme.org scheme
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
    + [ ] `.markdup` (e.g. markdown, fountain, maybe remarkjs)
    + [ ] `.series`
    + [ ] `.slug`
    + [ ] `.type` (e.g. post, article, homepage)
    + [ ] `.permalink`  (e.g. resolver URL)
    + [ ] `.language`
    + [ ] `.remarkjs` holds the settings for our remarkjs engine
+ [ ] mkpage Sitemap support
    + Current sitemap cli is too naive for sites more than a couple dozen pages
    + Need to support possibly nested sitemap XML references
    + Review Hugo's sitemap support for ideas
    + Need some sort of front matter to identify where/if content would show up in sitemap
+ [ ] mkpage slide support needs to align with remarkjs
    + See https://github.com/gnab/remark, https://remarkjs.com and https://github.com/gnab/remark/wiki
    + [ ] Consider merging mkpage and mkslide (fewer tools less to learn), consider front matter changes
    + [ ] Add support for slide notes delimited by `???`

## Someday, Maybe

+ [ ] Should **mkpage** continue to support Toml and Yaml frontmatter?
+ [ ] Should **mkpage** support Django 3/Jinja2 style templates?
+ [ ] Should **mkpage** support [RMarkdown](https://rmarkdown.rstudio.com/)?
    + See [RMarkdown, Definitive Guide](https://bookdown.org/yihui/rmarkdown/)
    + Requirements
        + Use of exec to pass content to RMarkdown via separate process
        + an R with RMarkdown installed
+ [ ] Should **mkpage** populate a `.Page` variable with page level metadata?
    + See https://gohugo.io/variables/page/ for definitions in Hugo
    + [ ] `.Alaises` aliases to this page (need to clarify this with mkpage's approach)
    + [ ] `.Content` the content itself defined after the front matter
    + [ ] `.Data` data specific to this type of page
    + [ ] `.Date` 
    + [ ] `.Description`
    + [ ] `.Dir`
    + [ ] `.Draft`
    + [ ] `.ExpiryDate`
    + [ ] `.File` see `.File` for sub-fields
    + [ ] `.FuzzyWordCount` (how fuzzy?)
    + [ ] `.MkPage` see `.MkPage` for sub-fields
    + [ ] `.IsHome`
    + [ ] `.IsNode`
    + [ ] `.IsPage` always true (why include it?)
    + [ ] `.IsTranslated` 
    + [ ] `.Keywords`
    + [ ] `.Kind`
    + [ ] `.Language`
    + [ ] `.Lastmod`
    + [ ] `.LinkTitle`
    + [ ] `.Next`
    + [ ] `.NextInSection`
    + [ ] `.OutputFormats`
    + [ ] `.Pages`
    + [ ] `.Permalink`
    + [ ] `.Plain` strip HTML tags (e.g. what you might want to index with Lunr)
    + [ ] `.PlainWords` content stripped for HTML tags and returned as a slice of string ([]string)
    + [ ] `.PreviousInSection`
    + [ ] `.PublishDate`
    + [ ] `.RawContent` (page content without front matter)
    + [ ] `.ReadingTime`
    + [ ] `.Resources`
    + [ ] `.Ref`
    + [ ] `.RelPermalink`
    + [ ] `.Summary`
    + [ ] `.TableOfContents`
    + [ ] `.Title`
    + [ ] `.Translations`
    + [ ] `.TranslationKey`
    + [ ] `.Truncated`
    + [ ] `.Type`
    + [ ] `.UniqueID` an MD5 Checksum of the file path (relative to content root)
    + [ ] `.Weight`
    + [ ] `.WordCount`
+ [ ] Should **mkpage** populate a `.File` variable with file related metadata
    + [ ] `.File.Path` (the path relative to the current working directory, content root)
    + [ ] `.File.LogicalName` (e.g. foo.en.md)
    + [ ] `.File.TranslationBaseName` (e.g. path.Base() in Go minus any language extension (e.g. foo from foo.en.md))
    + [ ] `.File.ContentBaseName` (e.g. same as TranslationalBaseName for compatibility with Hugo, not sure about Leaf bundle support)
    + [ ] `.File.BaseFileName` (e.g. path.Base(), without extention, e.g. foo.en.md -> foo.en)
    + [ ] `.File.Ext` (e.g. path.Ext())
    + [ ] `.File.Lang` (look at how Hugo handles this, probably uses the two letter language code)
    + [ ] `.File.Dir` (e.g. path.Dir() from current working directory, content root)
+ [ ] Should **mkpage** populate a set of `.MkPage` variables accessible to template or another varname populated and with option to override?
    + [ ] `.MkPage.Version` (semver)
    + [ ] `.MkPage.Generator` (generator string like used in RSS 2)
+ [ ] Add a tool to generate and search lunr indexes
+ [ ] Add support for metadata taken from Namaste (Name as text) in the directory
+ [ ] Remove the default template and ship distribution with a set of standard templates
+ [ ] Add support for integration with pygments for syntax highlighting
+ Investigate moving beyond Go templates 
    + [ ] review [Mustache](https://github.com/hoisie/mustache) or Handlebars via [velvet](https://github.com/gobuffalo/velvet) or [raymond](https://github.com/aymerick/raymond)
    + [ ] review [Pango2](https://github.com/flosch/pongo2) templates
+ [ ] Configurable Mime-Type assignments
+ [ ] Add a general purpose indexer that can process both Markdown files and metadata in JSON documents with same name (e.g. page.md, page.json)
+ [ ] Adds a csvblock that reads in a CSV file and converts to a GFM table like _csv2mdtable_
+ [ ] titleline and byline should use filenames on the command line if provided (not require -i)

