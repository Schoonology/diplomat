// Package lua contains all pre-packaged scripts: templates, validators, etc.
//
//go:generate ../../../bin/templify -p lua context.lua
//go:generate ../../../bin/templify -p lua generate.lua
//go:generate ../../../bin/templify -p lua util.lua
//go:generate ../../../bin/templify -p lua validate.lua
package lua

import (
	"net/http"
	"strings"

	"github.com/cjoudrey/gluahttp"
	json "github.com/layeh/gopher-json"
	"github.com/schoonology/diplomat/errors"
	jsonschema "github.com/xeipuuv/gojsonschema"
	"github.com/yuin/gluare"
	lua "github.com/yuin/gopher-lua"
)

var templates map[string]func() string

func init() {
	templates = map[string]func() string{
		"context.lua":  contextTemplate,
		"generate.lua": generateTemplate,
		"util.lua":     utilTemplate,
		"validate.lua": validateTemplate,
	}
}

// LoadAll loads all generated/packagd lua scripts into the provided `state`.
func LoadAll(state *lua.LState) error {
	for name, template := range templates {
		fn, err := state.Load(strings.NewReader(template()), name)
		if err != nil {
			return err
		}

		state.Push(fn)
		err = state.PCall(0, 1, nil)
		if err != nil {
			return err
		}
	}

	json.Preload(state)
	err := state.DoString("json = require('json')")
	if err != nil {
		return err
	}

	client := &http.Client{}
	state.PreloadModule("http", gluahttp.NewHttpModuleWithDo(
		func(req *http.Request) (*http.Response, error) {
			req.Header.Set("User-Agent", "Diplomat/0.0.1")
			return client.Do(req)
		},
	).Loader)
	err = state.DoString("http = require('http')")
	if err != nil {
		return err
	}

	state.PreloadModule("re", gluare.Loader)
	err = state.DoString("re = require('re')")
	if err != nil {
		return err
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
