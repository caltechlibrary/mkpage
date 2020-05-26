{
    "has_code": true
}


# urldecode

## USAGE

    urldecode [OPTIONS] [STRING_TO_ENCODE]

## SYNOPSIS

urldecode is a simple command line utility to URL decode content. By default
it reads from standard input and writes to standard out.  You can
also specifty the string to decode as a command line parameter.

## OPTIONS

```
	-h	display help
	-help	display help
	-i	set input filename
	-input	set input filename
	-l	display license
	-license	display license
	-o	set output filename
	-output	set output filename
	-v	display version
	-version	display version
```

## EXAMPLES

```shell
    echo 'This%20is%20the%20string%20to%20encode%20&%20nothing%20else%0A' | urldecode
```

would yield

```shell
    This is the string to encode & nothing else!
```

