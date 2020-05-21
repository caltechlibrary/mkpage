{
    "markup": "markdown",
    "title": "mkpage, a possible future",
    "byline": "By R. S. Doiel, 2020-05-19",
    "creator": [ "R. S. Doiel" ],
    "date": "2020-05-19"
}


# mkpage, a possible future

This is a distillation of my thinking about Markdown processors and **mkpage** as static website generator.  I am excited about the current state of [Pandoc](https://pandoc.org/) because of Mike Hucka's [Pangolin](https://github.com/mhucka/pangolin-notebook).

## What's Pandoc?

Pandoc is a document conversion tool. It converts to and from a ridiculous number of formats.  In addition to HTML it includes output formats I want to use in the near future like ePub2, ePub3, PDF. Pandoc has excellent support Mathematical text.  Pandoc is [R Studio](https://rstudio.com/)'s Markdown to HTML engine. This is important as that means people who use R Studio to create web content will know Pandoc's dialect of Markdown even if they don't know they are using Pandoc. Pandoc can be extended using [Lua](https://lua.org). Pandoc itself is written in [Haskell](https://www.haskell.org/) which is an interesting language.

On a selfish note, Pandoc supports the features of Markdown I regularly use, code blocks, tables and foot/end notes and Pandoc will convert from Markdown to Jira/Confluence's flavor of markup.  

### So why talk up Pandoc?

Based on Mike's Pangolin my thinking around Pandoc has evolved. There is Pandoc the document converter but there is also Pandoc the template engine.

Pandoc as a template engine is intriguing. I was able to read and understand Mike's templates without resorting to documentation. More importantly **Pandoc has good documentation**, including their simple template language!  This holds promise for future proofing our site generation with **mkpage**.

### mkpage's failure

A huge glaring failure in my `mkpage` project has been my reliance on Go's native template language. It was good bootstrap but not so fun now. 
While Go's template language is "enterprise" ready[^ready] Go's template language remains poorly documented. Go's template language has little outside adoption aside from Hugo[^Hugo]. This situation has not improved in the last decade and I think is unlikely to improve.  Go's various Markdown to HTML libraries have caused concerns. Like what I observed in the NodeJS community many Go packages fail to be updated with version changes and have become orphaned projects. My concern is that overtime 3rd Party Go packages I've relied on will loose development inertia or get abandoned. I am on my third Markdown library in Go at this stage of `mkpage`. Maybe it's time to step off the tread mill?

[^ready]:  In the sense of Java's various "enterprise" ready template languages

[^Hugo]:  Hugo starts with Go's template language but extends it into it's own dialect. In other words a niche in a niche. See https://gohugo.io/.


Integrating Pandoc into `mkpage` can be achieved initially using `os.exec` along with pipes. Pandoc has the advantage of being a standard Markdown to HTML converter. It has a simpler and much better documented template engine. Pandoc solves two issues in near term evolution of `mkpage`. I don't believe Pandoc's template language is well know but it is easy to pickup and is much better documented than the Go's template language. If Pandoc had its current feature set existed when I originally started `mkpage` I probably would not have created `mkpage`. I would have only focused on making things scriptable and metadata friendly which is where I think `mkpage` shines.


Where does this leave `mkpage` project? At a really interesting crossroads (at least to me).  I think `mkpage` collection of tools solves a number of things that are not addressed by Pandoc.

1. `mkpage` has decent front matter support for processing metadata related to Markdown documents
2. `mkpage` command line language (parameter language) is well suited to scripting and is cleaner than a traditional Unix options syntax for that purpose
3. `mkpage` as a set of tools, provides many additional features such as RSS generation, a crude Sitemap generator, title and byline extractors which assists in constructing many kindas of websites.

`mkpage` evolved from a different purpose than Pandoc and that purpose is still relevant. Content Management Systems  (CMS) whether custom, off the shelf, commercial or open sources tend to succumbed to excessive features, complexity and resource consumption.  `mkpage` remains a reflective reaction to the challenges and brittleness of those complex systems. Simplicity, at least at the tool level, is a virtue.

`mkpage` is an Un-CMS, a distillation of the rendering features of complex CMS.  It is a set of tools for leveraging metadata provided as JSON documents, command line options and front matter in plain text documents.  Like Pandoc it is highly scriptable from languages like Bash, Python and Make. With the option of switching template languages or content processors it remains a useful but simple abstraction in website/content rendering. When you know your "content structure" it is easy automate site construction with simple scripts.

## Where does **mkpage** go from here?

**mkpage** is adopting Pandoc has its Markdown to HTML engine. **mkpage** is adding Pandoc as a template engine as an option with an eye to it becoming the default template engine.  This will allow mkpage to be useful for other markup syntaxes like Textile and RestructText which Pandoc supports. Pandoc integration should be complete for the v0.1.0 release of **mkpage**.



