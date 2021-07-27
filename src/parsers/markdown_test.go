package parsers_test

import (
	"testing"

	"github.com/schoonology/diplomat/parsers"
	"github.com/stretchr/testify/assert"
)

func TestMarkdownText(t *testing.T) {
	subject := parsers.Markdown{}
	body := []string{
		"# Markdown",
		"```",
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"```",
	}

	paragraphs := subject.Parse(body)

	assert.Equal(t, parsers.Paragraph{
		Name: "Markdown",
		Body: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
		},
		LineNumber: 3,
	}, paragraphs[0])
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

// 	paragraphs := subject.Parse(body, errors)

// 	assertTest(t, builders.Test{
// 		Name: "Markdown",
// 		Request: fillRequest("METHOD", "path", map[string]string{
// 			"Header": "Request",
// 		}, ""),
// 		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
// 			"Header": "Response",
// 		}, ""),
// 	}, paragraphs, errors)
// }

func TestMarkdownDouble(t *testing.T) {
	subject := parsers.Markdown{}
	body := []string{
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
	}

	paragraphs := subject.Parse(body)

	assert.Equal(t, parsers.Paragraph{
		Name: "First request",
		Body: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
		},
		LineNumber: 4,
	}, paragraphs[0])

	assert.Equal(t, parsers.Paragraph{
		Name: "Second request",
		Body: []string{
			"> SECOND path",
			"> Header: Request 2",
			"< PROTO 1234 AGAIN",
			"< Header: Response 2",
		},
		LineNumber: 11,
	}, paragraphs[1])
}

func TestMarkdownTaggedCodeBlock(t *testing.T) {
	subject := parsers.Markdown{}
	body := []string{
		"# Markdown",
		"```tag",
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"```",
	}

	paragraphs := subject.Parse(body)

	assert.Equal(t, parsers.Paragraph{
		Name: "Markdown",
		Body: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
		},
		LineNumber: 3,
	}, paragraphs[0])
}

func TestMarkdownBlockQuote(t *testing.T) {
	subject := parsers.Markdown{}
	body := []string{
		"# Markdown",
		"> Quoting some paragraph or something",
		"```",
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"```",
	}

	paragraphs := subject.Parse(body)

	assert.Equal(t, parsers.Paragraph{
		Name: "Markdown",
		Body: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
		},
		LineNumber: 4,
	}, paragraphs[0])
}
