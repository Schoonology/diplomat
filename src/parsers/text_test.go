package parsers_test

import (
	"testing"

	"github.com/schoonology/diplomat/parsers"
	"github.com/stretchr/testify/assert"
)

func TestLoadText(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := []string{
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
	}

	paragraphs := subject.Parse(body)

	assert.Equal(t, parsers.Paragraph{
		Name: "",
		Body: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
		},
		LineNumber: 1,
	}, paragraphs[0])
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

// 	paragraphs := subject.Parse(body, errors)

// 	assertTest(t, []string{
// 		"> METHOD path",
// 		"> Header: Request",
// 		"< PROTO 1337 STATUS TEXT",
// 		"< Header: Response",
// 	}, paragraphs, errors)

// 	assertTest(t, []string{
// 		"> SECOND path",
// 		"> Header: Request 2",
// 		"< PROTO 1234 AGAIN",
// 		"< Header: Response 2",
// 	}, paragraphs, errors)
// }

func TestSingleLineBody(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := []string{
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"<",
		"Some response body",
	}

	paragraphs := subject.Parse(body)

	assert.Equal(t, parsers.Paragraph{
		Name: "",
		Body: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"<",
			"Some response body",
		},
		LineNumber: 1,
	}, paragraphs[0])
}

func TestMultiLineBodyWithIndentation(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := []string{
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"<",
		"This is the first line",
		"  This is the second line",
	}

	paragraphs := subject.Parse(body)

	assert.Equal(t, parsers.Paragraph{
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
		LineNumber: 1,
	}, paragraphs[0])
}

func TestMissingBracket(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := []string{
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"Some response body",
	}

	paragraphs := subject.Parse(body)

	assert.Equal(t, parsers.Paragraph{
		Name: "",
		Body: []string{
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"Some response body",
		},
		LineNumber: 1,
	}, paragraphs[0])
}

func TestCommentsAboveparagraph(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := []string{
		"comments!!",
		"> METHOD path",
		"> Header: Request",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"Some response body",
	}

	paragraphs := subject.Parse(body)

	assert.Equal(t, parsers.Paragraph{
		Name: "",
		Body: []string{
			"comments!!",
			"> METHOD path",
			"> Header: Request",
			"< PROTO 1337 STATUS TEXT",
			"< Header: Response",
			"Some response body",
		},
		LineNumber: 1,
	}, paragraphs[0])
}

func TestRequestBody(t *testing.T) {
	subject := parsers.PlainTextParser{}
	body := []string{
		"> METHOD path",
		"> Header: Request",
		"Some request body",
		"< PROTO 1337 STATUS TEXT",
		"< Header: Response",
		"<",
		"Some response body",
	}

	paragraphs := subject.Parse(body)

	assert.Equal(t, parsers.Paragraph{
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
		LineNumber: 1,
	}, paragraphs[0])
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

// 	paragraphs := subject.Parse(body, errors)

// 	assertTest(t, []string{
// 		"> FIRST path",
// 		"> Header: Request",
// 		"First request body",
// 		"< PROTO 1337 STATUS TEXT",
// 		"< Header: Response",
// 		"<",
// 		"First response body",
// 	}, paragraphs, errors)

// 	assertTest(t, []string{
// 		"> SECOND path",
// 		"> Header: Request",
// 		"Second request body",
// 		"< PROTO 1337 STATUS TEXT",
// 		"< Header: Response",
// 		"<",
// 		"Second response body",
// 	}, paragraphs, errors)
// }
