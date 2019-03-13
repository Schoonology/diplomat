package main

import (
	"os"

	"github.com/testdouble/diplomat/differs"
	"github.com/testdouble/diplomat/errors"
	"github.com/testdouble/diplomat/http"
	"github.com/testdouble/diplomat/loaders"
	"github.com/testdouble/diplomat/parsers"
	"github.com/testdouble/diplomat/printers"
	"github.com/testdouble/diplomat/runners"
	"github.com/testdouble/diplomat/scripting"
	"github.com/testdouble/diplomat/transforms"
	kingpin "gopkg.in/alecthomas/kingpin.v3-unstable"
)

var (
	app = kingpin.New("diplomat", "")

	debug   = app.Flag("debug", "Enable debug mode.").Bool()
	scripts = app.Flag("script", "Custom Lua script(s) to import.").PlaceHolder("FILE").Strings()
	tap     = app.Flag("tap", "Display results in TAP format.").Bool()

	filename = app.Arg("filename", "Treaty file to load.").Required().ExistingFile()
	address  = app.Arg("address", "Default base URL to use.").Required().String()
)

func init() {
	app.Version("0.0.1")
	app.UsageTemplate(`Usage: {{.App.Name}} {{.App.FlagSummary}} {{.App.ArgSummary}}

Flags:
{{.Context.Flags|FlagsToTwoColumns|FormatTwoColumns}}
Args:
{{.Context.Args|ArgsToTwoColumns|FormatTwoColumns}}`)
}

// Engine encapsulates all the behaviour of the tool as defined by the
// attached components.
type Engine struct {
	Loader     loaders.Loader
	Parser     parsers.SpecParser
	Transforms []transforms.Transformer
	Runner     runners.SpecRunner
	Printer    printers.ResultsPrinter
}

// Start runs the Engine.
func (r *Engine) Start(filename string, errors chan error) {
	lineChannel := r.Loader.Load(filename, errors)
	testChannel := r.Parser.Parse(lineChannel, errors)

	for _, transformer := range r.Transforms {
		testChannel = transformer.TransformAll(testChannel, errors)
	}

	resultChannel := r.Runner.Run(testChannel, errors)

	r.Printer.Print(resultChannel, errors)
}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	var printer printers.ResultsPrinter
	printer = &printers.Debug{}
	if *tap {
		printer = &printers.Tap{}
	}

	var differ differs.Differ
	differ = &differs.Smart{}
	if *debug {
		differ = &differs.Debug{}
	}

	for _, filename := range *scripts {
		err := scripting.LoadFile(filename)
		if err != nil {
			errors.Display(err)
			os.Exit(3)
		}
	}

	engine := Engine{
		Loader: &loaders.FileLoader{},
		Parser: parsers.GetParserFromFileName(*filename),
		Transforms: []transforms.Transformer{
			&transforms.TemplateRenderer{},
		},
		Runner: &runners.Serial{
			Client: http.NewClient(*address),
			Differ: differ,
		},
		Printer: printer,
	}

	errorStream := make(chan error)

	engine.Start(*filename, errorStream)

	for err := range errorStream {
		errors.Display(err)
		os.Exit(3)
	}
}
