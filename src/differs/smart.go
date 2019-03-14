package differs

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/testdouble/diplomat/scripting"

	"github.com/testdouble/diplomat/http"
)

var matcherRegex *regexp.Regexp

func init() {
	matcherRegex = regexp.MustCompile("^\\s*{\\? ([^\\?]+) \\?}\\s*$")
}

// The Smart differ provides looser restrictions on diffing, only printing
// diff output if an expected value is provided.
type Smart struct{}

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
		} else if actual.Headers[key] != value {
			output.WriteString(fmt.Sprintf("Invalid Header: %s\n", key))
			output.WriteString(fmt.Sprintf("	- %s\n", value))
			output.WriteString(fmt.Sprintf("	+ %s\n", actualValue))
		}
	}

	validatorMatch := matcherRegex.FindSubmatch(expected.Body)
	if len(validatorMatch) > 0 {
		validator := string(validatorMatch[1])
		valid, err := scripting.RunValidator(validator, string(actual.Body))
		if err != nil {
			return "", err
		}

		if !valid {
			output.WriteString(fmt.Sprintf("Body did not match validator: %s", validator))
		}
	} else if len(expected.Body) > 0 {
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
