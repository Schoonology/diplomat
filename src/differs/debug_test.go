package differs_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/diplomat/differs"
	"github.com/testdouble/diplomat/http"
)

func TestDiffSameNoHeaders(t *testing.T) {
	subject := differs.Debug{}
	diff, err := subject.Diff(http.NewResponse(200, "STATUS TEXT"), http.NewResponse(200, "STATUS TEXT"))

	assert.Nil(t, err)
	assert.Equal(t, "", diff)
}

func TestDiffSameFull(t *testing.T) {
	subject := differs.Debug{}
	actual := http.NewResponse(200, "OK")
	expected := http.NewResponse(200, "OK")

	actual.Headers["Some-Header"] = "Same!"
	expected.Headers["Some-Header"] = "Same!"

	diff, err := subject.Diff(expected, actual)

	assert.Nil(t, err)
	assert.Equal(t, "", diff)
}

func TestDiffWrongStatus(t *testing.T) {
	subject := differs.Debug{}
	diff, err := subject.Diff(http.NewResponse(204, "No Content"), http.NewResponse(200, "OK"))

	assert.Nil(t, err)
	assert.Equal(t, `{*http.Response}.StatusCode:
	-: 204
	+: 200
{*http.Response}.StatusText:
	-: "No Content"
	+: "OK"
`, diff)
}

func TestDiffWrongHeaderValue(t *testing.T) {
	subject := differs.Debug{}
	actual := http.NewResponse(200, "OK")
	expected := http.NewResponse(200, "OK")

	actual.Headers["Test-Header"] = "Actual"
	expected.Headers["Test-Header"] = "Expected"

	diff, err := subject.Diff(expected, actual)

	assert.Nil(t, err)
	assert.Equal(t, `{*http.Response}.Headers["Test-Header"]:
	-: "Expected"
	+: "Actual"
`, diff)
}

func TestDiffExtraHeader(t *testing.T) {
	subject := differs.Debug{}
	actual := http.NewResponse(200, "OK")
	expected := http.NewResponse(200, "OK")

	actual.Headers["Test-Header"] = "Extra"

	diff, err := subject.Diff(expected, actual)

	assert.Nil(t, err)
	assert.Equal(t, `{*http.Response}.Headers["Test-Header"]:
	-: <non-existent>
	+: "Extra"
`, diff)
}

func TestDiffMissingHeader(t *testing.T) {
	subject := differs.Debug{}
	actual := http.NewResponse(200, "OK")
	expected := http.NewResponse(200, "OK")

	expected.Headers["Test-Header"] = "Missing"

	diff, err := subject.Diff(expected, actual)

	assert.Nil(t, err)
	assert.Equal(t, `{*http.Response}.Headers["Test-Header"]:
	-: "Missing"
	+: <non-existent>
`, diff)
}
