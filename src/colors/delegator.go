package colors

import (
	"fmt"
)

// Delegator prints strings with colors using an Engine.
type Delegator struct {
	Engine Engine
}

// Print outputs a string with a specified color.
func (d *Delegator) Print(str string, color Color) {
	var coloredString string

	switch color {
	case Red:
		coloredString = d.Engine.Red(str)
	case Green:
		coloredString = d.Engine.Green(str)
	}

	fmt.Print(coloredString)
}
