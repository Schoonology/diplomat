// Package lua contains all pre-packaged scripts: templates, validators, etc.
//
//go:generate ../../bin/templify -p lua generate.lua
//go:generate ../../bin/templify -p lua validate.lua
package lua

import (
	jsonschema "github.com/xeipuuv/gojsonschema"
	lua "github.com/yuin/gopher-lua"
)

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

	state.SetGlobal("__validateJSON", state.NewFunction(validateJSONSchema))

	return nil
}

func validateJSONSchema(L *lua.LState) int {
	schema := L.ToString(1)
	value := L.ToString(2)

	result, err := jsonschema.Validate(
		jsonschema.NewStringLoader(schema),
		jsonschema.NewStringLoader(value),
	)
	if err != nil {
		// TODO(schoon) - Where does this go?
		panic(err)
	}

	// TODO(schoon) - Provide result.Errors() to Diplomat proper.

	L.Push(lua.LBool(result.Valid()))

	return 1
}
