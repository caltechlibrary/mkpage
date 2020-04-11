
# USAGE

	mkpongo [OPTIONS] [KEY/VALUE DATA PAIRS] TEMPLATE_FILENAME

## DESCRIPTION


SYNOPSIS

Using the key/value pairs populate the template(s) and render to stdout.


## OPTIONS

Below are a set of options available.

```
    -code               outout just code blocks for language, e.g. shell or json
    -codesnip           output just the code bocks
    -examples           display example(s)
    -generate-manpage   generate man page
    -generate-markdown  generate markdown documentation
    -h, -help           display help
    -i, -input          input filename
    -l, -license        display license
    -o, -output         output filename
    -quiet              suppress error messages
    -v, -version        display version
```


## EXAMPLES



EXAMPLE

Template (named "examples/pongo/weather.tmpl")

    Date: {{ now }}

    Hello {{ name }},
    
    The current weather is

    {{ weatherForecast.data.weather.0 }}

    Thank you

    {{ signature }}

Render the template above (i.e. examples/pongo/weather.tmpl) 
would be accomplished from the following data sources--

 + "now" and "name" are strings
 + "weatherForecast.data.weather.0" is a data path inside the 
    JSON document retrieved from a URL, ".0" references the
    "0"-th element of the array
 + "signature" comes from a file in our local disc 
   (i.e. examples/signature.txt)

That would be expressed on the command line as follows

    mkpongo "now=text:$(date)" "name=text:Little Frieda" \
        "weatherForecast=http://forecast.weather.gov/MapClick.php?lat=13.47190933300044&lon=144.74977715100056&FcstType=json" \
        signature=examples/signature.txt \
        examples/pongo/weather.tmpl     

Pongo2 is a Django/Jinga2 inspired template langauge
implemeted in Go.

      https://github.com/flosch/pongo2


