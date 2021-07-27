package builders_test

import (
	"testing"

	"github.com/schoonology/diplomat/builders"
	"github.com/schoonology/diplomat/errors"
	"github.com/schoonology/diplomat/parsers"
	"github.com/stretchr/testify/assert"
)

func TestNoBody(t *testing.T) {
	subject := builders.State{}
	body := parsers.Paragraph{
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
	body := parsers.Paragraph{
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
	body := parsers.Paragraph{
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
	body := parsers.Paragraph{
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
	body := parsers.Paragraph{
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
		LineNumber: 1,
	}

	_, err := subject.Build(body)

	assert.Equal(t, &errors.BuildError{
		Err:        &errors.MissingRequest{},
		LineNumber: 1,
	}, err)
}

func TestRequestBody(t *testing.T) {
	subject := builders.State{}
	body := parsers.Paragraph{
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
	body := parsers.Paragraph{
		Name: "",
		Body: []string{
			"> METHOD path",
			"> Header: Request",
		},
		LineNumber: 1,
	}

	_, err := subject.Build(body)

	assert.Equal(t, &errors.BuildError{
		Err:        &errors.MissingResponse{},
		LineNumber: 1,
	}, err)
}

func TestResponseOnly(t *testing.T) {
	subject := builders.State{}
	body := parsers.Paragraph{
		Name: "",
		Body: []string{
			"< HTTP/1.1 200 OK",
			"< Content-Length: 0",
			"<",
		},
		LineNumber: 1,
	}

	_, err := subject.Build(body)

	assert.Equal(t, &errors.BuildError{
		Err:        &errors.MissingRequest{},
		LineNumber: 1,
	}, err)
}

func TestNoVersion(t *testing.T) {
	subject := builders.State{}
	body := parsers.Paragraph{
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
