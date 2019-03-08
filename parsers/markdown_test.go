package parsers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/diplomat/loaders"
	"github.com/testdouble/diplomat/parsers"
)

func TestMarkdownText(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.Markdown{}
	body := loaders.Body{
		Lines: []string{
			"# Markdown",
			"```",
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"```",
		},
	}

	spec, err := subject.Parse(&body)
	assert.Nil(err)
	assert.NotNil(spec)
	assert.Equal(1, len(spec.Tests))

	test := spec.Tests[0]
	assert.Equal("Markdown", test.Name)
	assert.Equal("METHOD", test.Request.Method)
	assert.Equal("path", test.Request.Path)
	assert.Equal("Request", test.Request.Headers["Header"])
	assert.Equal(1337, test.Response.StatusCode)
	assert.Equal("STATUS TEXT", test.Response.StatusText)
	assert.Equal("Response", test.Response.Headers["Header"])
}

func TestMarkdownSplitReqRes(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.Markdown{}
	body := loaders.Body{
		Lines: []string{
			"# Markdown",
			"```",
			"> METHOD path",
			"> Header: Request",
			"```",
			"```",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"```",
		},
	}

	spec, err := subject.Parse(&body)
	assert.Nil(err)
	assert.NotNil(spec)
	assert.Equal(1, len(spec.Tests))

	test := spec.Tests[0]
	assert.Equal("Markdown", test.Name)
	assert.Equal("METHOD", test.Request.Method)
	assert.Equal("path", test.Request.Path)
	assert.Equal("Request", test.Request.Headers["Header"])
	assert.Equal(1337, test.Response.StatusCode)
	assert.Equal("STATUS TEXT", test.Response.StatusText)
	assert.Equal("Response", test.Response.Headers["Header"])
}

func TestMarkdownDouble(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.Markdown{}
	body := loaders.Body{
		Lines: []string{
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
		},
	}

	spec, err := subject.Parse(&body)
	assert.Nil(err)
	assert.NotNil(spec)
	assert.Equal(2, len(spec.Tests))

	test := spec.Tests[0]
	assert.Equal("First request", test.Name)
	assert.Equal("METHOD", test.Request.Method)
	assert.Equal("path", test.Request.Path)
	assert.Equal("Request", test.Request.Headers["Header"])
	assert.Equal(1337, test.Response.StatusCode)
	assert.Equal("STATUS TEXT", test.Response.StatusText)
	assert.Equal("Response", test.Response.Headers["Header"])

	test = spec.Tests[1]
	assert.Equal("Second request", test.Name)
	assert.Equal("SECOND", test.Request.Method)
	assert.Equal("path", test.Request.Path)
	assert.Equal("Request 2", test.Request.Headers["Header"])
	assert.Equal(1234, test.Response.StatusCode)
	assert.Equal("AGAIN", test.Response.StatusText)
	assert.Equal("Response 2", test.Response.Headers["Header"])
}

func TestMarkdownTaggedCodeBlock(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.Markdown{}
	body := loaders.Body{
		Lines: []string{
			"# Markdown",
			"```tag",
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"```",
		},
	}

	spec, err := subject.Parse(&body)
	assert.Nil(err)
	assert.NotNil(spec)
	assert.Equal(1, len(spec.Tests))

	test := spec.Tests[0]
	assert.Equal("Markdown", test.Name)
	assert.Equal("METHOD", test.Request.Method)
	assert.Equal("path", test.Request.Path)
	assert.Equal("Request", test.Request.Headers["Header"])
	assert.Equal(1337, test.Response.StatusCode)
	assert.Equal("STATUS TEXT", test.Response.StatusText)
	assert.Equal("Response", test.Response.Headers["Header"])
}

func TestMarkdownBlockQuote(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.Markdown{}
	body := loaders.Body{
		Lines: []string{
			"# Markdown",
			"> Quoting some spec or something",
			"```",
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"```",
		},
	}

	spec, err := subject.Parse(&body)
	assert.Nil(err)
	assert.NotNil(spec)
	assert.Equal(1, len(spec.Tests))

	test := spec.Tests[0]
	assert.Equal("Markdown", test.Name)
	assert.Equal("METHOD", test.Request.Method)
	assert.Equal("path", test.Request.Path)
	assert.Equal("Request", test.Request.Headers["Header"])
	assert.Equal(1337, test.Response.StatusCode)
	assert.Equal("STATUS TEXT", test.Response.StatusText)
	assert.Equal("Response", test.Response.Headers["Header"])
}

func TestMarkdownFallbackName(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.Markdown{}
	body := loaders.Body{
		Lines: []string{
			"```",
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"```",
		},
	}

	spec, err := subject.Parse(&body)
	assert.Nil(err)
	assert.NotNil(spec)
	assert.Equal(1, len(spec.Tests))

	test := spec.Tests[0]
	assert.Equal("METHOD path -> 1337", test.Name)
	assert.Equal("METHOD", test.Request.Method)
	assert.Equal("path", test.Request.Path)
	assert.Equal("Request", test.Request.Headers["Header"])
	assert.Equal(1337, test.Response.StatusCode)
	assert.Equal("STATUS TEXT", test.Response.StatusText)
	assert.Equal("Response", test.Response.Headers["Header"])
}
