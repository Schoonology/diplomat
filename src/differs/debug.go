package differs

import (
	"github.com/google/go-cmp/cmp"
	"github.com/schoonology/diplomat/http"
)

// The Debug differ returns any and all differences.
type Debug struct{}

// Diff returns the difference between `expected` and `actual`.
func (s *Debug) Diff(expected *http.Response, actual *http.Response) (string, error) {
	return cmp.Diff(expected, actual), nil
}
