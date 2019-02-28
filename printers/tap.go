package printers

import (
	"fmt"

	"github.com/testdouble/http-assertion-tool/runners"
)

type Tap struct{}

func (t *Tap) Print(result *runners.Result) error {
	fmt.Println("TAP version 13")
	fmt.Printf("1..%d\n", len(result.Results))

	for idx, result := range result.Results {
		failed := len(result.Diff) > 0
		status := "ok"
		if failed {
			status = "not ok"
		}

		fmt.Printf("%s %d %s\n", status, idx, result.Name)

		if failed {
			fmt.Print(result.Diff)
		}
	}

	return nil
}
