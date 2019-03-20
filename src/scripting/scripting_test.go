package scripting_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/diplomat/scripting"
)

func Test__Test(t *testing.T) {
	result, err := scripting.RunPipeline("__test")

	assert.Nil(t, err)
	assert.Equal(t, "this is a test", result)
}

func TestIsTrue(t *testing.T) {
	result, err := scripting.RunValidator("is_test", "test")

	assert.Nil(t, err)
	assert.True(t, result)
}

func TestJsonSchema(t *testing.T) {
	result, err := scripting.RunValidator(`json_schema([[{
		"type": "object",
		"properties": {
			"test": {
				"type": "boolean"
			}
		}
	}]])`, `{
		"test": true
	}`)

	assert.Nil(t, err)
	assert.True(t, result)
}

func Test__RoundTripJSON(t *testing.T) {
	json := `{"key":"value"}`

	result, err := scripting.RunPipeline(
		fmt.Sprintf("json.encode(json.decode('%s'))", json),
	)

	assert.Nil(t, err)
	assert.Equal(t, json, result)
}

func Test__BasicHTTP(t *testing.T) {
	result, err := scripting.RunPipeline(`http.post("http://localhost:7357/anything", {
		headers = {
			['Basic-Test'] = "yup",
		},
	}).body`)

	assert.Nil(t, err)

	match, err := regexp.MatchString(`"method": "POST"`, result)
	assert.Nil(t, err)
	assert.True(t, match)

	match, err = regexp.MatchString(`"Basic-Test": "yup"`, result)
	assert.Nil(t, err)
	assert.True(t, match)

	match, err = regexp.MatchString(`"User-Agent": "Diplomat/0.0.1"`, result)
	assert.Nil(t, err)
	assert.True(t, match)
}
