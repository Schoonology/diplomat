package scripting_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/http-assertion-tool/scripting"
)

func Test__Test(t *testing.T) {
	result, err := scripting.RunPipeline("__test")

	assert.Nil(t, err)
	assert.Equal(t, "this is a test", result)
}
