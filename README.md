[![Project Status: Active – The project has reached a stable, usable state and is being actively developed.](https://www.repostatus.org/badges/latest/active.svg)](https://www.repostatus.org/#active)

# mkpage

_mkpage_ is a deconstructed, post modern, content management system for generating static websites.
It is suited to building sites hosted on services like GitHub Pages or Amazon's S3. It is 
comprised of a set of command line utilities that augment the standard suite of Unix/POSIX commands 
available on most POSIX based operating systems (e.g. Linux, Mac OS X, Raspberry Pi and Windows systems that 
have a port of Bash).

_mkpage_ can run on machines as small as a Raspberry Pi.  Its small foot print and minimal 
dependencies means installation usually boils down to copying the precompiled binaries to a bin directory 
in your path. Precompiled binaries are available for Linux, Windows and Mac OS X running on Intel as 
well as for the ARM7 versions of Raspbian running on Raspberry Pi.  _mkpage_ is built on Go's text templates.
The template markup is similar to the [Mustache](https://mustache.github.io/) and 
[Handlebars](http://handlebarsjs.com/). It is also similar to [Hugo](https://gohugo.io)'s template markup.   _mkpage_ has been easier for us to support when compared with 
more established static site generators like [Jekyll](https://jekyllrb.com/) [Hugo](https://gohugo.io) and [Assemble](http://assemble.io/).

_mkpage_'s minimalism turns into an advantage when you combine _mkpage_ with the standard suite of text processing tools available under your typical Unix/POSIX like operating systems. This makes scripting a _mkpage_ project using languages like Bash and Python relatively straight forward.  Each _mkpage_ utility is independent. You can use as few or as many as you like when you script your website creation process. You wind up with a workflow and build process that best fits your needs.


The following command line tools come with _mkpage_ 

+ [mkpage](docs/mkpage.html) -- a single page renderer with support for Markdown, [Mmark](https://mmark.nl),  [Fountain](https://fountain.io), JSON and Go text templates
+ [mkslides](docs/mkslides.html) -- a HTML slide generator based on the approach in _mkpage_
+ [mkrss](docs/mkrss.html) -- an RSS feed generator for content authored in Markdown and rendered to HTML
+ [sitemapper](docs/sitemapper.html) -- an XML Sitemap generator
+ [frontmatter](docs/frontmatter.html) -- a front matter extractor
+ [byline](docs/byline.html) -- a tool for extracting bylines from Markdown files
+ [titleline](docs/titleline.html) -- a tool for extracting the first title (H1) in a Markdown document
+ [reldocpath](docs/reldocpath.html) -- a relative path calculator, useful for pathing hrefs and src attributes in a website
+ [ws](docs/ws.html) -- a fast, small, web server for site development or deployment
+ [mkpongo](docs/mkpongo.html) -- a experimental version of mkpage using [Pongo2](https://github.com/flosch/pongo2) style templates.

## A quick tour

_mkpage_ command accepts key value pairs and applies them to a Golang [text/template](https://golang.org/pkg/text/template/).  
The key side of a pair corresponds to the template element names that will be replaced in the render 
version of the document. If a key was called "Content" the template element would look like `{{ .Content }}`. 
The value of "Content" would replace `{{ .Content }}`. Go text/templates elements can do more than 
that but the is the core idea.  On the value side of the key/value pair you have strings of one of 
four formats - plain text, markdown, [fountain](https://fountain.io) and JSON.  These four formatted strings can be explicit strings, 
data from a file or content retrieved via a URL.  Here's a basic demonstration of sampling of capabilities
and integrating data from the [NOAA weather website](http://weather.gov).

### a basic template

```template
    {{ define "weather.tmpl" }}
    Date: {{- .now}}
    
    Hello {{.name -}},
        
    The current weather is
    
    {{index .weatherForecast.data.weather 0}}
    
    Thank you
    
    {{.signature}}
    
    {{ end }}
```

To render the template above (i.e. [weather.tmpl](examples/weather.tmpl)) is expecting values from various data sources.
This break down is as follows.

+ "now" and "name" are explicit strings
    + "now" integrates getting data from the Unix _date_ command
+ "weatherForecast" comes from a URL which returns a JSON document
    + ".data.weather" is the path into the JSON document
    + _index_ is a function that lets us pull out the initial value in the array
+ "signature" comes from a file in our local disc

### the _mkpage_ command

Here is how we would express the key/value pairs on the command line.

```shell
    mkpage "now=text:$(date)" \
        "name=text:Little Frieda" \
        "weather=http://forecast.weather.gov/MapClick.php?lat=13.47190933300044&lon=144.74977715100056&FcstType=json" \
        signature=examples/signature.txt \
        examples/weather.tmpl
```

Notice the two explicit strings are prefixed with "text:" (other possible formats are "markdown:", "json:").
Values without a prefix are assumed to be file paths. We see that in testdata/signature.txt.  Likewise the 
weather data is coming from a URL. *mkpage* distinguishes that by the prefixes "http://" and "https://". 
Since a HTTP response contains headers describing the content type (e.g.  "Content-Type: text/markdown") we 
do not require any other prefix. Likewise a filename's extension can give us an inference of the data format 
it contains. ".json" is a JSON document, ".md" is a Markdown document and everything else is just plain text.


Since we are leveraging Go's [text/template](https://golang.org/pkg/text/template/) the template itself
can be more than a simple substitution. It can contain conditional expressions, ranges for data and even
include blocks from other templates.



### About Go text/template

*mkpage* template engine is the Go [text/template](https://golang.org/pkg/text/template/) package.  You can 
get a feel for working with Go templates and _mkpage_ by exploring _mkpage_'s [How To](how-to/). A good place
to start is [how to/the basics](how-to/the-basics.html) and then proceed to [How To/One element](how-to/one-element/).


### companion utilities

#### mkpage

*mkpage* comes with some helper utilities that make scripting a deconstructed
content management system from Bash easier.

#### mkslides

*mkslides* generates a set of HTML5 slides from a single Markdown file. It uses
the same template engine as *mkpage*

#### mkrss

*mkrss* will scan a directory tree for Markdown files and add each markdown file with
a corresponding HTML file to the RSS feed generated.

#### frontmatter

*frontmatter* will extract a Markdown files' front matter so you can
process it with another tool. When you used in conjunction with *mkpage*
you can render the same file into metadata about the file and 
HTML output. This is handy if you're using the front matter to build
up metadata in an HTML template or building a corpus JSON document
for use with browser side search engines like [Lunrjs](https://lunrjs.com).

#### byline

*byline* will look inside a markdown file and return the first _byline_ it finds
or an empty string if it finds none. The _byline_ is identified with a regular
expression. This regular expression can be overridden with a command line option.

#### titleline

*titleline* will look inside a markdown file and return the first h1 equivalent title
it finds or an empty string if it finds none. 

#### reldocpath

*reldocpath* is intended to simplify the calculation of relative
asset paths (e.g. common CSS files, images, feeds) when working from
a common project directory.

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

*ws* is a simple static file web server.  It is suitable for viewing your local copy
of your static website on your machine.  It runs with minimal resources and by default
will serve content out to the URL http://localhost:8000.  It can also be used to host
a static website and has run well on small Amazon virtual machines as well as Raspberry Pi
computers acting as local private network web servers.

##### Example

```shell
    ws Sites/mysite.example.org
```

This would start the web server up listen for browser requests on _http://localhost:8000_.
The content viewable by your web browser would be the files inside the _Sites/mysite.example.org_
directory.

```shell
    ws -url http://mysite.example.org:80 Sites/mysite.example.org
```

Assume the machine where you are running *ws* has the name mysite.example.org then your could
point your web browser at _http://mysite.example.org_ and see the web content you have in 
_Site/mysite.example.org_ directory.
