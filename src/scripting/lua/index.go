// Package lua contains all pre-packaged scripts: templates, validators, etc.
//
//go:generate ../../../bin/templify -p lua generate.lua
//go:generate ../../../bin/templify -p lua util.lua
//go:generate ../../../bin/templify -p lua validate.lua
package lua

import (
	"github.com/testdouble/diplomat/errors"
	jsonschema "github.com/xeipuuv/gojsonschema"
	lua "github.com/yuin/gopher-lua"
)

var templates []func() string

func init() {
	templates = []func() string{
		generateTemplate,
		utilTemplate,
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
		L.RaiseError(err.Error())
	}

	if len(result.Errors()) > 0 {
		L.Error(errors.BuildErrorTable(L, result), 1)
	}

	L.Push(lua.LBool(result.Valid()))

	return 1
}