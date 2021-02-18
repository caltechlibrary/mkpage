
USAGE
=====

	mkpage [OPTIONS] [KEY/VALUE DATA PAIRS] [TEMPLATE_FILENAMES]

DESCRIPTION
-----------


Using the key/value pairs populate the template(s) and render to stdout.
MkPage renders markdown using Pandoc (version >= v2.10). 


OPTIONS
-------

Below are a set of options available.

```
    -code                outout just code blocks for specific language, e.g. shell or json, reads from standard input
    -codesnip            output just the code bocks, reads from standard input
    -examples            display example(s)
    -f, -from            set the from value (e.g. markdown) used by pandoc
    -generate-markdown   generate markdown documentation
    -h, -help            display help
    -i, -input           input filename
    -l, -license         display license
    -o, -output          output filename
    -pandoc-version      display Pandoc version found
    -t, -to              set the from value (e.g. html) used by pandoc
    -v, -version         display version
```


EXAMPLES
--------


Template (named "examples/weather.tmpl")
    
    Date: $now$
    
    Hello $name$,
        
    The weather forcast is
    
    $if(weather.data.weather)$
      $weather.data.weather[; ]$
    $endif$
    
    Thank you
    
    $signature$

Render the template above (i.e. examples/weather.tmpl) would be 
accomplished from the following data sources--

 + "now" and "name" are strings
 + "weather" is JSON data retrieved from a URL
 	+ ".data.weather" is a data path inside the JSON document
	+ "index" let's us pull our the "0"-th element (i.e. the initial element of the array)
 + "signature" comes from a file in our local disc (i.e. examples/signature.txt)

That would be expressed on the command line as follows

    mkpage "now=text:$(date)" "name=text:Little Frieda" \
        "weatherForecast=http://forecast.weather.gov/MapClick.php?lat=13.47190933300044&lon=144.74977715100056&FcstType=json" \
        signature=examples/signature.txt \
        examples/weather.tmpl     



mkpage v0.2.4
