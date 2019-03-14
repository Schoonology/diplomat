package builders_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/diplomat/builders"
	"github.com/testdouble/diplomat/errors"
	"github.com/testdouble/diplomat/http"
	"github.com/testdouble/diplomat/parsers"
)

func assertTest(t *testing.T, expected builders.Test, actual builders.Test, err error) {
	if err != nil {
		assert.FailNow(t, "Error should not exist.")
		return
	}
	assert.Equal(t, expected, actual)
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

func TestNoBody(t *testing.T) {
	subject := builders.State{}
	body := parsers.Spec{
		"",
		[]string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
		},
	}

	test, err := subject.Build(body)

	assertTest(t, builders.Test{
		Request: fillRequest("METHOD", "path", map[string]string{
			"Header": "Request",
		}, ""),
		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
			"Header": "Response",
		}, ""),
	}, test, err)
}

func TestSingleLineBody(t *testing.T) {
	subject := builders.State{}
	body := parsers.Spec{
		"",
		[]string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"<",
			"Some response body",
		},
	}

	test, err := subject.Build(body)

	assertTest(t, builders.Test{
		Request: fillRequest("METHOD", "path", map[string]string{
			"Header": "Request",
		}, ""),
		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
			"Header": "Response",
		}, "Some response body\n"),
	}, test, err)
}

func TestMultiLineBodyWithIndentation(t *testing.T) {
	subject := builders.State{}
	body := parsers.Spec{
		"",
		[]string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"<",
			"This is the first line",
			"  This is the second line",
		},
	}

	test, err := subject.Build(body)

	assertTest(t, builders.Test{
		Request: fillRequest("METHOD", "path", map[string]string{
			"Header": "Request",
		}, ""),
		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
			"Header": "Response",
		}, "This is the first line\n  This is the second line\n"),
	}, test, err)
}

func TestMissingBracket(t *testing.T) {
	subject := builders.State{}
	body := parsers.Spec{
		"",
		[]string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"Some response body",
		},
	}

	test, err := subject.Build(body)

	assertTest(t, builders.Test{
		Request: fillRequest("METHOD", "path", map[string]string{
			"Header": "Request",
		}, ""),
		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
			"Header": "Response",
		}, "Some response body\n"),
	}, test, err)
}

func TestMultiLineBodyWithAngleBracket(t *testing.T) {
	subject := builders.State{}
	body := parsers.Spec{
		"",
		[]string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"<",
			"This is the first line",
			"<   This is the second line",
		},
	}

	_, err := subject.Build(body)
	_, ok := err.(*errors.MissingRequest)

	assert.True(t, ok)
}

func TestRequestBody(t *testing.T) {
	subject := builders.State{}
	body := parsers.Spec{
		"",
		[]string{
			"> METHOD path",
			"> Header: Request",
			"Some request body",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"<",
			"Some response body",
		},
	}

	test, err := subject.Build(body)

	assertTest(t, builders.Test{
		Request: fillRequest("METHOD", "path", map[string]string{
			"Header": "Request",
		}, "Some request body\n"),
		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
			"Header": "Response",
		}, "Some response body\n"),
	}, test, err)
}
