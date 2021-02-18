
USAGE
=====

	ws [OPTIONS] [DOCROOT]

DESCRIPTION
-----------



	a nimble web server

ws is a command line utility for developing and testing static websites.
It uses Go's standard http libraries and can supports both http 1 and 2
out of the box.  It is intended as a minimal wrapper for Go's standard
http libraries supporting http/https versions 1 and 2 out of the box.


OPTIONS
-------

Below are a set of options available.

```
    -c, -cert            Set the path for the SSL Cert
    -cors-origin         Set the CORS Origin Policy to a specific host or *
    -d, -docs            Set the htdocs path
    -example             display example(s)
    -generate-markdown   generate markdown documentation
    -h                   display help
    -help                display help
    -k, -key             Set the path for the SSL Key
    -l                   display license
    -license             display license
    -quiet               suppress error messages
    -redirects-csv       Use target,destination replacement paths defined in CSV file
    -u, -url             The protocol and hostname listen for as a URL
    -v                   display version
    -version             display version
```


EXAMPLES
--------


Run web server using the content in the current directory

   ws

Run web server using a specified directory

   ws /www/htdocs


ws v0.2.4
