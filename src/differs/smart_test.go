package differs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/diplomat/differs"
	"github.com/testdouble/diplomat/http"
)

func TestSmartSameNoHeaders(t *testing.T) {
	subject := differs.Smart{}
	diff, err := subject.Diff(http.NewResponse(200, "STATUS TEXT"), http.NewResponse(200, "STATUS TEXT"))

	assert.Nil(t, err)
	assert.Equal(t, "", diff)
}

func TestSmartSameFull(t *testing.T) {
	subject := differs.Smart{}
	actual := http.NewResponse(200, "OK")
	expected := http.NewResponse(200, "OK")

	actual.Headers["Some-Header"] = "Same!"
	expected.Headers["Some-Header"] = "Same!"

	diff, err := subject.Diff(expected, actual)

	assert.Nil(t, err)
	assert.Equal(t, "", diff)
}

func TestSmartWrongStatus(t *testing.T) {
	subject := differs.Smart{}
	diff, err := subject.Diff(http.NewResponse(204, "No Content"), http.NewResponse(200, "OK"))

	assert.Nil(t, err)
	assert.Equal(t, `Status:
	- 204 No Content
	+ 200 OK
`, diff)
}

func TestSmartWrongHeaderValue(t *testing.T) {
	subject := differs.Smart{}
	actual := http.NewResponse(200, "OK")
	expected := http.NewResponse(200, "OK")

	actual.Headers["Test-Header"] = "Actual"
	expected.Headers["Test-Header"] = "Expected"

	diff, err := subject.Diff(expected, actual)

	assert.Nil(t, err)
	assert.Equal(t, `Invalid Header: Test-Header
	- Expected
	+ Actual
`, diff)
}

func TestSmartExtraHeader(t *testing.T) {
	subject := differs.Smart{}
	actual := http.NewResponse(200, "OK")
	expected := http.NewResponse(200, "OK")

	actual.Headers["Test-Header"] = "Extra"

	diff, err := subject.Diff(expected, actual)

	assert.Nil(t, err)
	assert.Equal(t, "", diff)
}

func TestSmartMissingHeader(t *testing.T) {
	subject := differs.Smart{}
	actual := http.NewResponse(200, "OK")
	expected := http.NewResponse(200, "OK")

	expected.Headers["Test-Header"] = "Missing"

	diff, err := subject.Diff(expected, actual)

	assert.Nil(t, err)
	assert.Equal(t, `Missing Header: Test-Header
`, diff)
}

func TestSmartSimilarJson(t *testing.T) {
	subject := differs.Smart{}
	actual := http.NewResponse(200, "OK")
	expected := http.NewResponse(200, "OK")

	expected.Body = []byte(`{"key":"value"}`)
	actual.Headers["Content-Type"] = "application/json"
	actual.Body = []byte(`{
	"key": "value"
}`)

	diff, err := subject.Diff(expected, actual)

	assert.Nil(t, err)
	assert.Equal(t, "", diff)
}

func TestSmartWrongJson(t *testing.T) {
	subject := differs.Smart{}
	actual := http.NewResponse(200, "OK")
	expected := http.NewResponse(200, "OK")

	expected.Body = []byte(`{"key":"value"}`)
	actual.Headers["Content-Type"] = "application/json"
	actual.Body = []byte(`{"key": "wrong"}`)

	diff, err := subject.Diff(expected, actual)

	assert.Nil(t, err)
	assert.Equal(t, `Invalid Body:
root["key"]:
	-: "value"
	+: "wrong"
`, diff)
}

func TestSmartCorrectText(t *testing.T) {
	subject := differs.Smart{}
	actual := http.NewResponse(200, "OK")
	expected := http.NewResponse(200, "OK")

	expected.Body = []byte(`This is correct!`)
	actual.Headers["Content-Type"] = "plain/text"
	actual.Body = []byte(`This is correct!`)

	diff, err := subject.Diff(expected, actual)

	assert.Nil(t, err)
	assert.Equal(t, "", diff)
}

func TestSmartIncorrectText(t *testing.T) {
	subject := differs.Smart{}
	actual := http.NewResponse(200, "OK")
	expected := http.NewResponse(200, "OK")

	expected.Body = []byte(`This is correct!`)
	actual.Headers["Content-Type"] = "plain/text"
	actual.Body = []byte(`This is incorrect!`)

	diff, err := subject.Diff(expected, actual)

	assert.Nil(t, err)
	assert.Equal(t, `Invalid Body:
	-: This is incorrect!
	+: This is correct!
`, diff)
}

func TestSmartMultipleCustomValidators(t *testing.T) {
	subject := differs.Smart{}
	actual := http.NewResponse(200, "OK")
	expected := http.NewResponse(200, "OK")

	expected.Body = []byte(`{? regexp("^te") ?}{? regexp("st$") ?}`)
	actual.Body = []byte(`tempest`)

	diff, err := subject.Diff(expected, actual)

	assert.Nil(t, err)
	assert.Equal(t, "", diff)
}

func TestSmartMultipleCustomValidatorsWhitespace(t *testing.T) {
	subject := differs.Smart{}
	actual := http.NewResponse(200, "OK")
	expected := http.NewResponse(200, "OK")

	expected.Body = []byte(`  {? regexp("^te") ?}	{? regexp("st$") ?}  `)
	actual.Body = []byte(`tempest`)

	diff, err := subject.Diff(expected, actual)

	assert.Nil(t, err)
	assert.Equal(t, "", diff)
}
