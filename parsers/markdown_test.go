package parsers_test

import (
	"testing"

	"github.com/testdouble/diplomat/http"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/diplomat/parsers"
)

func streamBody(body []string) chan string {
	lines := make(chan string)

	go func() {
		for _, line := range body {
			lines <- line
		}

		close(lines)
	}()

	return lines
}

func assertTest(t *testing.T, expected parsers.Test, tests chan parsers.Test, errors chan error) {
	select {
	case err := <-errors:
		assert.FailNow(t, "Error should not exist.", err)
	case test := <-tests:
		assert.Equal(t, expected, test)
	}
}

func fillRequest(method string, path string, headers map[string]string, body string) *http.Request {
	request := http.NewRequest(method, path)
	for key, value := range headers {
		request.Headers[key] = value
	}
	request.Body = []byte(body)
	return request
}

func fillResponse(code int, status string, headers map[string]string, body string) *http.Response {
	response := http.NewResponse(code, status)
	for key, value := range headers {
		response.Headers[key] = value
	}
	response.Body = []byte(body)
	return response
}

func TestMarkdownText(t *testing.T) {
	subject := parsers.Markdown{}
	body := streamBody([]string{
		"# Markdown",
		"```",
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"```",
	})
	errors := make(chan error)

	tests := subject.Parse(body, errors)

	assertTest(t, parsers.Test{
		Name: "Markdown",
		Request: fillRequest("METHOD", "path", map[string]string{
			"Header": "Request",
		}, ""),
		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
			"Header": "Response",
		}, ""),
	}, tests, errors)
	// TODO(schoon) - Assert that the channel is closed here.
}

func TestMarkdownSplitReqRes(t *testing.T) {
	subject := parsers.Markdown{}
	body := streamBody([]string{
		"# Markdown",
		"```",
		"> METHOD path",
		"> Header: Request",
		"```",
		"```",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"```",
	})
	errors := make(chan error)

	tests := subject.Parse(body, errors)

	assertTest(t, parsers.Test{
		Name: "Markdown",
		Request: fillRequest("METHOD", "path", map[string]string{
			"Header": "Request",
		}, ""),
		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
			"Header": "Response",
		}, ""),
	}, tests, errors)
}

func TestMarkdownDouble(t *testing.T) {
	subject := parsers.Markdown{}
	body := streamBody([]string{
		"# Markdown",
		"## First request",
		"```",
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"```",
		"## Second request",
		"```",
		"> SECOND path",
		"> Header: Request 2",
		"< PROTO 1234 AGAIN",
		"< Header: Response 2",
		"```",
	})
	errors := make(chan error)

	tests := subject.Parse(body, errors)

	assertTest(t, parsers.Test{
		Name: "First request",
		Request: fillRequest("METHOD", "path", map[string]string{
			"Header": "Request",
		}, ""),
		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
			"Header": "Response",
		}, ""),
	}, tests, errors)

	assertTest(t, parsers.Test{
		Name: "Second request",
		Request: fillRequest("SECOND", "path", map[string]string{
			"Header": "Request 2",
		}, ""),
		Response: fillResponse(1234, "AGAIN", map[string]string{
			"Header": "Response 2",
		}, ""),
	}, tests, errors)
}

func TestMarkdownTaggedCodeBlock(t *testing.T) {
	subject := parsers.Markdown{}
	body := streamBody([]string{
		"# Markdown",
		"```tag",
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"```",
	})
	errors := make(chan error)

	tests := subject.Parse(body, errors)

	assertTest(t, parsers.Test{
		Name: "Markdown",
		Request: fillRequest("METHOD", "path", map[string]string{
			"Header": "Request",
		}, ""),
		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
			"Header": "Response",
		}, ""),
	}, tests, errors)
}

func TestMarkdownBlockQuote(t *testing.T) {
	subject := parsers.Markdown{}
	body := streamBody([]string{
		"# Markdown",
		"> Quoting some spec or something",
		"```",
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"```",
	})
	errors := make(chan error)

	tests := subject.Parse(body, errors)

	assertTest(t, parsers.Test{
		Name: "Markdown",
		Request: fillRequest("METHOD", "path", map[string]string{
			"Header": "Request",
		}, ""),
		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
			"Header": "Response",
		}, ""),
	}, tests, errors)
}

func TestMarkdownFallbackName(t *testing.T) {
	subject := parsers.Markdown{}
	body := streamBody([]string{
		"```",
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"```",
	})
	errors := make(chan error)

	tests := subject.Parse(body, errors)

	assertTest(t, parsers.Test{
		Name: "METHOD path -> 1337",
		Request: fillRequest("METHOD", "path", map[string]string{
			"Header": "Request",
		}, ""),
		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
			"Header": "Response",
		}, ""),
	}, tests, errors)
}
