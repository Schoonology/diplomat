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
func (s *Serial) Run(tests chan parsers.Test, errors chan error) chan TestResult {
	results := make(chan TestResult)

	go func() {
		for test := range tests {
			response, err := s.Client.Do(test.Request)
			if err != nil {
				errors <- err
				return
			}

			diff, err := s.Differ.Diff(test.Response, response)
			if err != nil {
				errors <- err
				return
			}

			results <- TestResult{
				Name: test.Name,
				Diff: diff,
			}
		}

		close(results)
	}()

	return results
}
