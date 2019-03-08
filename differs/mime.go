package differs

import (
	"encoding/json"
	"fmt"
	"mime"
	"strings"

	"github.com/google/go-cmp/cmp"
)

func diffBody(expected []byte, actual []byte, contentType string) (string, error) {
	mediaType, _, _ := mime.ParseMediaType(contentType)

	switch mediaType {
	case "application/json":
		return diffJSON(expected, actual)
	default:
		return diffText(expected, actual)
		// TODO(schoon) - Once we have a "strict" mode:
		// return "", errors.New(fmt.Sprintf("Unsupported media type: %s", mediaType))
	}
}

func diffText(expected []byte, actual []byte) (string, error) {
	expectedBody := strings.TrimSpace(string(expected))
	actualBody := strings.TrimSpace(string(actual))

	if expectedBody != actualBody {
		output := strings.Builder{}

		output.WriteString(fmt.Sprintf("	-: %v\n", actualBody))
		output.WriteString(fmt.Sprintf("	+: %v\n", expectedBody))

		return output.String(), nil
	}

	return "", nil
}

func diffJSON(expected []byte, actual []byte) (string, error) {
	parsedExp := make(map[string]interface{})
	parsedAct := make(map[string]interface{})

	err := json.Unmarshal(expected, &parsedExp)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(actual, &parsedAct)
	if err != nil {
		return "", err
	}

	return cmp.Diff(parsedExp, parsedAct), nil
}
