package parsers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/diplomat/errors"
	"github.com/testdouble/diplomat/http"
	"github.com/testdouble/diplomat/parsers"
)

func TestLoadText(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := streamBody([]string{
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
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

func TestLoadEmpty(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := streamBody([]string{})
	errors := make(chan error)

	tests := subject.Parse(body, errors)

	_, more := <-tests
	assert.False(t, more)
}

func TestLoadDouble(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := streamBody([]string{
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"> SECOND path",
		"> Header: Request 2",
		"< PROTO 1234 AGAIN",
		"< Header: Response 2",
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

	assertTest(t, parsers.Test{
		Name: "SECOND path -> 1234",
		Request: fillRequest("SECOND", "path", map[string]string{
			"Header": "Request 2",
		}, ""),
		Response: fillResponse(1234, "AGAIN", map[string]string{
			"Header": "Response 2",
		}, ""),
	}, tests, errors)
}

func TestSingleLineBody(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := streamBody([]string{
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"<",
		"Some response body",
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
		}, "Some response body\n"),
	}, tests, errors)
}

func TestMultiLineBodyWithIndentation(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := streamBody([]string{
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"<",
		"This is the first line",
		"  This is the second line",
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
		}, "This is the first line\n  This is the second line\n"),
	}, tests, errors)
}

func TestMissingBracket(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := streamBody([]string{
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"Some response body",
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
		}, "Some response body\n"),
	}, tests, errors)
}

func TestMultiLineBodyWithAngleBracket(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := streamBody([]string{
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"<",
		"This is the first line",
		"<   This is the second line",
	})
	errorChannel := make(chan error)

	_ = subject.Parse(body, errorChannel)

	err := <-errorChannel
	_, ok := err.(*errors.MissingRequest)
	assert.True(t, ok)
}

func TestCommentsAboveSpec(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := streamBody([]string{
		"comments!!",
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"Some response body",
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
		}, "Some response body\n"),
	}, tests, errors)
}

func TestRequestBody(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := streamBody([]string{
		"> METHOD path",
		"> Header: Request",
		"Some request body",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"<",
		"Some response body",
	})
	errors := make(chan error)

	tests := subject.Parse(body, errors)

	assertTest(t, parsers.Test{
		Name: "METHOD path -> 1337",
		Request: fillRequest("METHOD", "path", map[string]string{
			"Header": "Request",
		}, "Some request body\n"),
		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
			"Header": "Response",
		}, "Some response body\n"),
	}, tests, errors)
}

func TestKitchenSink(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := streamBody([]string{
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
	})
	errors := make(chan error)

	tests := subject.Parse(body, errors)

	assertTest(t, parsers.Test{
		Name: "FIRST path -> 1337",
		Request: &http.Request{
			Method:  "FIRST",
			Path:    "path",
			Headers: map[string]string{"Header": "Request"},
			Body:    []byte("First request body\n"),
		},
		Response: &http.Response{
			StatusCode: 1337,
			StatusText: "STATUS TEXT",
			Headers:    map[string]string{"Header": "Response"},
			Body:       []byte("First response body\n"),
		},
	}, tests, errors)

	assertTest(t, parsers.Test{
		Name: "SECOND path -> 1337",
		Request: &http.Request{
			Method:  "SECOND",
			Path:    "path",
			Headers: map[string]string{"Header": "Request"},
			Body:    []byte("Second request body\n"),
		},
		Response: &http.Response{
			StatusCode: 1337,
			StatusText: "STATUS TEXT",
			Headers:    map[string]string{"Header": "Response"},
			Body:       []byte("Second response body\n"),
		},
	}, tests, errors)
}
