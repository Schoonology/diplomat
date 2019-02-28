package differs

import "github.com/testdouble/http-assertion-tool/http"

type Differ interface {
	Diff(*http.Response, *http.Response) (string, error)
}
