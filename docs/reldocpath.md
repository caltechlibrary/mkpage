{
    "has_code": true
}


# reldocpath

## USAGE

    reldocpath SOURCE_DOC_PATH TARGET_DOC_PATH 

## SYNOPSIS

Given a source document path, a target document path calculate and
the implied common base path calculate the relative path for target.

## OPTIONS

```
	-h	display help
	-help	display help
	-l	display license
	-license	display license
	-v	display version
	-version	display version
```

## EXAMPLE

Given

```
    reldocpath chapter-01/lesson-03.html css/site.css
```

would output

```
    .../css/site.css
```

