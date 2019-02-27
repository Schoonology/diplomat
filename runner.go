package main

import (
	"github.com/testdouble/http-assertion-tool/loaders"
	"github.com/testdouble/http-assertion-tool/parsers"
	"github.com/testdouble/http-assertion-tool/printers"
	"github.com/testdouble/http-assertion-tool/runners"
)

type Runner struct {
	loader  loaders.FileLoader
	parser  parsers.SpecParser
	runner  runners.SpecRunner
	printer printers.ResultsPrinter
}

type Results struct{}

func (r *Runner) Run(filename string) error {
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
