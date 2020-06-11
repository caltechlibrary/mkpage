
# USAGE

>	reldocpath \[OPTIONS\] SOURCE_DOC_PATH TARGET_DOC_PATH

## DESCRIPTION

Given a source document path, a target document path calculate and
the implied common base path calculate the relative path for target.


## OPTIONS

Below are a set of options available.

```
    -examples            display example(s)
    -generate-markdown   generate markdown documentation
    -h, -help            display help
    -l, -license         display license
    -quiet               suppress error messages
    -v, -version         display version
```


## EXAMPLES

Given

```
    reldocpath chapter-01/lesson-03.html css/site.css
```

would output

```
    .../css/site.css
```

reldocpath v0.1.1

