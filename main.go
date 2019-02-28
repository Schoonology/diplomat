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
)

type Args struct {
	Address  string
	Filename string
}

func loadArgs() (args Args) {
	flag.Parse()

	args.Filename = flag.Arg(0)
	args.Address = flag.Arg(1)

	return
}

type Engine struct {
	Loader  loaders.Loader
	Parser  parsers.SpecParser
	Runner  runners.SpecRunner
	Printer printers.ResultsPrinter
}

func (r *Engine) Start(filename string) error {
	file, err := r.Loader.Load(filename)
	if err != nil {
		return err
	}

	spec, err := r.Parser.Parse(file)
	if err != nil {
		return err
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
	engine := Engine{
		Loader: &loaders.FileLoader{},
		Parser: &parsers.PlainTextParser{},
		Runner: &runners.Serial{
			Client: &http.NativeClient{
				Address: args.Address,
			},
			Differ: &differs.Debug{},
		},
		Printer: &printers.Debug{},
	}

	err := engine.Start(args.Filename)
	if err != nil {
		fmt.Printf("Failed with: %v", err)
	}
}
