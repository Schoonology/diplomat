package main

import (
	"flag"
	"fmt"

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
	loader  loaders.Loader
	parser  parsers.SpecParser
	runner  runners.SpecRunner
	printer printers.ResultsPrinter
}

func (r *Engine) Start(filename string) error {
	file, err := r.loader.Load(filename)
	if err != nil {
		return err
	}

	spec, err := r.parser.Parse(file)
	if err != nil {
		return err
	}

	result, err := r.runner.Run(spec)
	if err != nil {
		return err
	}

	err = r.printer.Print(result)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	args := loadArgs()
	engine := Engine{
		loader: &loaders.FileLoader{},
		parser: &parsers.PlainTextParser{},
		runner: &runners.Serial{
			Address: args.Address,
		},
		printer: &printers.Test{},
	}

	err := engine.Start(args.Filename)
	if err != nil {
		fmt.Printf("Failed with: %v", err)
	}
}
