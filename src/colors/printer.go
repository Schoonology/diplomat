package colors

import (
	"fmt"
)

// Printer prints strings with colors using an Engine.
type Printer struct {
	Engine Engine
}

// Print outputs a string with a specified color.
func (p *Printer) Print(str string, color Color) {
	var coloredString string

	switch color {
	case Red:
		coloredString = p.Engine.Red(str)
	case Green:
		coloredString = p.Engine.Green(str)
	}

	fmt.Print(coloredString)
}
