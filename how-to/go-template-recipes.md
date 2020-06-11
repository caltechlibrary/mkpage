
# Go text/template recipes

*mkpage* template engine is Go's [text/template](https://golang.org/pkg/text/template/). Go's templates provide a flexible and simple [DSL](https://en.wikipedia.org/wiki/Domain-specific_language) describing how to assemble a document based on a data structure passed to it.  *mkpage* uses a list of key/value pairs on the command line to populate the data structure the template package expects.  This includes support for JSON formatted text from strings, files and URL response. It also support transforming markdown content into HTML before assembling the final template.


While Go's template package is not complicated to use it doesn't come with allot of examples or tutorials.  Most articles you find on Go's template packages either focus on web server code or are for sophisticated static content generators like [Hugo](http://gohugo.io). Hugo extends Go's template DSL providing capabilities that rival and surpass older static content generators like [Jekyll](https://jekyllrb.com/) and [Jade](http://jade-lang.com/).

*mkpage* uses Go v1.8's text/template as is providing little in the way of extensions.  *mkpage* is meant to be a trivially easy system for producing simple content from plain text, markdown text, and JSON. It deliberately implements a minimal feature set targetting a scripting environment like a Bash shell.



## only three data formats are supported

*mkpage* supports three document formats

+ text/plain
+ text/markdown
+ application/json


## only three data sources are supported

*mkpage* supports three data sources 

+ files (the default data source)
+ explicit strings (using a hint prefix, e.g. "text:", "markdown:", "json:")
+ URLs as data sources (prefixed with http:// and https:// as appropriate)

## Examples

### Rendering a Markdown as HTML

This is a minimal example of using *mkpage* to render a Markdown file as an HTML page.
It features no navigation, just a wrapping HTML document head (with link to CSS file) and body.

```
    {{ define "page.tmpl" }}
    <!DOCTYPE html>
    <html>
        <head><link rel="stylesheet" href="/css/site.css"></head>
        <body>
        {{ .Content }}
        </body>
    </html>
    {{ end }}
```

Rendering a markdown document named _myfile.md_ as _myfile.html_ would look like

```shell
    mkpage -pandoc=false Content=myfile.md page.tmpl > myfile.html
```



### explicit stings, a get well card

In this example we want to add a name to a simple get well message.

Our template is called **get-well.tmpl**. It looks like 

```go
    {{ define "get-well.tmpl" }}
    Dear {{ .name -}},

    Hope you are feeling better today.

    Sencerly,

    Mojo Sam
    
    {{ end }}
```

On the command line we can run *mkpage* with the following options

```shell
    mkpage -pandoc=false "name=text:Little Frieda" get-well.tmpl
```

The output would look like

```text
    Little Frieda,

    Hope you are feeling better today.

    Sencerly,

    Mojo Sam

```

#### Explanation

The key "name" has a string value of "Little Frieda".  The template indicates this needs to be included after the word "Dear". The key "name" is proceeded by a period or dot.  The substitution happens between the opening "{{" and closing "}}".  Notice the "-" before the closing "}}". This tells the template engine to not allow spacas after the value and the next non-space character (i.e. the comma of the opening line).

### JSON data, a key/value blob report

In this example we construct a JSON object as part of the key/value pairs on the command line and pass it through the blob.tmpl template that displays they pairs.

The command envokation looks like

```shell
    mkpage -pandoc=false 'blob=json:{"one":1,"two":2}'  blob.tmpl
```

The template is a simple range construct

```go
    {{ define "blob.tmpl" }}
    {{range $key,$val := .blob }}
        Key: {{ $key }} Value: {{ $val -}}
    {{end}}
    {{ end }}
```

Results in text like

```text
    
       Key: one Value: 1
       Key: two Value: 2

```

#### Explanation

We use the range function to iterate over the key/value pairs of our JSON object. Additionally we assign those values to the template variables called "`$key`" and "`$val`". These are then used to format our output. Also notice the trailing values "-" which supresses and extra new line.

## Files are data source

### Wraping a Markdown document in HTML

In this example we want to embed a "story" in a simple HTML document. The *story* is written in Markdown format. Here's the simple template

```go
    {{ define "simple-page.tmpl" }}
    <!DOCTYPE html>
    <html>
        <head><title>Stories</title></head>
        <body>
        {{ .story }}
        </body>
    </html>
    {{ end }}
```

The command line would look something like

```shell
    mkpage "story=my-story.md" simple-page.tmpl > my-story.html
```

#### Explanation

On the command line *story* is assumed to point to a file named "my-story.md". The reason a file is assumed is because there is no hint prefix or URL prefix at the start of the value. Because the file ends in the file extension ".md" it is assume to be a Markdown file and processed accordingly before being assemble in the template.


## URL as data source

### JSON data, a weather forecast

In this example we get the current weather forecast for Guam.  The source of the weather information is [NOAA](http://noaa.gov)'s [National Weather Services](http://weather.gov) website.  By including the parameter "FcstType=json" at the end of the URL you get a JSON version of the weather forecast rather than the HTML or XML alternatives.

+ data source: http://forecast.weather.gov/MapClick.php?lat=13.47190933300044&lon=144.74977715100056&FcstType=json

Our template will be call **forecast.tmpl**. It will be used to produce a Markdown file of weather related information obtained from the JSON response.

```go
    {{ define "forecast.tmpl" }}
    {{with $co := .forecast.currentobservation}}
    Current Observation:

        + {{ $co.name }}
        + Elevation: {{ $co.elev }}
        + Latitude: {{ $co.latitude }}
        + Longitude: {{ $co.longitude }}
        + Date: {{ $co.Date }}
        + Temp: {{ $co.Temp }}
        + Dew Point: {{ $co.Dewp }}
        + Relative Humidity: {{ $co.Relh }}
        + Winds: {{ $co.Winds }}
        + Wind direction: {{ $co.Windd }}
        + Gust: {{ $co.Gust }}
        + Visibility: {{ $co.Visibility }}

    {{end}}

    Forecast:
    {{range .forecast.data.text }}
        + {{ . -}}
    {{end}}

    {{ end }}
```

The command line for *mkpage* would look like

```shell
    mkpage "forecast=http://forecast.weather.gov/MapClick.php?lat=13.47190933300044&lon=144.74977715100056&FcstType=json" forecast.tmpl
```

The resulting page should look something like 

```text

    Current Observation:

        + Agana, Guam International Airport
        + Elevation: 299
        + Latitude: 13.48
        + Longitude: 144.8
        + Date: 5 Aug 08:54 am ChST
        + Temp: 82
        + Dew Point: 79
        + Relative Humidity: 89
        + Winds: 12
        + Wind direction: 220
        + Gust: 0
        + Visibility: 10.00


    Forecast:

        + Scattered showers and thunderstorms.  Mostly cloudy, with a high near 84. Breezy, with a southwest wind 23 to 25 mph, with gusts as high as 32 mph.  Chance of precipitation is 40%.
        + Scattered showers and thunderstorms.  Mostly cloudy, with a low around 79. Breezy, with a southwest wind 15 to 20 mph, with gusts as high as 25 mph.  Chance of precipitation is 40%.
        + Mostly cloudy, with a high near 88. Heat index values as high as 99. Breezy, with a southwest wind 17 to 21 mph, with gusts as high as 26 mph. 
        + Mostly cloudy, with a low around 79. Southwest wind 13 to 17 mph, with gusts as high as 22 mph. 
        + Mostly cloudy, with a high near 88. Southwest wind 14 to 17 mph, with gusts as high as 22 mph. 
        + Mostly cloudy, with a low around 79.
        + Mostly sunny, with a high near 89.
        + Partly cloudy, with a low around 80.
        + Scattered showers and thunderstorms.  Mostly cloudy, with a high near 89. Chance of precipitation is 40%.
        + Scattered showers and thunderstorms.  Mostly cloudy, with a low around 79. Chance of precipitation is 40%.
        + Scattered showers and thunderstorms.  Mostly cloudy, with a high near 89. Chance of precipitation is 40%.
        + Scattered showers and thunderstorms.  Mostly cloudy, with a low around 79. Chance of precipitation is 40%.
        + Partly sunny, with a high near 89.
```

## JSON with Dashes in property names

Some JSON APIs come with objects containing property names with dashes in them.  This is valid JSON but throws a monkey wrench into addressing them with Go template's usual dot notation. Go's templates provides a *index* function to address those types of property names.

```JSON
    [
        {
            "journal-title":"Favorite adventures", 
            "article-title": "Zamborra and Beyond",
            "author-list":[
                {"family-name": "Frieda", "other-name": "Little"},
                {"family-name": "Sam", "other-name": "Mojo"},
                {"family-name": "Flanders", "other-name": "Jack"}
            ]
        }
    ]
```

Iterating over this list in a template

```html
    <ul>
        {{with .data -}}
        <li>
           {{with (index "article-title" .)}}Article: {{ . }}{{end -}}
           {{with (index "journal-title" .)}} Journal: {{ . }}{{end -}}
           {{with (index "author-list" .)}}<br />
           Authors: {{range $i, $author := (index "author-list" .)}}
                    {{if gt $i 0}}; {{end-}}
                    {{index "family-name" $author}}, {{index "other-name" $author}}
                {{end}}
           {{end}}
        </li>
        {{- end}}
    </ul>
```

The *index* function can be used to build a nested path too

```
    {{index "top-level" "middle-level" "bottom-level" .someData}}
```

## Including sub-templates

In complex pages it is nice to be able to include sub templates. In this example
We have a _signature.tmpl_ and _postscript.tmpl_ files as sub templates to _letter.tmpl_.

#### letter.tmpl

```
    {{ define "letter.tmpl" }}
    Dear {{ .ToName }},

    Hope all is well.  I will be with you shortly though not necessarily
    on the same plane of existance.

    {{template "signature.tmpl" .}}

    {{template "postscript.tmpl" .}}
    {{ end }}
```

#### signature.tmpl

```
    {{ define "signature.tmpl" }}
    Sincerly,

    {{ .Name }}, somewhere next door to reality

    {{ end }}
```   

#### postscript.tmpl

```
    {{ define "postscript.tmpl" }}
    (P.S. What is comming at you is coming from you, {{rfc3339 "now"}})
    {{ end }}
```

#### Putting it all together

```
    mkpage "ToName=text:Mojo Sam" "Name=text:Jack Flanders" letter.tmpl signature.tmpl postscript.tmpl
```

This should output somthing like

```
    Dear Mojo Sam,

    Hope all is well.  I will be with you shortly though not necessarily
    on the same plane of existance.

    Sincerly,

    Jack Flanders, somewhere next door to reality

    (P.S. What is comming at you is coming from you, 2016-12-12T16:17:14-08:00)
```

