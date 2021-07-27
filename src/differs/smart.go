package differs

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/schoonology/diplomat/scripting"

	"github.com/schoonology/diplomat/http"
)

var isValidatorRegex *regexp.Regexp
var expressionRegex *regexp.Regexp

func init() {
	// The syntax {? func ?} is used to embed Diplomat validators.
	isValidatorRegex = regexp.MustCompile("^([\\s]*{\\?[^\\?]+\\?}[\\s]*)+$")
	expressionRegex = regexp.MustCompile("{\\?[\\s]*([^\\?]+?)[\\s]*\\?}")
}

// The Smart differ provides looser restrictions on diffing, only printing
// diff output if an expected value is provided.
type Smart struct{}

func forAllValidators(source string, iterator func(expression string) (bool, error)) ([]bool, error) {
	if !isValidatorRegex.MatchString(source) {
		return []bool{}, nil
	}

	matches := expressionRegex.FindAllStringSubmatch(source, -1)
	results := make([]bool, len(matches))
	for idx, match := range matches {
		expression := string(match[1])
		valid, err := iterator(expression)
		if err != nil {
			return nil, err
		}

		results[idx] = valid
	}

	return results, nil
}

// Diff returns the difference between `expected` and `actual`.
func (s *Smart) Diff(expected *http.Response, actual *http.Response) (string, error) {
	output := strings.Builder{}

	if actual.StatusCode != expected.StatusCode || actual.StatusText != expected.StatusText {
		output.WriteString("Status:\n")
		output.WriteString(fmt.Sprintf("	- %d %s\n", expected.StatusCode, expected.StatusText))
		output.WriteString(fmt.Sprintf("	+ %d %s\n", actual.StatusCode, actual.StatusText))
	}

	for key, value := range expected.Headers {
		actualValue, present := actual.Headers[key]

		if !present {
			output.WriteString(fmt.Sprintf("Missing Header: %s\n", key))
			continue
		}

		results, err := forAllValidators(value, func(expr string) (bool, error) {
			valid, err := scripting.RunValidator(expr, actualValue)
			if err == nil && !valid {
				output.WriteString(fmt.Sprintf("Header `%s` did not match validator: %s\n", key, expr))
			}
			return valid, err
		})
		if err != nil {
			return "", err
		}

		if len(results) == 0 && actual.Headers[key] != value {
			output.WriteString(fmt.Sprintf("Invalid Header: %s\n", key))
			output.WriteString(fmt.Sprintf("	- %s\n", value))
			output.WriteString(fmt.Sprintf("	+ %s\n", actualValue))
		}
	}

	results, err := forAllValidators(string(expected.Body), func(expr string) (bool, error) {
		valid, err := scripting.RunValidator(expr, string(actual.Body))
		if err == nil && !valid {
			output.WriteString(fmt.Sprintf("Body did not match validator: %s", expr))
		}
		return valid, err
	})
	if err != nil {
		return "", err
	}

	if len(results) == 0 && len(expected.Body) > 0 {
		contentType, present := expected.Headers["Content-Type"]
		if !present {
			contentType, present = actual.Headers["Content-Type"]
		}

		if !present {
			output.WriteString("Missing Content-Type with Body assertion.")
		} else {
			bodyDiff, err := diffBody(expected.Body, actual.Body, contentType)
			if err != nil {
				return "", err
			}

			if len(bodyDiff) > 0 {
				output.WriteString("Invalid Body:\n")
				output.WriteString(bodyDiff)
			}
		}
	}

	return output.String(), nil
}
