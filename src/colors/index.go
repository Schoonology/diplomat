package colors

// A Colorizer prints strings in colors, by delegating to an Engine.
type Colorizer interface {
	Paint(str string, color Color) string
}

// DefaultColorizer returns the standard color printer using the Aurora engine.
func DefaultColorizer(enabled bool) Colorizer {
	return &Delegator{
		Engine: NewAuroraEngine(enabled),
	}
}

// An Engine implements methods for each printable color.
type Engine interface {
	Red(string) string
	Green(string) string
}

// A Color is a representation of a printable color
type Color string

const (
	// Red for errors
	Red Color = "RED"
	// Green for success
	Green Color = "GREEN"
)
