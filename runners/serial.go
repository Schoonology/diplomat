package runners

import (
	"github.com/testdouble/http-assertion-tool/http"
	"github.com/testdouble/http-assertion-tool/parsers"
)

type Serial struct {
	Differ http.Differ
	Client http.Client
}

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

		result.Results = append(result.Results, diff)
	}

	return result, nil
}
