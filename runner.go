package main

type Runner struct {
	parser SpecParser
}

func NewRunner(parser SpecParser) Runner {
	return Runner{
		parser: parser,
	}
}

func (r *Runner) Run(filename string) error {
	_, err := r.parser.Parse(filename)
	if err != nil {
		return err
	}

	return nil
}
