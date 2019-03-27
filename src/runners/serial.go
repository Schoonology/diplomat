package runners

import (
	"github.com/testdouble/diplomat/builders"
	"github.com/testdouble/diplomat/differs"
	"github.com/testdouble/diplomat/http"
)

// Serial runs all tests one at a time.
type Serial struct {
	Differ differs.Differ
	Client http.Client
}

// Run returns the results of running all tests in `test`.
func (s *Serial) Run(test builders.Test) (TestResult, error) {
	response, err := s.Client.Do(test.Request)
	if err != nil {
		return TestResult{Name: test.Name}, err
	}

	diff, err := s.Differ.Diff(test.Response, response)
	if err != nil {
		return TestResult{Name: test.Name}, err
	}

	return TestResult{
		Name: test.Name,
		Diff: diff,
	}, nil
}

// RunAll returns the results of running all tests in the provided channel.
func (s *Serial) RunAll(tests chan builders.Test) chan TestResult {
	results := make(chan TestResult)

	go func() {
		defer close(results)

		for test := range tests {
			if test.Err != nil {
				results <- TestResult{
					Name: test.Name,
					Err:  test.Err,
				}
				continue
			}

			result, err := s.Run(test)
			if err != nil {
				result.Err = err
				results <- result
				continue
			}

			results <- result
		}
	}()

	return results
}
