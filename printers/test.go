package printers

import "github.com/testdouble/http-assertion-tool/runners"

type Test struct{}

func (t *Test) Print(*runners.Result) error {
	return nil
}
