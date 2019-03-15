package parsers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

	specs := subject.Parse(body, errors)

	assertTest(t, parsers.Spec{
		Name: "",
		Body: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
		},
	}, specs, errors)
}

func TestLoadEmpty(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := streamBody([]string{})
	errors := make(chan error)

	specs := subject.Parse(body, errors)

	_, more := <-specs
	assert.False(t, more)
}

// TODO: Is this test valid?
// func TestLoadDouble(t *testing.T) {
// 	subject := parsers.PlainTextParser{}
// 	body := streamBody([]string{
// 		"> METHOD path",
// 		"> Header: Request",
// 		"< PROTO 1337 STATUS TEXT",
// 		"< Header: Response",
// 		"> SECOND path",
// 		"> Header: Request 2",
// 		"< PROTO 1234 AGAIN",
// 		"< Header: Response 2",
// 	})
// 	errors := make(chan error)

// 	specs := subject.Parse(body, errors)

// 	assertTest(t, []string{
// 		"> METHOD path",
// 		"> Header: Request",
// 		"< PROTO 1337 STATUS TEXT",
// 		"< Header: Response",
// 	}, specs, errors)

// 	assertTest(t, []string{
// 		"> SECOND path",
// 		"> Header: Request 2",
// 		"< PROTO 1234 AGAIN",
// 		"< Header: Response 2",
// 	}, specs, errors)
// }

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

	specs := subject.Parse(body, errors)

	assertTest(t, parsers.Spec{
		Name: "",
		Body: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"<",
			"Some response body",
		},
	}, specs, errors)
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

	specs := subject.Parse(body, errors)

	assertTest(t, parsers.Spec{
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
	}, specs, errors)
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

	specs := subject.Parse(body, errors)

	assertTest(t, parsers.Spec{
		Name: "",
		Body: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"Some response body",
		},
	}, specs, errors)
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

	specs := subject.Parse(body, errors)

	assertTest(t, parsers.Spec{
		Name: "",
		Body: []string{
			"comments!!",
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"Some response body",
		},
	}, specs, errors)
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

	specs := subject.Parse(body, errors)

	assertTest(t, parsers.Spec{
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
	}, specs, errors)
}

// TODO: Is this test valid?
// func TestKitchenSink(t *testing.T) {
// 	subject := parsers.PlainTextParser{}
// 	body := streamBody([]string{
// 		"> FIRST path",
// 		"> Header: Request",
// 		"First request body",
// 		"< PROTO 1337 STATUS TEXT",
// 		"< Header: Response",
// 		"<",
// 		"First response body",
// 		"> SECOND path",
// 		"> Header: Request",
// 		"Second request body",
// 		"< PROTO 1337 STATUS TEXT",
// 		"< Header: Response",
// 		"<",
// 		"Second response body",
// 	})
// 	errors := make(chan error)

// 	specs := subject.Parse(body, errors)

// 	assertTest(t, []string{
// 		"> FIRST path",
// 		"> Header: Request",
// 		"First request body",
// 		"< PROTO 1337 STATUS TEXT",
// 		"< Header: Response",
// 		"<",
// 		"First response body",
// 	}, specs, errors)

// 	assertTest(t, []string{
// 		"> SECOND path",
// 		"> Header: Request",
// 		"Second request body",
// 		"< PROTO 1337 STATUS TEXT",
// 		"< Header: Response",
// 		"<",
// 		"Second response body",
// 	}, specs, errors)
// }
