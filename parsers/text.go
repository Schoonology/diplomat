package parsers

import (
	"strconv"
	"strings"

	"github.com/testdouble/http-assertion-tool/http"
	"github.com/testdouble/http-assertion-tool/loaders"
)

const (
	AwaitingRequest = iota
	InRequestHeaders
	AwaitingRequestBody
	InResponseHeaders
	AwaitingResponseBody
)

type PlainTextParser struct{}

func (m *PlainTextParser) Parse(body *loaders.Body) (*Spec, error) {
	mode := AwaitingRequest
	tests := make([]Test, 0)
	currentTest := Test{}
	for _, line := range body.Lines {
		switch {
		case len(line) == 0:
			continue
		case line[0] == '>':
			line = strings.TrimSpace(line[1:])

			if mode == AwaitingRequest || mode == InResponseHeaders || mode == AwaitingResponseBody {
				if mode == InResponseHeaders || mode == AwaitingResponseBody {
					tests = append(tests, currentTest)
					currentTest = Test{}
				}

				currentTest.Request = RequestFromLine(line)
				mode = InRequestHeaders
			} else if mode == InRequestHeaders && len(line) == 0 {
				mode = AwaitingRequestBody
			} else if mode == InRequestHeaders {
				key, value := HeaderFromLine(line)
				currentTest.Request.Headers[key] = value
			}
		case line[0] == '<':
			line = strings.TrimSpace(line[1:])

			if mode == InRequestHeaders || mode == AwaitingRequestBody {
				currentTest.Response = ResponseFromLine(line)
				mode = InResponseHeaders
			} else if mode == InResponseHeaders && len(line) == 0 {
				mode = AwaitingResponseBody
			} else if mode == InResponseHeaders {
				key, value := HeaderFromLine(line)
				currentTest.Response.Headers[key] = value
			}
		}
	}

	if mode != AwaitingRequest {
		tests = append(tests, currentTest)
	}

	spec := Spec{
		Tests: tests,
	}

	return &spec, nil
}

func RequestFromLine(line string) *http.Request {
	pieces := strings.Split(line, " ")

	return http.NewRequest(pieces[0], pieces[1])
}

func HeaderFromLine(line string) (string, string) {
	pieces := strings.Split(line, ":")

	return strings.TrimSpace(pieces[0]), strings.TrimSpace(pieces[1])
}

func ResponseFromLine(line string) *http.Response {
	pieces := strings.Split(line, " ")
	code, _ := strconv.Atoi(pieces[1])

	return http.NewResponse(code, strings.Join(pieces[2:], " "))
}
