package runners

import "github.com/testdouble/http-assertion-tool/parsers"

type SpecRunner interface {
	Run(*parsers.Spec) (*Result, error)
}

type Result struct {
	Results []string
}
