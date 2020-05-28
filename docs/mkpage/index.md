
# mkpage

## USAGE

	mkpage [OPTIONS] [KEY/VALUE DATA PAIRS] [TEMPLATE_FILENAMES]

## DESCRIPTION

A Pandoc preprossor supporting front matter in YAML, TOML and JSON.
The command parameter language use the key/value pairs populate 
template(s) and render to stdout. Values can be explicitly typed,
derived from files and their extensions or retrieved from the net.

## OPTIONS

Below are a set of options available. Options will override any corresponding environment settings.

```
    -code               outout just code blocks for language, 
                        e.g. shell or json, reads from standard input
    -codesnip           output just the code bocks, reads from
                        standard input
    -examples           display example(s)
    -generate-manpage   generate man page
    -generate-markdown  generate markdown documentation
    -gt, -go-templates  (DEPRECIATED) use Go's template engine instead of 
                        Pandoc's template engine
    -h, -help           display help
    -i, -input          input filename
    -l, -license        display license
    -o, -output         output filename
    -quiet              suppress error messages
    -s, -show-template  display source for a default page template
    -t, -templates      (DEPRECIATED) colon delimited list of 
                        templates to use
    -v, -version        display version
```


## EXAMPLES



EXAMPLE

Template (named "examples/weather.tmpl")
    
```
    Date: ${now}
    
    Hello ${name},
        
    The weather forcast is
    
    $if(weather.data.weather)$
      $weather.data.weather[; ]$
    $endif$
    
    Thank you
    
    ${signature}

```

Render the template above (i.e. examples/weather.tmpl) would be 
accomplished from the following data sources--

 + "now" and "name" are strings
 + "weather" is JSON data retrieved from a URL
 	+ ".data.weather" is a data path inside the JSON document
	+ "index" let's us pull our the "0"-th element (i.e. the initial element of the array)
 + "signature" comes from a file in our local disc (i.e. examples/signature.txt)

That would be expressed on the command line as follows

```shell
    mkpage "now=text:$(date)" "name=text:Little Frieda" \
        "weatherForecast=http://forecast.weather.gov/MapClick.php?lat=13.47190933300044&lon=144.74977715100056&FcstType=json" \
        signature=examples/signature.txt \
        examples/weather.tmpl     
```



