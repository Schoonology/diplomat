package main

type Spec struct{}

type SpecParser struct{}

func (s *SpecParser) Parse(filename string) (spec Spec, err error) {
	return
}
