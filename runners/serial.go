package runners

import (
	"github.com/testdouble/diplomat/differs"
	"github.com/testdouble/diplomat/http"
	"github.com/testdouble/diplomat/parsers"
)

// Serial runs all tests one at a time.
type Serial struct {
	Differ differs.Differ
	Client http.Client
}

// Run returns the results of running all tests in `spec`.
func (s *Serial) Run(spec *parsers.Spec) (*Result, error) {
	result := new(Result)

	for _, test := range spec.Tests {
		response, err := s.Client.Do(test.Request)
		if err != nil {
			return nil, err
		}

		diff, err := s.Differ.Diff(test.Response, response)
		if err != nil {
			return nil, err
		}

		result.Results = append(result.Results, TestResult{
			Name: test.Name,
			Diff: diff,
		})
	}

	return result, nil
}
