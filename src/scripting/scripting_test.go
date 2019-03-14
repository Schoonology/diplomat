package scripting_test

import (
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
