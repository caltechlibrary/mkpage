{
    "has_code": true
}


# titleline

## USAGE

    titleline [OPTIONS]

## SYNOPSIS

titleline extracts the first title line from a Markdown file. By default it reads
from standard in and writes to standard out but can read/write
to specific files using an option.

## OPTIONS

```
	-h	display help
	-help	display help
	-i	input filename
	-input	input filename
	-l	display license
	-license	display license
	-o	output filename
	-output	output filename
	-t	set title regexp
	-title	set title regexp
	-v	display version
	-version	display version
```

## EXAMPLE

```shell
    cat article.md | titleline
```

This will display the title of an article.md.

