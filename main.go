package main

import (
	"flag"
	"fmt"

	"github.com/testdouble/diplomat/differs"
	"github.com/testdouble/diplomat/http"
	"github.com/testdouble/diplomat/loaders"
	"github.com/testdouble/diplomat/parsers"
	"github.com/testdouble/diplomat/printers"
	"github.com/testdouble/diplomat/runners"
	"github.com/testdouble/diplomat/scripting"
	"github.com/testdouble/diplomat/transforms"
)

type customScripts []string

func (i *customScripts) String() string {
	return "Scripts"
}

func (i *customScripts) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// Args contains all CLI arguments passed to the tool.
type Args struct {
	Tap      bool
	Debug    bool
	Address  string
	Filename string
	Scripts  customScripts
}

func loadArgs() (args Args) {
	flag.BoolVar(&args.Tap, "tap", false, "Display results in TAP format")
	flag.BoolVar(&args.Debug, "debug", false, "Display results using the debug differ")
	flag.Var(&args.Scripts, "script", "Custom scripts")

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

	var differ differs.Differ
	differ = &differs.Smart{}
	if args.Debug {
		differ = &differs.Debug{}
	}

	for _, filename := range args.Scripts {
		err := scripting.LoadFile(filename)
		if err != nil {
			panic(err)
		}
	}

	engine := Engine{
		Loader: &loaders.FileLoader{},
		Parser: parsers.GetParserFromFileName(args.Filename),
		Transforms: []transforms.Transform{
			transforms.RenderTemplates,
		},
		Runner: &runners.Serial{
			Client: http.NewClient(args.Address),
			Differ: differ,
		},
		Printer: printer,
	}

	err := engine.Start(args.Filename)
	if err != nil {
		fmt.Printf("Failed with: %v", err)
	}
}
