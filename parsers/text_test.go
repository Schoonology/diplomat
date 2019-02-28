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
