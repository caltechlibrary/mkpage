
# USAGE

>	titleline \[OPTIONS\]

## DESCRIPTION

titleline extracts the first title line from a Markdown file. 
By default it reads from standard in and writes to standard out 
but can read/write to specific files using an option.


## OPTIONS

Below are a set of options available.

```
    -examples            display example(s)
    -generate-markdown   generate markdown documentation
    -h, -help            display help
    -i, -input           input filename
    -l, -license         display license
    -o, -output          output filename
    -quiet               suppress error messages
    -t, -title           set title regexp
    -v, -version         display version
```


## EXAMPLES

```
    cat article.md | titleline
```

This will display the title of an article.md.


titleline v0.1.1

