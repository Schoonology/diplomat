package differs

import (
	"fmt"
	"strings"

	"github.com/testdouble/http-assertion-tool/http"
)

type Smart struct{}

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

	return output.String(), nil
}
