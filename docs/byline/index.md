
USAGE
=====

	byline [OPTIONS]

DESCRIPTION
-----------


byline extracts a byline from a Markdown file. By default it reads
from standard in and writes to standard out but can read/write
to specific files using an option.


OPTIONS
-------

Below are a set of options available.

```
    -b, -byline          set byline regexp
    -examples            display example(s)
    -generate-markdown   generate Markdown documentation
    -h, -help            display help
    -i, -input           input filename
    -l, -license         display license
    -o, -output          output filename
    -quiet               suppress error messages
    -v, -version         display version
```


EXAMPLES
--------


Extract a byline from article.md.

    cat article.md | byline

This will display the byline if one is found in article.md.


byline v0.2.4
