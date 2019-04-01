package colors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/diplomat/colors"
)

func TestIntegrationPaintRed(t *testing.T) {
	subject := colors.DefaultColorizer(true)

	output := subject.Paint("abc", colors.Red)

	assert.Equal(t, "\x1b[31mabc\x1b[0m", output)
}

func TestIntegrationPaintGreen(t *testing.T) {
	subject := colors.DefaultColorizer(true)

	output := subject.Paint("abc", colors.Green)

	assert.Equal(t, "\x1b[32mabc\x1b[0m", output)
}

func TestIntegrationEngineIsDisabled(t *testing.T) {
	subject := colors.DefaultColorizer(false)

	output := subject.Paint("abc", colors.Red)

	assert.Equal(t, "abc", output)
}
