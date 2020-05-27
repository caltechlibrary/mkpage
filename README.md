[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)

# MkPage Project

**MkPage Project** is a collection of tools for rendering static websites.
Featured is the _mkpage_ command, a front end to
[Pandoc](https://pandoc.org). Pandoc supports converting from many
[lightweight markup languages](https://pandoc.org/ "Pandoc's list of supported formats"). _mkpage_ supports metadata encoded as
front matter[^frontmatter] using [YAML](https://yaml.org),
[TOML](https://github.com/toml-lang/toml) and
[JSON](https://www.json.org/json-en.html)
as well as additional data sources expressed on the command line
in a simple command language. Content is rendered using Pandoc's
[template language](https://pandoc.org/MANUAL.html#templates).


[^frontmatter]: Front matter in light weight markup languages like
Markdown start at the top of the file and begin and end with a simple
set of delimiters. `---` (three dashes) is used for YAML, `+++`
(three plus signs) for TOML, open and close curly braces are used by
JSON.

**MkPage Project** was inspired by deconstructing more complex
content management systems and distilling the rendering functions
down to a core set of simple command line tools.  It is well suited for
building sites hosted on services like GitHub Pages or Amazon's S3.
It is comprised of a set of command line utilities that augment the
standard suite of Unix/POSIX commands available on most POSIX based
operating systems for text processing (e.g. Linux, macOS,
Raspberry Pi and Windows systems that have a port of Bash). It uses
the widely adopted [Pandoc](https://pandoc.org) as its markup
conversion engine and template engine.  As such you can create your
website using a variety of light weight markup languages such as
[Markdown](https://daringfireball.net/projects/markdown/),
[Textile](http://redcloth.org/textile),
[ReStructureText](https://docutils.sourceforge.io/docs/ref/rst/introduction.html) and [Jira's](https://jira.atlassian.com/secure/WikiRendererHelpAction.jspa?section=all) wiki markup.


The _mkpage_ tools can run on machines as small as a Raspberry Pi.
Their small foot print and minimal dependencies (only Pandoc) means
installation usually boils down to copying the precompiled binaries
to a bin directory in your path after a installing Pandoc.
Precompiled binaries are available for Linux, Windows and macOS
running on Intel as well as for the ARM7 versions of Raspbian running
on Raspberry Pi.

The _mkpage_ command supports
[Pandoc](https://pandoc.org/MANUAL.html#templates)'s template language.
This language is easy to learn and well documented. It is
generally easier to use than more ambitious template engines like
[Jekyll](https://jekyllrb.com/) or [Hugo](https://gohugo.io).

_mkpage_'s minimalism is an advantage. It plays nice with
the standard suite of text processing tools available with
most Unix/POSIX compatible operating systems[^posix]. This makes
scripting a **MkPage Project** using languages like Python, Make or
Bash straight forward.  Each _mkpage_ utility is independent. You
can use as few or as many or as few as you like. You determine the
workflow and build process that best fits your needs.

[^posix]: Common POSIX compatible systems include macOS, Linux, and recent versions of Windows 10

## A quick tour _mkpage_ command

The _mkpage_ command accepts key/value pairs as command line parameters.
The pairs can be explicit data types, files on disc or resources from the
web. Additionally _mkpage_ will merge in any front matter found in
your light weight markup such as Markdown documents.
_mkpage_ assembles the all metadata into a JSON structure which will be
processed by Pandoc when rendering a Pandoc template.
Additionally _mkpage_ understands the [Fountain](https://fountain.io)
markup language will will handle conversion before passing to onto Pandoc.

### _mkpage_'s command language

The "key" in our key/value pairs is used to map into the
[Pandoc](https://pandoc.org/MANUAL.html) templates you want rendered.
If a key was called "content" the template element would be like
`${content}`.  The value of "content" would replace `${content}` in
the Pandoc template.  Pandoc templates can combine logic and iteration
to make more complex pages. See the Pandoc [User Guide](https://pandoc.org/MANUAL.html#templates) for more details.

On the "value" side of the key/value pair you have strings of one of
several formats - plain text, markdown, [fountain](https://fountain.io),
ReStructureText, Jira text and JSON. The values can be from
explicit strings associated with a data type, data from a file where the
file extension identifies the content type, or
content retrieved via a URL based on the mime-type sent from the web
service.  Here's a basic demonstration of sampling of capabilities
and integrating data from the [NOAA weather website](http://weather.gov).

Next is an example of a basic Pandoc template followed by an
example command line for invoking _mkpage_.

#### A basic template

```template

    Date: ${now}
    
    Hello ${name},
        
    The weather forecast is
    
    ${if(weather.data.weather)}
      ${weather.data.weather[; ]}
    ${endif}
    
    Thank you
    
    ${signature}

```

To render the template above (i.e. [weather.tmpl](examples/weather.tmpl))
is expecting values from various data sources broken down as follows.

+ "now" and "name" will be explicit strings
    + "now" integrates getting data from the Unix _date_ command
+ "weather" will come from a URL which returns a JSON document
    + ".data.weather" is the path into the JSON document
+ "signature" will come from a plain text file in your local disc

#### typing the _mkpage_ command

Here is how we would express the key/value pairs on the command line.

```shell
    mkpage "now=text:$(date)" \
        'name=text:Little Frieda' \
        'weather=http://forecast.weather.gov/MapClick.php?lat=13.47190933300044&lon=144.74977715100056&FcstType=json' \
        'signature=examples/signature.txt' \
        'examples/weather.tmpl'
```

Notice the two explicit strings are prefixed with "text:" (other formats
include "markdown:" and "json:").  Values without a prefix are assumed
to be local file paths. We see that in testdata/signature.txt is one.
Likewise the weather data is coming from a URL identified by the
"http:" protocol reference . *mkpage* uses the "protocol"
prefix to distinguish between literals, file paths and URL based
based content. "http:" and "https:" returns an HTTP header
the header is used to identify the content type for processing by
_mkpage_ before handing off to Pandoc. E.g. "Content-Type: text/markdown"
tells us to use Pandoc to translate from Markdown to HTML. For data
contained in files we rely on the file extension to identify content
type, e.g. ".md" is markdown, ".rst" is ReStructureText, ".json" is a
JSON document.  If no content type is discernible then we assume the
content is plain text.

### MkPage Project Tools

#### mkpage

[mkpage](docs/mkpage/) is a page renderer and front end to Pandoc.
It is used to aggregate metadata and templates into a complete website page.
It serves as a pre-processor for Markdown, [Fountain](https://fountain.io),
ReStructureText, Textile, Jira markup, JSON using
[Pandoc](https://pandoc.org) as conversion and template engine.

#### blogit

[blogit](docs/blogit/) performs two tasks, first if given a
filename and date (in YYYY-MM-DD format) blogit will copy the file
into an appropriate blog path based on the date provided. The second task
it performs is to maintain a `blog.json` file describing the content of
the blog.  This is placed in the same folder as the where the year
folders for the blog are create.

#### mkrss

[mkrss](docs/mkrss/) is an RSS feed generator for content authored
in Markdown.  It can read a `blog.json` file created with the _blogit_
and produce an RSS feed from it or scan the directory tree for Markdown
files with corresponding HTML files and generate an RSS feed.

#### frontmatter

[frontmatter](docs/frontmatter/) will extract a light weight markup
files' (e.g. Markdown) front matter so you can process it with another
tool. It can optional convert the front matter into JSON even if the
front matter was defined in YAML or TOML.

#### byline

[byline](docs/byline/) inside a Markdown file's front matter
for a "byline" field and return it before scanning the file using
a regular expression for the byline.  If nothing is found in either
the front matter or the regular expression then it'll return an empty
string.

#### titleline

[titleline](docs/titleline/) will look inside a markdown file's
front matter for a "title" field it if not present in the front matter
it'll look for the h1 demlimiter (Markdown '# ') and return it's content.
It will return an empty string if it finds none.


#### reldocpath

[reldocpath](docs/reldocpath/) is intended to simplify the
calculation of relative asset paths (e.g. common CSS files, images,
RSS files) when working from a common project directory.


##### Example reldocpath

You know the path from the source document to target document from the project root folder.

+ Source is *course/week/01/readings.html*
+ Target is *css/site.css*.

In Bash this would look like--

```shell
    # We know the paths relative to the project directory
    DOC_PATH="course/week/01/readings.html"
    CSS_PATH="css/site.css"
    echo $(reldocpath $DOC_PATH $CSS_PATH)
```

the output would look like

```shell
    ../../../css/site.css
```

#### sitemapper

[sitemapper](docs/sitemapper/) a simplistic XML Sitemap generator.
Sitemaps are used by web crawls to find content in your website and
can help your website be more search-able by modern full text search
engines.


#### ws

[ws](docs/ws/) is a simple static file web server.  It is suitable
for viewing your local copy of your static website on your machine.
It runs with minimal resources and by default will serve content out to
the URL http://localhost:8000. It is a fast, small, web server for
site development and copyedit work.

##### Example

```shell
    ws Sites/mysite.example.org
```

This would start the web server up listen for browser requests on
_http://localhost:8000_.  The content viewable by your web browser would
be the files inside the _Sites/mysite.example.org_ directory.

```shell
    ws -url http://mysite.example.org:80 Sites/mysite.example.org
```

Assume the machine where you are running *ws* has the name
mysite.example.org then your could point your web browser at
_http://mysite.example.org_ and see the web content you have in
_Site/mysite.example.org_ directory.

## Problem Reporting and lending a hand

**MkPage** project is hosted at [GitHub](https://github.com/caltechlibrary/mkpage) and bugs can be reported via the [Issue Tracker](https://github.com/caltechlibrary/mkpage/issues). As an open source project 
pull requests as well as bug reports are appreciated.

## Getting your copy of **MkPage Project**

You can find releases of **MkPage Project** at [github.com/caltechlibrary/mkpage](https://github.com/caltechlibrary/mkpage/releases)

## License

**MkPage Project** is released under an open source [license](license.html).

