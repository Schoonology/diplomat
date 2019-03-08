package scripting

import (
	"fmt"

	scripts "github.com/testdouble/http-assertion-tool/scripting/lua"
	lua "github.com/yuin/gopher-lua"
)

var state *lua.LState

func init() {
	state = lua.NewState()

	if err := scripts.LoadAll(state); err != nil {
		panic(err)
	}
}

// RunPipeline returns the result of the function at `src`, which should
// return a string.
func RunPipeline(src string) (string, error) {
	err := state.DoString(fmt.Sprintf("return %s()", src))
	if err != nil {
		return "", err
	}

	ret := state.Get(-1)
	state.Pop(1)

	return ret.String(), nil
}

// RunValidator returns the result of the function at `src`, given `value`,
// which should return a boolean.
func RunValidator(src string, value string) (bool, error) {
	err := state.DoString(fmt.Sprintf("return %s(\"%s\")", src, value))
	if err != nil {
		return false, err
	}

	ret := state.Get(-1)
	state.Pop(1)

	return ret.(lua.LBool) == lua.LTrue, nil
}
