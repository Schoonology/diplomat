package main

import (
	"fmt"
	"os"

	"github.com/schoonology/diplomat/builders"
	"github.com/schoonology/diplomat/colors"
	"github.com/schoonology/diplomat/differs"
	"github.com/schoonology/diplomat/errors"
	"github.com/schoonology/diplomat/formatters"
	"github.com/schoonology/diplomat/http"
	"github.com/schoonology/diplomat/loaders"
	"github.com/schoonology/diplomat/loggers"
	"github.com/schoonology/diplomat/parsers"
	"github.com/schoonology/diplomat/runners"
	"github.com/schoonology/diplomat/scripting"
	"github.com/schoonology/diplomat/transforms"
	kingpin "gopkg.in/alecthomas/kingpin.v3-unstable"
)

var (
	app = kingpin.New("diplomat", "")

	color   = app.Flag("color", "Enable colored output.").Bool()
	debug   = app.Flag("debug", "Enable debug mode.").Bool()
	scripts = app.Flag("script", "Custom Lua script(s) to import.").PlaceHolder("FILE").Strings()
	tap     = app.Flag("tap", "Display results in TAP format.").Bool()
	address = app.Flag("address", "Default base URL to use.").Required().String()

	filenames = app.Arg("filenames", "Treaty file to load.").Required().ExistingFiles()
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
	Parser     parsers.ParseDelegator
	Builder    builders.SpecBuilder
	Transforms []transforms.Transformer
	Runner     runners.SpecRunner
	Formatter  formatters.ResultsFormatter
	Logger     loggers.Logger
}

// Start runs the Engine.
func (r *Engine) Start(filenameChannel chan string, errorChannel chan error) {
	fileChannel := r.Loader.LoadAll(filenameChannel, errorChannel)
	paragraphChannel := r.Parser.ParseAll(fileChannel)
	testChannel := r.Builder.BuildAll(paragraphChannel)

	for _, transformer := range r.Transforms {
		testChannel = transformer.TransformAll(testChannel)
	}

	resultChannel := r.Runner.RunAll(testChannel)

	outputChannel := r.Formatter.FormatAll(resultChannel, errorChannel)
	r.Logger.PrintAll(outputChannel, errorChannel)
}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	colorizer := colors.DefaultColorizer(*color)

	var formatter formatters.ResultsFormatter
	formatter = &formatters.Pretty{
		Colorizer: colorizer,
	}
	if *tap {
		formatter = &formatters.Tap{}
	} else if *debug {
		formatter = &formatters.Debug{}
	}

	var differ differs.Differ
	differ = &differs.Smart{}
	if *debug {
		differ = &differs.Debug{}
	}

	for _, filename := range *scripts {
		err := scripting.LoadFile(filename)
		if err != nil {
			fmt.Print(errors.Format(err))
			os.Exit(1)
		}
	}

	engine := Engine{
		Loader:     &loaders.FileLoader{},
		Parser:     &parsers.Delegator{},
		Builder:    &builders.State{},
		Transforms: []transforms.Transformer{},
		Runner: &runners.Serial{
			Client: http.NewClient(*address),
			Differ: differ,
			Transforms: []transforms.Transformer{
				&transforms.TemplateRenderer{},
			},
		},
		Formatter: formatter,
		Logger:    &loggers.StandardOutput{},
	}

	filenameStream := make(chan string)
	errorStream := make(chan error)

	// send all the filenames to the channel
	go func() {
		for _, f := range *filenames {
			filenameStream <- f
		}
		close(filenameStream)
	}()

	engine.Start(filenameStream, errorStream)

	errorCount := 0
	for range errorStream {
		errorCount++
	}

	os.Exit(errorCount)
}
