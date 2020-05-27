
# frontmatter

## USAGE

	frontmatter [OPTIONS]

## DESCRIPTION


frontmatter extracts a front matter from a Markdown file. If no front matter is present then an empty file is returned. Note frontmatter doesn’t process the data extracted. It returns it unprocessed. Other tools can be used to process the front matter appropriately. By default frontmatter reads from standard in and writes to standard out. This makes it very suitable for pipeline processing or for passing JSON formatted front matter back to mkpage for integration into the templates processed.


## OPTIONS

Below are a set of options available.

```
    -examples           display example(s)
    -generate-manpage   generate man page
    -generate-markdown  generate Markdown documentation
    -h, -help           display help
    -i, -input          input filename
    -j, -json           output as JSON
    -l, -license        display license
    -o, -output         output filename
    -quiet              suppress error messages
    -v, -version        display version
```


## EXAMPLES


Extract a front matter from article.md.

```shell
    cat article.md | frontmatter
```

This will display the front matter if found in article.md.

```shell
    frontmatter -i article.md
```

Will also do the same.


