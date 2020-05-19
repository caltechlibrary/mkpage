{
    "markup": "mmark",
    "has_code": true
}

# Make slides with mkpage

## USAGE

    mkpage --mkslides -i MARKDOWN_FILENAME [OPTIONS] [KEY/VALUE DATA PAIRS] [TEMPLATE_FILENAMES]

## SYNOPSIS

_mkslides_ converts a Markdown file into a sequence of HTML5 slides using the
key/value pairs to populate the templates and render to stdout.

## Features

+ Use Markdown to write your presentation in one file
+ Separate slides by "--" and a new line (e.g. \n versus \r\n)
+ Apply the default template or use your own
+ Control Layout and display with HTML5 and CSS

_mkslides_ option is based on _mkpage_ with the difference that 
multiple pages result from a single Markdown file. To manage the linkage 
between slides some predefined template variables is used.

+ *title* which would hold the page title for presentation
+ *csspath* which would hold the path to your CSS File.
+ *content* holds the extracted for each slide
+ *cur_no* which holds the current page number
+ *first_no* which holds the first slide's page number (e.g. 00)
+ *last_no* which holds the last slides page number (e..g length of slide deck minus one)
+ *prev_no* which holds the previous slide number if CurNo is create than 0
+ *next_no* which holds the next slide number if CurNo is not the last slide
+ *filename* is the filename for presentation
+ *nav* is to replace the default presentation navigation
+ *header* is to add a header to the presentation
+ *footer* is to add a footer to the presentation

In your custom templates these should be exist to link everything together
as expected.  In addition you may want to include JavaScript to allow mapping
actions like "next slide" to the space bar or mourse click.

## CONFIGURATION

+ MKPAGE_TEMPLATES - specify where to find the template(s) to use for slides

## EXAMPLE

In this example we're using the default slide template.
Here's an example of a three slide presentation

```

    Welcome to mkpage used to make slides.

    by R. S. Doiel, <rsdoiel@library.caltech.edu>

    --

    # mkslides

    _mkslides_ can generate multiple HTML5 pages from
    one markdown file.  It splits the markdown file
    on each "--" 

    --

    Thank you

    Hope you enjoy [mkpage](https://github.com/caltechlbrary/mkpage)

```

If you saved this as presentation.md you can run the following
command to generate slides

```shell
    mkpage --mkslides -i presentation.md \
        "title=text:My Presentation" \
	    "csspath=text:css/slides.css"
```

Using a custom template would look like

```shell
    mkpage --mkslides -i presentation.md \
        -t custom-slides.tmpl \
        "title=text:My Presentation" \
	    "csspath=text:css/slides.css"
```

This would result in the following webpages

+ 00-presentation.html
+ 01-presentation.html
+ 02-presentation.html

