// Package lua contains all pre-packaged scripts: templates, validators, etc.
//
//go:generate ../../bin/templify -p lua generate.lua
//go:generate ../../bin/templify -p lua validate.lua
package lua

import lua "github.com/yuin/gopher-lua"

var templates []func() string

func init() {
	templates = []func() string{
		generateTemplate,
		validateTemplate,
	}
}

// LoadAll loads all generated/packagd lua scripts into the provided `state`.
func LoadAll(state *lua.LState) error {
	for _, template := range templates {
		if err := state.DoString(template()); err != nil {
			return err
		}
	}

	return nil
}
