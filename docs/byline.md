{
    "has_code": true
}


# byline

### USAGE

    byline [OPTIONS]

## SYNOPSIS

byline extracts a byline from a Markdown file. By default it reads
from standard in and writes to standard out but can read/write
to specific files using an option.

## OPTIONS

```
	-b	set byline regexp
	-byline	set byline regexp
	-h	display help
	-help	display help
	-i	input filename
	-input	input filename
	-l	display license
	-license	display license
	-o	output filename
	-output	output filename
	-v	display version
	-version	display version
```

## EXAMPLE

```shell
    cat article.md | byline
```

This will display the byline of article.md.

