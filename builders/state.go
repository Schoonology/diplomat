package builders

import (
	"strconv"
	"strings"

	"github.com/testdouble/diplomat/errors"
	"github.com/testdouble/diplomat/http"
)

// The State builder converts a set of lines into a test.
type State struct{}

const (
	modeAwaitingRequest = iota
	modeInRequestHeaders
	modeInRequestBody
	modeInResponseHeaders
	modeInResponseBody
)

func isRequestMetadataIndicator(char byte) bool {
	return char == '>'
}

func isResponseMetadataIndicator(char byte) bool {
	return char == '<'
}

func isRequestMode(mode int) bool {
	return mode == modeInRequestHeaders || mode == modeInRequestBody
}

func isResponseMode(mode int) bool {
	return mode == modeInResponseHeaders || mode == modeInResponseBody
}

func requestFromLine(line string) (*http.Request, error) {
	pieces := strings.Split(line, " ")

	if len(pieces) < 2 {
		return nil, &errors.BadRequestLine{Line: line}
	}

	return http.NewRequest(pieces[0], pieces[1]), nil
}

func headerFromLine(line string) (string, string, error) {
	pieces := strings.Split(line, ":")

	if len(pieces) < 2 {
		return "", "", &errors.BadHeader{Header: line}
	}

	return strings.TrimSpace(pieces[0]), strings.TrimSpace(pieces[1]), nil
}

func responseFromLine(line string) (*http.Response, error) {
	pieces := strings.Split(line, " ")

	if len(pieces) < 3 {
		return nil, &errors.BadResponseStatus{Line: line}
	}

	code, _ := strconv.Atoi(pieces[1])

	return http.NewResponse(code, strings.Join(pieces[2:], " ")), nil
}

// Build parses a set of lines into a test.
func (s *State) Build(lines []string) (Test, error) {
	mode := modeAwaitingRequest
	test := Test{}

	for _, line := range lines {
		char := byte(0)
		trimmedLine := ""

		if len(line) > 0 {
			char = line[0]
			trimmedLine = strings.TrimSpace(line[1:])
		}

		if isRequestMetadataIndicator(char) {
			if isResponseMode(mode) {
				return test, &errors.MisplacedRequest{}
			}
			if mode == modeAwaitingRequest {
				mode = modeInRequestHeaders
				request, err := requestFromLine(trimmedLine)
				if err != nil {
					return test, err
				}

				test.Request = request
			} else if mode == modeInRequestHeaders && len(trimmedLine) == 0 {
				mode = modeInRequestBody
			} else if mode == modeInRequestHeaders {
				key, value, err := headerFromLine(trimmedLine)
				if err != nil {
					return test, err
				}

				test.Request.Headers[key] = value
			}
		} else if isResponseMetadataIndicator(char) {
			if isRequestMode(mode) {
				mode = modeInResponseHeaders
				response, err := responseFromLine(trimmedLine)
				if err != nil {
					return test, err
				}

				test.Response = response
			} else if mode == modeInResponseHeaders && len(trimmedLine) == 0 {
				mode = modeInResponseBody
			} else if mode == modeInResponseHeaders {
				key, value, err := headerFromLine(trimmedLine)
				if err != nil {
					return test, err
				}

				test.Response.Headers[key] = value
			} else if mode == modeAwaitingRequest || mode == modeInResponseBody {
				// This can happen if an expected response body contains a response
				// indicator (e.g. `<`) as the first character. This explicitly
				// disallows that, preventing unexpected issues with missing requests.
				return test, &errors.MissingRequest{}
			}
		} else if isRequestMode(mode) {
			mode = modeInRequestBody
			test.Request.Body = append(test.Request.Body, []byte(line+"\n")...)
		} else if isResponseMode(mode) {
			mode = modeInResponseBody
			test.Response.Body = append(test.Response.Body, []byte(line+"\n")...)
		}
	}

	return test, nil
}

// BuildAll returns the results of running Build on all string arrays in a channel.
// spec in this case refers to a group of lines
func (s *State) BuildAll(specs chan []string, errors chan error) chan Test {
	tests := make(chan Test)

	go func() {
		for spec := range specs {
			test, err := s.Build(spec)
			if err != nil {
				errors <- err
				return
			}

			tests <- test
		}

		close(tests)
	}()

	return tests
}
