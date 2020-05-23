[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)

# mkpage

_mkpage_ is a collection of tools that form a static website renderer.
It was inspired by deconstructing more complex content management
systems and distilling the rendering functions down to a core set
of command line tools.  It is well suited to building sites hosted 
on services like GitHub Pages or Amazon's S3. It is comprised of a 
set of command line utilities that augment the standard suite of 
Unix/POSIX commands available on most POSIX based operating systems 
(e.g. Linux, Mac OS X, Raspberry Pi and Windows systems that 
have a port of Bash). It uses the widely adopted 
[Pandoc](https://pandoc.org) as its markup conversion tool 
(e.g. Markdown to HTML). As such you can create your website from
Markdown, ReStructureText or Jira text.

_mkpage_ can run on machines as small as a Raspberry Pi.  Its small 
foot print and minimal dependencies (only Pandoc) means installation 
usually boils down to copying the precompiled binaries to a bin directory 
in your path after a installing Pandoc. Precompiled binaries are 
available for Linux, Windows and Mac OS X running on Intel as well as 
for the ARM7 versions of Raspbian running on Raspberry Pi.  _mkpage_
supports [Pandoc](https://pandoc.org/MANUAL.html#templates)'s template 
language.  This language is easy to learn and well documented. It is 
generally easier to use than more ambitious template engines like 
[Jekyll](https://jekyllrb.com/), [Hugo](https://gohugo.io)'s Go 
templates and [Assemble](http://assemble.io/).

_mkpage_'s minimalism is an advantage when you combine _mkpage_ 
with the standard suite of text processing tools available under your 
typical Unix/POSIX like operating systems. This makes scripting a _mkpage_ 
project using languages like Python, Make or Bash straight forward.  
Each _mkpage_ utility is independent. You can use as few or as many 
or as few as you like. You determine the workflow and build process 
that best fits your needs.


The following command line tools come with _mkpage_ 

+ [mkpage](docs/mkpage.html) -- a single page renderer and processor for Markdown, [Fountain](https://fountain.io), ReStructureText, Textile, Jira markup, JSON using [Pandoc](https://pandoc.org) as convertion and template engine
+ [mkrss](docs/mkrss.html) -- an RSS feed generator for content authored in Markdown and rendered to HTML
+ [sitemapper](docs/sitemapper.html) -- an XML Sitemap generator
+ [frontmatter](docs/frontmatter.html) -- a front matter extractor
+ [byline](docs/byline.html) -- a tool for extracting bylines from Markdown files and front matter
+ [titleline](docs/titleline.html) -- a tool for extracting the first title (H1) in a Markdown document or from front matter
+ [reldocpath](docs/reldocpath.html) -- a relative path calculator, useful for pathing hrefs and src attributes in a website
+ [ws](docs/ws.html) -- a fast, small, web server for site development or deployment

## A quick tour

_mkpage_ command accepts key/value pairs. The pairs can be explicit data 
types, [front matter] in files or resources from the web.
_mkpage_ assembles the metadata and content and sends them along to 
Pandoc for processing. In this example and the template 
is implemented as a Pandoc's template language. 

The "key" in our key/value pairs is used to map into the 
[Pandoc](https://pandoc.org/MANUAL.html) templates you want rendered. 
If a key was called "content" the template element would be like 
`${content}`.  The value of "content" would replace `${content}` in
the Pandoc template.  Pandoc templates are combine logic and iteration 
to make more complex pages.

On the "value" side of the key/value pair you have strings of one of 
several formats - plain text, markdown, [fountain](https://fountain.io),
ReStructureText, Jira text and JSON. The values can be from 
explicit strings associated with a data type, data from a file where the
file extension identifies the content type, or 
content retrieved via a URL based on the mime-type sent from the web 
service.  Here's a basic demonstration of sampling of capabilities
and integrating data from the [NOAA weather website](http://weather.gov).


### A basic template

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

### typing the _mkpage_ command

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
Likewise the weather data is coming from a URL idendifitied by the 
"http:" protocol reference . *mkpage* uses the "protocol" 
prefix to distinguish between literals, file paths and URL based 
based content. "http:" and "https:" returns an HTTP header 
the header is used to identify the content type for processing by
_mkpage_ before handing off to Pandoc. E.g. "Content-Type: text/markdown" 
tells us to use Pandoc to translate from Markdown to HTML. For data 
contained in files we rely on the file extension to identify content 
type, e.g. ".md" is markdown, ".rst" is ReStructureText, ".json" is a 
JSON document.  If no content type is decernable then we assume the 
content is plain text.

### the utilities

#### mkpage

*mkpage* is a page renderer.  It comes with some helper utilities 
that make scripting a deconstructed content management system from 
Python/Bash easier.

#### mkrss

*mkrss* will scan a directory tree for Markdown files and add each 
markdown file with a corresponding HTML file to the RSS feed generated.

#### frontmatter

*frontmatter* will extract a Markdown files' front matter so you can
process it with another tool. When you used in conjunction with *mkpage*
you can render the same file into metadata about the file and 
HTML output. This is handy if you're using the front matter to build
up metadata in an HTML template or building a corpus JSON document
for use with browser side search engines like [Lunrjs](https://lunrjs.com).

#### byline

*byline* will look inside a markdown file and return the first _byline_ fromfront matter or one idendified by a regular expression. 
It returns an empty string if it finds none. The default regular
expression fallback can be overridden with a command line option.

#### titleline

*titleline* will look inside a markdown file's front matter or 
the return the first h1 equivalent title it finds or an empty string 
if it finds none. 

#### reldocpath

*reldocpath* is intended to simplify the calculation of relative
asset paths (e.g. common CSS files, images, feeds) when working from
a common project directory.

#### blogit

*blogit* performs two tasks, first is given a filename and date 
(in YYYY-MM-DD format) it will copy the file into an appropriate
blog path based on the date provided. The second task is it
maintains a `blog.json` file describing the content of the blog.
This is placed in the same folder as the where the year folder for
the blog is create.

##### Example

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

#### ws

*ws* is a simple static file web server.  It is suitable for viewing your 
local copy of your static website on your machine.  It runs with minimal 
resources and by default will serve content out to the URL 
http://localhost:8000.  It can also be used to host a static website 
and has run well on small Amazon virtual machines as well as Raspberry Pi
computers acting as local private network web servers.

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

