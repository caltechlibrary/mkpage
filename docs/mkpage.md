---
{
    "has_code": true
}
---

# mkpage

## USAGE

    mkpage [OPTION] [KEY/VALUE DATA PAIRS] TEMPLATE_FILENAME [TEMPLATE_FILENAMES]

## SYNOPSIS

Using the key/value pairs populate the template(s) and render to stdout.

## CONFIGURATION

You can set a local default template path by using environment variables.

+ MKPAGE_TEMPLATES - is the colon delimited list of template paths

## OPTIONS

```
	-h	show help
	-help	show help
	-l	show license
	-license	show license
	-s	display the default template
	-show-template	display the default template
	-t	colon delimited list of templates to use
	-templates	colon delimited list of templates to use
	-v	show version
	-version	show version
```

## EXAMPLE

Template (named "examples/weather.tmpl")

```
    {{ define "weather.tmpl" }}
    Date: {{- .now}}

    Hello {{.name -}},
    
    The current weather is

    {{index .weatherForecast.data.weather 0}}

    Thank you

    {{.signature}}
    {{ end }}
```

Render the template above (i.e. examples/weather.tmpl) would be accomplished from 
the following data sources--

+ "now" and "name" are strings
+ "weatherForecast" is JSON data retrieved from a URL
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

Golang's text/template docs can be found at 

```
      https://golang.org/pkg/text/template/
```

