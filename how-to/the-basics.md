{
    "has_code": true
}


# The Basics

_mkpage_ uses Pandoc's template engine to render content. This template 
is documented on the [Pandoc](https://pandoc.org/MANUAL.html#templates).
It is a simple template system compared to systems like [Jykell]() and
[Go text/templates]().

## Basic element

Like a templated data element is wrapped in a `${` .. `}` or in
`$` .. `$`.  A template that says "Hello" followed by a value
named "world" would look something like this

```
    Hello ${world}
```

We can use the template to say Hello to "Georgina"

```shell
    echo 'Hello ${world}' > hello-world.tmpl
    mkpage "world=text:Georgina" hello-world.tmpl
```

Running these two command should result in output like

```
    Hello Georgina
```

The line with the `echo` is just creating our template and saving it 
as the file _hello-world.tmpl_.  In the template the only special part 
is `${world}`. "world" is a variable will be replaced by something we 
define.  In the line with `mkpage` we define the value for "world" as 
plain text. Note we don't need to prefix "world" in the command line. 
The 'text:' before the variable name to indicate the type of object, 
in this case plain text.  The line tells the template to replace 
`${world}` with the text value "Georgina".  The last part of the 
command instructs _mkpage_ to use our _hello-world.tmpl_ template.

If we did not include `world=...` with the _mkpage_ command using 
the _hello-world.tmpl_ template _mkpage_ would return output like 

```
    Hello
```

If we included other key/value pairs not mentioned in the template 
they would be silently ignored. 

If we made a typo in _hello-world.tmpl_ then we would see an error 
message from the [pandoc template engine](https://pandoc.org/MANUAL.html#exit-codes). NOTE: "5	PandocTemplateError" is the one that tells you
you have a template error.


Try the following to get a feel for how key/value pairs work with 
_mkpage_. The first two will render but display `Hello`. The first 
example fails because no value is provided and the second fails 
because the value provided doesn't match what is in the template 
(notice the typo "wrold" vs. "world").  the next one will display an 
error because _text:_ wasn't included on the value side of the 
key/value pair.  By default _mkpage_ assumes the value is refering 
to a file and in this case can't find the file Georgina in your 
current directory.  The last two will display `Hello Georgina` 
should display since the value for "world" is provided. The last one 
just ignores "Name=text:Fred" because "Name" isn't found in the template.

```shell
    mkpage hello-world.tmpl
    mkpage "wrold=text:Georgina" hello-world.tmpl
    mkpage "world=Georgina" hello-world.tmpl
    mkpage "world=text:Georgina" "name=text:Fred" hello-world.tmpl
    mkpage "world=text:Georgina" hello-world.tmpl
```


### Conditional elements

Pandoc templates have conditional elements.  `${if()} ... ${endif}` 
and `${for()} ... ${endfor}` are simimlar.  The "if" will display
the contents if the variable passing to the "if" function is non-empty.
If the object is an array it will concatenate the elements together,
if it is a map then it will display true. The "for" function is similar
but iterates over the elements provided and those elements can be accessed
with an alias of `${it}`.  See the Pandoc Manual for details.    

```
   ${if(title)}Title: ${title}${endif}
```

or using "for"

```
    By ${for(authors)}${it}$sep$, ${endfor}
```

Let's create a template file with both these statements called _title-demo.tmpl_ and run the 
_mkpage_ command.

```shell
    echo '${if(title)Title: ${title}${endif}' > title-demo.tmpl
    echo 'By ${for(authors)}${it}$sep$ and ${endif}' >> title-demo.tmpl 
    mkpage 'title=text:This is a title demo' \
           'authors=json:[ 'Jane', 'Carol' ]' \
           title-demo.tmpl
```

The output should look like

```
    Title: This is a title demo
    By Jane and Carol
```

In the first line with the *if* we use "title" as the variable, 
just like "world" in our first example. The second line uses
a *for* to iterate over authors names and insert an "and" between
them.  What happens if you run this command?

```shell
    mkpage title-demo.tmpl
```

This produces two lines of output that look something 
. The reason we don't see something like

```
    
    By
```

We see the "By" because it is outside the *for* loop, but we don't
see a title because "Title: " is inside the *if* block.


### Template and sub templates

Pandoc provides for a simple mechanism to include sub templates.
If the same directory as your main template any template can be
included as its own function. This is easier to see in practice
than to describe.

Create two templates, "header.tmpl" and "footer.tmpl".

This would be the header.

```
   <header>${if(title)}<h1>${title}</h1>${endif}</header>
```

This would be for the footer.

```
    <footer>By ${for(authors)}${it}$sep$, ${endfor}</footer>
```

Now create a third temaplte called "main.tmpl"

```
   ${header()}
   <section>
   This is the title: ${title}
   <p>
   Authors are: ${authors[, ]}
   </section>
   ${footer()}
```

Now run our *mkpage* command providing our "title" and "authors" values
and the template named "main.tmpl".

```
    mkpage 'title=text:This is a big template project' \
           'authors=text:Me, Myself and I' \
           main.tmpl
```

This output should look like

```html
    <header><h1>Hi there</h1></header><section>
    This is the title: Hi there
    <p>
    Authors are: Me, Myself and I
    </section>
    <footer>
```

In this case "header.tmpl" and "footer.tmpl" are found because
they are in the same directory as "main.tmpl".


## Content formats and data sources

*mkpage* understands several content formats

+ text/plain (e.g. "text:" strings and any file expect those having the extension ".md" or ".json")
+ text/markdown (e.g. "markdown:" strings and file extension ".md")
+ application/json (e.g. "json:" strings and file extension ".json")
+ text/restructured (e.g. "rst:" strings and file extension ".rst")
+ text/textile (e.g. "textile:" strings and file extension ".textile")
+ text/fountain (e.g. "fountain:" strings and file extension ".fountain")

It also supports three data sources

+ an explicit string (prefixed with a format, e.g. "text:", "markdown:", "json:")
+ a filepath and filename (the default data source)
+ a URL (identified by the URL prefixes http:// and https://)

Content type is evaluated, transformed (if necessary), and sent to the Pandoc template engine.

```
    This is a plain text string: "${string}"

    Below is a an included file:
    ${file}
    
    Finally below is data from a URL:
    ${url}
    ${end}
```

Create a text file named _hello.md_.

```
    # this is a file

    Hello World!
```

Type the following

```shell
    mkpage "string=text:Hello World" "file=hello.md" \
      "url=https://raw.githubusercontent.com/caltechlibrary/mkpage/master/nav.md" \
      data-source-demo.tmpl
```

What do you see?

## A note about Markdown dialect

_mkpage_ uses [Pandoc](https://pandoc.org) as its markdown (markup)
and template (output) engine.

The markdown processor is invoked for values with the "markdown:" 
hint prefix, files ending in ".md" extension or URL content with the 
content type returned as "text/markdown" (i.e.  content with a type 
of "text/plain" does not use the markdown process and is treated as 
plain text).

