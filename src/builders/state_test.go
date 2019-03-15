package builders_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/diplomat/builders"
	"github.com/testdouble/diplomat/errors"
	"github.com/testdouble/diplomat/parsers"
)

func TestNoBody(t *testing.T) {
	subject := builders.State{}
	body := parsers.Spec{
		Name: "",
		Body: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
		},
	}

	test, err := subject.Build(body)

	assertTest(t, builders.Test{
		Name: "METHOD path -> 1337",
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
		Name: "",
		Body: []string{
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
		Name: "METHOD path -> 1337",
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
		Name: "",
		Body: []string{
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
		Name: "METHOD path -> 1337",
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
		Name: "",
		Body: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"Some response body",
		},
	}

	test, err := subject.Build(body)

	assertTest(t, builders.Test{
		Name: "METHOD path -> 1337",
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
		Name: "",
		Body: []string{
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
		Name: "",
		Body: []string{
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
		Name: "METHOD path -> 1337",
		Request: fillRequest("METHOD", "path", map[string]string{
			"Header": "Request",
		}, "Some request body\n"),
		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
			"Header": "Response",
		}, "Some response body\n"),
	}, test, err)
}

func TestRequestOnly(t *testing.T) {
	subject := builders.State{}
	body := parsers.Spec{
		Name: "",
		Body: []string{
			"> METHOD path",
			"> Header: Request",
		},
	}

	_, err := subject.Build(body)

	assert.Equal(t, &errors.MissingResponse{}, err)
}

func TestResponseOnly(t *testing.T) {
	subject := builders.State{}
	body := parsers.Spec{
		Name: "",
		Body: []string{
			"< HTTP/1.1 200 OK",
			"< Content-Length: 0",
			"<",
		},
	}

	_, err := subject.Build(body)

	assert.Equal(t, &errors.MissingRequest{}, err)
}

func TestNoVersion(t *testing.T) {
	subject := builders.State{}
	body := parsers.Spec{
		Name: "",
		Body: []string{
			"> METHOD path",
			"> Header: Request",
			"< 1337 STATUS TEXT",
			"< Header: Response",
		},
	}

	test, err := subject.Build(body)

	assertTest(t, builders.Test{
		Name: "METHOD path -> 1337",
		Request: fillRequest("METHOD", "path", map[string]string{
			"Header": "Request",
		}, ""),
		Response: fillResponse(1337, "STATUS TEXT", map[string]string{
			"Header": "Response",
		}, ""),
	}, test, err)
}
