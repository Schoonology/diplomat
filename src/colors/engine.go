package colors

import (
	"github.com/logrusorgru/aurora"
)

// Aurora is an Engine implemented using the aurora package.
type Aurora struct {
	au aurora.Aurora
}

// NewAuroraEngine returns an Aurora engine that is enabled or disabled.
func NewAuroraEngine(enabled bool) *Aurora {
	return &Aurora{
		au: aurora.NewAurora(enabled),
	}
}

// Red colors a string red.
func (a *Aurora) Red(str string) string {
	return a.au.Red(str).String()
}

// Green colors a string green.
func (a *Aurora) Green(str string) string {
	return a.au.Green(str).String()
}
