package parsers_test

import (
	"testing"

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

func assertTest(t *testing.T, expected []string, specs chan []string, errors chan error) {
	select {
	case err := <-errors:
		assert.FailNow(t, "Error should not exist.", err)
	case spec := <-specs:
		assert.Equal(t, expected, spec)
	}
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

	specs := subject.Parse(body, errors)

	assertTest(t, []string{
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
	}, specs, errors)
	// TODO(schoon) - Assert that the channel is closed here.
}

// TODO: is this still valid?
// func TestMarkdownSplitReqRes(t *testing.T) {
// 	subject := parsers.Markdown{}
// 	body := streamBody([]string{
// 		"# Markdown",
// 		"```",
// 		"> METHOD path",
// 		"> Header: Request",
// 		"```",
// 		"```",
// 		"< PROTO 1337 STATUS TEXT",
// 		"< Header: Response",
// 		"```",
// 	})
// 	errors := make(chan error)

// 	specs := subject.Parse(body, errors)

// 	assertTest(t, builders.Test{
// 		Name: "Markdown",
// 		Request: fillRequest("METHOD", "path", map[string]string{
// 			"Header": "Request",
// 		}, ""),
// 		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
// 			"Header": "Response",
// 		}, ""),
// 	}, specs, errors)
// }

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

	specs := subject.Parse(body, errors)

	assertTest(t, []string{
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
	}, specs, errors)

	assertTest(t, []string{
		"> SECOND path",
		"> Header: Request 2",
		"< PROTO 1234 AGAIN",
		"< Header: Response 2",
	}, specs, errors)
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

	specs := subject.Parse(body, errors)

	assertTest(t, []string{
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
	}, specs, errors)
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

	specs := subject.Parse(body, errors)

	assertTest(t, []string{
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
	}, specs, errors)
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

	specs := subject.Parse(body, errors)

	assertTest(t, []string{
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
	}, specs, errors)
}
