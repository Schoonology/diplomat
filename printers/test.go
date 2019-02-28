package printers

import (
	"fmt"

	"github.com/testdouble/http-assertion-tool/runners"
)

type Test struct{}

func (t *Test) Print(result *runners.Result) error {
	fmt.Print(result)
	return nil
}
