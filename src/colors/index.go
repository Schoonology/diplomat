package colors

import (
	"fmt"

	"github.com/logrusorgru/aurora"
)

type engine aurora.Aurora
type engineFunc func(interface{}) aurora.Value

// NewColorizer returns a configured colorizer
func NewColorizer(enabled bool) Colorizer {
	c := Colorizer{
		enabled: enabled,
	}
	c.init()
	return c
}

// Colorizer prints strings with colors.
type Colorizer struct {
	enabled bool
	colors  map[Color]engineFunc
}

func (c *Colorizer) init() {
	engine := aurora.NewAurora(c.enabled)

	c.colors = map[Color]engineFunc{
		Red:   engine.Red,
		Green: engine.Green,
	}
}

// Print outputs a string with a specified color.
func (c *Colorizer) Print(str string, color Color) {
	colorFunc := c.colors[color]
	fmt.Print(colorFunc(str))
}

// A Color is a representation of a printable color
type Color string

const (
	// Red for errors
	Red Color = "RED"
	// Green for success
	Green Color = "GREEN"
)
