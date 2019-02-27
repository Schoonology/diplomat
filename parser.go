package main

type Spec struct{}

type SpecParser struct{}

func (s *SpecParser) Parse(bytes []byte) (spec Spec, err error) {
	return
}
