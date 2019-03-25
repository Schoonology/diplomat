package main

import (
	"os"

	"github.com/testdouble/diplomat/builders"
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
	address = app.Flag("address", "Default base URL to use.").Required().String()

	filename = app.Arg("filename", "Treaty file to load.").Required().ExistingFile()
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
	Parser     parsers.ParagraphParser
	Builder    builders.SpecBuilder
	Transforms []transforms.Transformer
	Runner     runners.SpecRunner
	Printer    printers.ResultsPrinter
}

// Start runs the Engine.
func (r *Engine) Start(filename string, errorChannel chan error) {
	lineChannel := r.Loader.Load(filename, errorChannel)
	paragraphChannel := r.Parser.Parse(lineChannel)
	testChannel := r.Builder.BuildAll(paragraphChannel)

	for _, transformer := range r.Transforms {
		testChannel = transformer.TransformAll(testChannel)
	}

	resultChannel := r.Runner.RunAll(testChannel)

	r.Printer.Print(resultChannel, errorChannel)
}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	var printer printers.ResultsPrinter
	printer = &printers.Pretty{}
	if *tap {
		printer = &printers.Tap{}
	} else if *debug {
		printer = &printers.Debug{}
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
			os.Exit(1)
		}
	}

	engine := Engine{
		Loader:  &loaders.FileLoader{},
		Parser:  parsers.GetParserFromFileName(*filename),
		Builder: &builders.State{},
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

	errorCount := 0
	for range errorStream {
		errorCount++
	}

	os.Exit(errorCount)
}
