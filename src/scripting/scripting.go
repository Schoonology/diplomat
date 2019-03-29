package scripting

import (
	"fmt"
	"strings"

	"github.com/testdouble/diplomat/errors"
	scripts "github.com/testdouble/diplomat/scripting/lua"
	lua "github.com/yuin/gopher-lua"
)

var state *lua.LState

func init() {
	state = lua.NewState()

	if err := scripts.LoadAll(state); err != nil {
		panic(err)
	}
}

// LoadFile runs the Lua file in `filename`. Theoretically this file
// should add custom functions into the global namespace.
func LoadFile(filename string) error {
	return state.DoFile(filename)
}

// RunPipeline returns the result of the function at `src`, which should
// return a string.
func RunPipeline(src string) (string, error) {
	fn, err := state.Load(strings.NewReader(fmt.Sprintf("return %s", src)), src)
	if err != nil {
		return "", err
	}

	state.Push(fn)
	err = state.PCall(0, 1, nil)
	if err != nil {
		return "", nil
	}

	value := state.Get(-1)

	if value.Type() != lua.LTNil {
		state.Pop(1)
	}

	switch value.Type() {
	case lua.LTFunction:
		state.Push(value)
		err = state.PCall(0, 1, nil)
		if err != nil {
			return "", err
		}

		value = state.Get(-1)
		state.Pop(1)
	case lua.LTNil:
		return "", &errors.MissingTemplate{Template: src}
	}

	return value.String(), nil
}

// RunValidator returns the result of the function at `src`, given `value`,
// which should return a boolean.
func RunValidator(src string, value string) (bool, error) {
	fn, err := state.Load(strings.NewReader(fmt.Sprintf("return %s", src)), src)
	if err != nil {
		return false, err
	}

	state.Push(fn)
	err = state.PCall(0, 1, nil)
	if err != nil {
		return false, err
	}

	validator := state.Get(-1)
	if validator.Type() == lua.LTNil {
		return false, &errors.MissingTemplate{Template: src}
	}

	state.Push(lua.LString(value))

	err = state.PCall(1, 1, nil)
	if err != nil {
		return false, err
	}

	ret := state.Get(-1)
	state.Pop(1)

	return lua.LVAsBool(ret), nil
}
