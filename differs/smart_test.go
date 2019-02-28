package differs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/http-assertion-tool/differs"
	"github.com/testdouble/http-assertion-tool/http"
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
