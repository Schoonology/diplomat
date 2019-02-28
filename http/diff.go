package http

import "github.com/google/go-cmp/cmp"

// TODO(schoon) - Is there a better way to allow mocking for functional
// components? Should dependents take a `func` member instead?
type DebugDiffer struct{}

func (s *DebugDiffer) Diff(expected *Response, actual *Response) (string, error) {
	return cmp.Diff(expected, actual), nil
}
