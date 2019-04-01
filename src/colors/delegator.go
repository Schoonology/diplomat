package colors

// Delegator prints strings with colors using an Engine.
type Delegator struct {
	Engine Engine
}

// Paint changes the print color of a string.
func (d *Delegator) Paint(str string, color Color) string {
	var coloredString string

	switch color {
	case Red:
		coloredString = d.Engine.Red(str)
	case Green:
		coloredString = d.Engine.Green(str)
	default:
		coloredString = str
	}

	return coloredString
}
