
# mkpage

## USAGE

	mkpage [OPTIONS] [KEY/VALUE DATA PAIRS] [TEMPLATE_FILENAMES]

## DESCRIPTION

Using the key/value pairs populate the template(s) and render to stdout.


## ENVIRONMENT

Environment variables can be overridden by corresponding options

```
    MKPAGE_TEMPLATES   # set the default template path
```

## OPTIONS

Below are a set of options available. Options will override any 
corresponding environment settings.

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
    -pandoc             use Pandoc's template engine
    -quiet              suppress error messages
    -s                  display the default template
    -show-template      display the default template
    -t                  colon delimited list of templates to use
    -templates          colon delimited list of templates to use
    -v, -version        display version
```


## EXAMPLES

Template (named "examples/weather.tmpl")
    
```
    Date: ${now}
    
    Hello ${name},
        
    The weather forcast is
    
    ${if(weather.data.weather)}
      ${weather.data.weather[; ]}
    ${endif}
    
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

```
    mkpage "now=text:$(date)" "name=text:Little Frieda" \
        "weather=http://forecast.weather.gov/MapClick.php?lat=13.47190933300044&lon=144.74977715100056&FcstType=json" \
        signature=examples/signature.txt \
        examples/weather.tmpl     
```

mkpage v0.0.33
