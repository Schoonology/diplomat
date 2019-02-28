package runners

import "github.com/testdouble/http-assertion-tool/parsers"

type SpecRunner interface {
	Run(*parsers.Spec) (*Result, error)
}

type TestResult struct {
	Name string
	Diff string
}

type Result struct {
	Results []TestResult
}
