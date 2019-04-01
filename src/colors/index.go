package colors

// A ColorPrinter prints strings in colors, by delegating to an Engine.
type ColorPrinter interface {
	Print(str string, color Color)
}

// DefaultPrinter returns the standard color printer using the Aurora engine.
func DefaultPrinter(enabled bool) ColorPrinter {
	return &Printer{
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
