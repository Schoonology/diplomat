package differs

import "github.com/testdouble/http-assertion-tool/http"

// A Differ returns a string representing the difference between an expected
// and actual Response.
type Differ interface {
	Diff(*http.Response, *http.Response) (string, error)
}
