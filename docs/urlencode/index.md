
# urlencode

## USAGE

	urlencode [OPTIONS] [STRING_TO_ENCODE]

## DESCRIPTION

urlencode is a simple command line utility to URL encode content. By default
it reads from standard input and writes to standard out.  You can
also specifty the string to encode as a command line parameter.


## OPTIONS

Below are a set of options available.

```
    -examples            display example(s)
    -generate-manpage    generate man page
    -generate-markdown   generate markdown documentation
    -h, -help            display help
    -i, -input           set input filename
    -l, -license         display license
    -nl, -newline        add a trailing newline to output
    -o, -output          set output filename
    -q, -query           use query escape (pluses for spaces)
    -quiet               suppress error messages
    -v, -version         display version
```


## EXAMPLE

```shell
    echo "This is the string to encode & nothing else!" | urlencode
```

would yield

```
    This%20is%20the%20string%20to%20encode%20&%20nothing%20else%0A
```


