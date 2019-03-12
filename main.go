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
	Transforms []transforms.Transform
	Runner     runners.SpecRunner
	Printer    printers.ResultsPrinter
}

// Start runs the Engine.
func (r *Engine) Start(filename string) error {
	file, err := r.Loader.Load(filename)
	if err != nil {
		return err
	}

	spec, err := r.Parser.Parse(file)
	if err != nil {
		return err
	}

	for _, transform := range r.Transforms {
		err = transform(spec)
		if err != nil {
			return err
		}
	}

	result, err := r.Runner.Run(spec)
	if err != nil {
		return err
	}

	err = r.Printer.Print(result)
	if err != nil {
		return err
	}

	return nil
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
			panic(err)
		}
	}

	engine := Engine{
		Loader: &loaders.FileLoader{},
		Parser: parsers.GetParserFromFileName(*filename),
		Transforms: []transforms.Transform{
			transforms.RenderTemplates,
		},
		Runner: &runners.Serial{
			Client: http.NewClient(*address),
			Differ: differ,
		},
		Printer: printer,
	}

	err := engine.Start(*filename)
	if err != nil {
		errors.Display(err)
		os.Exit(3)
	}
}
