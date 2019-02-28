package printers

import (
	"fmt"

	"github.com/testdouble/http-assertion-tool/runners"
)

type Debug struct{}

func (t *Debug) Print(result *runners.Result) error {
	for _, diff := range result.Results {
		fmt.Print(diff)
	}

	return nil
}
