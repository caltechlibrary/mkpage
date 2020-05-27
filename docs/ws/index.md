
# ws

## USAGE

	ws [OPTIONS] [DOCROOT]

## DESCRIPTION

	a nimble web server

ws is a command line utility for developing and testing static websites.
It uses Go's standard http libraries and can supports both http 1 and 2
out of the box.  It is intended as a minimal wrapper for Go's standard
http libraries supporting http/https versions 1 and 2 out of the box.

### CONFIGURATION

ws can be configurated through environment settings. The following are
supported.

+ MKPAGE_URL  - sets the URL to listen on (e.g. http://localhost:8000)
+ MKPAGE_DOCROOT - sets the document path to use
+ MKPAGE_SSL_KEY - the path to the SSL key if using https
+ MKPAGE_SSL_CERT - the path to the SSL cert if using https



## ENVIRONMENT

Environment variables can be overridden by corresponding options

```
    MKPAGE_DOCROOT   # set the htdoc root
    MKPAGE_SSL_CERT  # set the path to the SSL Certificate
    MKPAGE_SSL_KEY   # set the path to the SSL KEY
    MKPAGE_URL       # set the URL to listen on, defaults to http://localhost:8000
```

## OPTIONS

Below are a set of options available. Options will override any corresponding environment settings.

```
    -c                  Set the path for the SSL Cert
    -cert               Set the path for the SSL Cert
    -cors-origin        Set the CORS Origin Policy to a specific host or *
    -d                  Set the htdocs path
    -docs               Set the htdocs path
    -example            display example(s)
    -generate-manpage   generate man page
    -generate-markdown  generate markdown documentation
    -h                  display help
    -help               display help
    -k                  Set the path for the SSL Key
    -key                Set the path for the SSL Key
    -l                  display license
    -license            display license
    -quiet              suppress error messages
    -redirects-csv      Use target,destination replacement paths defined in CSV file
    -u                  The protocol and hostname listen for as a URL
    -url                The protocol and hostname listen for as a URL
    -v                  display version
    -version            display version
```


## EXAMPLE

Run web server using the content in the current directory
(assumes the environment variables MKPAGE_DOCROOT are not defined).

```shell
   ws
```

Run web server using a specified directory

```shell
   ws /www/htdocs
```


