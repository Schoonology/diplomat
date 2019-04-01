package colors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/diplomat/colors"
	"github.com/testdouble/diplomat/mocks"
)

func TestUnitPaintRed(t *testing.T) {
	engine := mocks.Engine{}
	engine.
		On("Red", "abc").
		Return("red abc")
	subject := colors.Delegator{
		Engine: &engine,
	}

	output := subject.Paint("abc", colors.Red)

	assert.Equal(t, "red abc", output)
}

func TestUnitPaintGreen(t *testing.T) {
	engine := mocks.Engine{}
	engine.
		On("Green", "abc").
		Return("green abc")
	subject := colors.Delegator{
		Engine: &engine,
	}

	output := subject.Paint("abc", colors.Green)

	assert.Equal(t, "green abc", output)
}

func TestUnitPaintUnknownColor(t *testing.T) {
	engine := mocks.Engine{}
	subject := colors.Delegator{
		Engine: &engine,
	}

	output := subject.Paint("abc", "unknown")

	assert.Equal(t, "abc", output)
}
