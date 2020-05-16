//
// mkpandoc.go is a MkPage front end to pandoc. It is intended for experimenting
// with unifying features developed in the mkpage project with pandoc as the rendering engine.
//
// Concepts to map to pandoc
//
// + mkpage supports frontmatter in JSON, YAML, TOML and pandoc supports metadata in YAML and JSON
//   + we can split the frontmatter into a temporary metadata file in JSON format and pass to pandoc vi command line.
// + pandoc provides a simple and readable template language
// 	 + leverage front matter and mkpage key/value pairs passing it to pandoc
//   + leverage pandoc's template engine for rendering output
//
package main

import (
	"fmt"
	"log"
	"os/exec"

	// Caltech Library Packages
	"github.com/caltechlibrary/cli"
	"github.com/caltechlibrary/mkpage"
)

var (
	description = `
SYNOPSIS

Using the key/value pairs populate the template(s) and render to stdout.
`

	examples = `

EXAMPLE

Template (named "examples/pandoc/weather.tmpl")

    Date: $now$

    Hello $name$,

    The current weather is

    $weatherForecast.data.weather.0$

    Thank you

    $signature$

Render the template above (i.e. examples/pandoc/weather.tmpl)
would be accomplished from the following data sources--

 + "now" and "name" are strings
 + "weatherForecast.data.weather.0" is a data path inside the
    JSON document retrieved from a URL, ".0" references the
    "0"-th element of the array
 + "signature" comes from a file in our local disc
   (i.e. examples/signature.txt)

That would be expressed on the command line as follows

    %s "now=text:$(date)" "name=text:Little Frieda" \
        "weatherForecast=http://forecast.weather.gov/MapClick.php?lat=13.47190933300044&lon=144.74977715100056&FcstType=json" \
        signature=examples/signature.txt \
        examples/pongo/weather.tmpl

Pandoc has a simple template engine we mkpandoc uses.
See https://pandoc.org/MANUAL.html#templates

`

	// Standard Options
	showHelp         bool
	showVersion      bool
	showLicense      bool
	showExamples     bool
	inputFName       string
	outputFName      string
	quiet            bool
	generateMarkdown bool
	generateManPage  bool

	// Application Options
	templateFName string
	codesnip      bool
	codeType      string
)

func main() {
	pandoc, err := exec.LookPath("pandoc")
	if err != nil {
		log.Fatal("Can't find pandoc, see https://pandoc.org")
	}

	app := cli.NewCli(mkpage.Version)
	appName := app.AppName()

	// Document expected parameters
	app.SetParams(`[KEY/VALUE DATA PAIRS]`, `TEMPLATE_FILENAME`)

	// Add Help docs
	app.AddHelp("license", []byte(fmt.Sprintf(mkpage.LicenseText, appName, mkpage.Version)))
	app.AddHelp("description", []byte(fmt.Sprintf(description)))
	app.AddHelp("examples", []byte(fmt.Sprintf(examples, appName)))

	// Standard Options
	app.BoolVar(&showHelp, "h,help", false, "display help")
	app.BoolVar(&showVersion, "v,version", false, "display version")
	app.BoolVar(&showExamples, "examples", false, "display example(s)")
	app.BoolVar(&showLicense, "l,license", false, "display license")
	app.StringVar(&inputFName, "i,input", "", "input filename")
	app.StringVar(&outputFName, "o,output", "", "output filename")
	app.BoolVar(&quiet, "quiet", false, "suppress error messages")
	app.BoolVar(&generateMarkdown, "generate-markdown", false, "generate markdown documentation")
	app.BoolVar(&generateManPage, "generate-manpage", false, "generate man page")

	// Application specific options
	app.BoolVar(&codesnip, "codesnip", false, "output just the code bocks")
	app.StringVar(&codeType, "code", "", "outout just code blocks for language, e.g. shell or json")

	app.Parse()
	args := app.Args()

	if showHelp || showExamples {
		if len(args) > 0 {
			fmt.Fprintln(app.Out, app.Help(args...))
		} else {
			app.Usage(app.Out)
		}
		os.Exit(0)
	}
	if showLicense {
		fmt.Fprintln(app.Out, app.License())
		os.Exit(0)
	}
	if showVersion {
		fmt.Fprintln(app.Out, app.Version())
		os.Exit(0)
	}

	// Default template name is page.tmpl
	var (
		templateName string
		err          error
	)

	// Setup IO
	app.Eout = os.Stderr
	app.In, err = cli.Open(inputFName, os.Stdin)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(inputFName, app.In)

	app.Out, err = cli.Create(outputFName, os.Stdout)
	cli.ExitOnError(app.Eout, err, quiet)
	defer cli.CloseFile(outputFName, app.Out)

	// Process flags and update the environment as needed.
	if generateMarkdown {
		app.GenerateMarkdown(app.Out)
		os.Exit(0)
	}
	if generateManPage {
		app.GenerateManPage(app.Out)
		os.Exit(0)
	}

	if codesnip || codeType != "" {
		err = mkpage.Codesnip(app.In, app.Out, codeType)
		cli.ExitOnError(app.Eout, err, quiet)
		os.Exit(0)
	}

	data := map[string]string{}
	for i, arg := range args {
		if strings.Contains(arg, "=") == true {
			// Update data map
			pair := strings.SplitN(arg, "=", 2)
			if len(pair) != 2 {
				fmt.Fprintf(app.Eout, "Can't read pair (%d) %s\n", i+1, arg)
				os.Exit(1)
			}
			data[pair[0]] = pair[1]
		} else {
			templateName = arg[:]
		}
	}

	err = mkpage.MakePandoc(app.Out, templateName, data)
	cli.ExitOnError(app.Eout, err, quiet)
}
