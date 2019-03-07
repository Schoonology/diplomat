//go:generate ../../bin/templify -p lua generate.lua
package lua

import lua "github.com/yuin/gopher-lua"

var templates []func() string

func init() {
	templates = []func() string{
		generateTemplate,
	}
}

func LoadAll(state *lua.LState) error {
	for _, template := range templates {
		if err := state.DoString(template()); err != nil {
			return err
		}
	}

	return nil
}
