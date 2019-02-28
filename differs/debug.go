package differs

import (
	"github.com/google/go-cmp/cmp"
	"github.com/testdouble/http-assertion-tool/http"
)

// TODO(schoon) - Is there a better way to allow mocking for functional
// components? Should dependents take a `func` member instead?
type Debug struct{}

func (s *Debug) Diff(expected *http.Response, actual *http.Response) (string, error) {
	return cmp.Diff(expected, actual), nil
}
