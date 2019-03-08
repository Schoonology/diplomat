package parsers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/http-assertion-tool/loaders"
	"github.com/testdouble/http-assertion-tool/parsers"
)

func TestLoadText(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.PlainTextParser{}
	body := loaders.Body{
		Lines: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
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

func TestLoadEmpty(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.PlainTextParser{}
	body := loaders.Body{
		Lines: []string{},
	}

	spec, err := subject.Parse(&body)
	assert.Nil(err)
	assert.NotNil(spec)
	assert.Equal(0, len(spec.Tests))
}

func TestLoadDouble(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.PlainTextParser{}
	body := loaders.Body{
		Lines: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"> SECOND path",
			"> Header: Request 2",
			"< PROTO 1234 AGAIN",
			"< Header: Response 2",
		},
	}

	spec, err := subject.Parse(&body)
	assert.Nil(err)
	assert.NotNil(spec)
	assert.Equal(2, len(spec.Tests))

	test := spec.Tests[0]
	assert.Equal("METHOD", test.Request.Method)
	assert.Equal("path", test.Request.Path)
	assert.Equal("Request", test.Request.Headers["Header"])
	assert.Equal(1337, test.Response.StatusCode)
	assert.Equal("STATUS TEXT", test.Response.StatusText)
	assert.Equal("Response", test.Response.Headers["Header"])

	test = spec.Tests[1]
	assert.Equal("SECOND", test.Request.Method)
	assert.Equal("path", test.Request.Path)
	assert.Equal("Request 2", test.Request.Headers["Header"])
	assert.Equal(1234, test.Response.StatusCode)
	assert.Equal("AGAIN", test.Response.StatusText)
	assert.Equal("Response 2", test.Response.Headers["Header"])
}

func TestSingleLineBody(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.PlainTextParser{}
	body := loaders.Body{
		Lines: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"<",
			"Some response body",
		},
	}

	spec, err := subject.Parse(&body)
	assert.Nil(err)
	assert.NotNil(spec)
	assert.Equal(1, len(spec.Tests))

	test := spec.Tests[0]
	assert.Equal("METHOD", test.Request.Method)
	assert.Equal("path", test.Request.Path)
	assert.Equal("Request", test.Request.Headers["Header"])
	assert.Equal(1337, test.Response.StatusCode)
	assert.Equal("STATUS TEXT", test.Response.StatusText)
	assert.Equal("Response", test.Response.Headers["Header"])
	assert.Equal("Some response body\n", string(test.Response.Body))
}

func TestMultiLineBodyWithIndentation(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.PlainTextParser{}
	body := loaders.Body{
		Lines: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"<",
			"This is the first line",
			"  This is the second line",
		},
	}

	spec, err := subject.Parse(&body)
	assert.Nil(err)
	assert.NotNil(spec)
	assert.Equal(1, len(spec.Tests))

	test := spec.Tests[0]
	assert.Equal("METHOD", test.Request.Method)
	assert.Equal("path", test.Request.Path)
	assert.Equal("Request", test.Request.Headers["Header"])
	assert.Equal(1337, test.Response.StatusCode)
	assert.Equal("STATUS TEXT", test.Response.StatusText)
	assert.Equal("Response", test.Response.Headers["Header"])
	assert.Equal("This is the first line\n  This is the second line\n", string(test.Response.Body))
}

func TestMissingBracket(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.PlainTextParser{}
	body := loaders.Body{
		Lines: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"Some response body",
		},
	}

	spec, err := subject.Parse(&body)
	assert.Nil(err)
	assert.NotNil(spec)

	test := spec.Tests[0]
	assert.Equal("Some response body\n", string(test.Response.Body))
}

func TestMultiLineBodyWithAngleBracket(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.PlainTextParser{}
	body := loaders.Body{
		Lines: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"<",
			"This is the first line",
			"<   This is the second line",
		},
	}

	spec, err := subject.Parse(&body)
	assert.Nil(err)
	assert.NotNil(spec)
	assert.Equal(1, len(spec.Tests))

	test := spec.Tests[0]
	assert.Equal("METHOD", test.Request.Method)
	assert.Equal("path", test.Request.Path)
	assert.Equal("Request", test.Request.Headers["Header"])
	assert.Equal(1337, test.Response.StatusCode)
	assert.Equal("STATUS TEXT", test.Response.StatusText)
	assert.Equal("Response", test.Response.Headers["Header"])
	assert.Equal("This is the first line\n<   This is the second line\n", string(test.Response.Body))
}

func TestCommentsAboveSpec(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.PlainTextParser{}
	body := loaders.Body{
		Lines: []string{
			"comments!!",
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"<",
			"This is the first line",
			"<   This is the second line",
		},
	}

	spec, err := subject.Parse(&body)
	assert.Nil(err)
	assert.NotNil(spec)
}

func TestRequestBody(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.PlainTextParser{}
	body := loaders.Body{
		Lines: []string{
			"> METHOD path",
			"> Header: Request",
			"Some request body",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"<",
			"Some response body",
		},
	}

	spec, err := subject.Parse(&body)
	assert.Nil(err)
	assert.NotNil(spec)
	assert.Equal(1, len(spec.Tests))

	test := spec.Tests[0]
	assert.Equal("METHOD", test.Request.Method)
	assert.Equal("path", test.Request.Path)
	assert.Equal("Request", test.Request.Headers["Header"])
	assert.Equal("Some request body\n", string(test.Request.Body))
	assert.Equal(1337, test.Response.StatusCode)
	assert.Equal("STATUS TEXT", test.Response.StatusText)
	assert.Equal("Response", test.Response.Headers["Header"])
	assert.Equal("Some response body\n", string(test.Response.Body))
}

func TestKitchenSink(t *testing.T) {
	assert := assert.New(t)
	subject := parsers.PlainTextParser{}
	body := loaders.Body{
		Lines: []string{
			"> FIRST path",
			"> Header: Request",
			"First request body",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"<",
			"First response body",
			"> SECOND path",
			"> Header: Request",
			"Second request body",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"<",
			"Second response body",
		},
	}

	spec, err := subject.Parse(&body)
	assert.Nil(err)
	assert.NotNil(spec)
	assert.Equal(2, len(spec.Tests))

	test := spec.Tests[0]
	assert.Equal("FIRST", test.Request.Method)
	assert.Equal("path", test.Request.Path)
	assert.Equal("Request", test.Request.Headers["Header"])
	assert.Equal("First request body\n", string(test.Request.Body))
	assert.Equal(1337, test.Response.StatusCode)
	assert.Equal("STATUS TEXT", test.Response.StatusText)
	assert.Equal("Response", test.Response.Headers["Header"])
	assert.Equal("First response body\n", string(test.Response.Body))

	test = spec.Tests[1]
	assert.Equal("SECOND", test.Request.Method)
	assert.Equal("path", test.Request.Path)
	assert.Equal("Request", test.Request.Headers["Header"])
	assert.Equal("Second request body\n", string(test.Request.Body))
	assert.Equal(1337, test.Response.StatusCode)
	assert.Equal("STATUS TEXT", test.Response.StatusText)
	assert.Equal("Response", test.Response.Headers["Header"])
	assert.Equal("Second response body\n", string(test.Response.Body))
}
