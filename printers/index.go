package printers

import "github.com/testdouble/http-assertion-tool/runners"

type ResultsPrinter interface {
	Print(*runners.Result) error
}
