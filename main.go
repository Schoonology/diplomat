package main

import (
	"flag"
	"fmt"

	"github.com/testdouble/http-assertion-tool/differs"
	"github.com/testdouble/http-assertion-tool/http"
	"github.com/testdouble/http-assertion-tool/loaders"
	"github.com/testdouble/http-assertion-tool/parsers"
	"github.com/testdouble/http-assertion-tool/printers"
	"github.com/testdouble/http-assertion-tool/runners"
	"github.com/testdouble/http-assertion-tool/transforms"
)

// Args contains all CLI arguments passed to the tool.
type Args struct {
	Tap      bool
	Address  string
	Filename string
}

func loadArgs() (args Args) {
	flag.BoolVar(&args.Tap, "tap", false, "Display results in TAP format")

	flag.Parse()

	args.Filename = flag.Arg(0)
	args.Address = flag.Arg(1)

	return
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
	args := loadArgs()

	var printer printers.ResultsPrinter
	printer = &printers.Debug{}
	if args.Tap {
		printer = &printers.Tap{}
	}

	engine := Engine{
		Loader: &loaders.FileLoader{},
		Parser: parsers.GetParserFromFileName(args.Filename),
		Transforms: []transforms.Transform{
			transforms.RenderTemplates,
		},
		Runner: &runners.Serial{
			Client: http.NewClient(args.Address),
			Differ: &differs.Debug{},
		},
		Printer: printer,
	}

	err := engine.Start(args.Filename)
	if err != nil {
		fmt.Printf("Failed with: %v", err)
	}
}
