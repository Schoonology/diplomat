package errors

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"syscall"

	jsonschema "github.com/xeipuuv/gojsonschema"
	lua "github.com/yuin/gopher-lua"
	luaParse "github.com/yuin/gopher-lua/parse"
)

func formatURLError(err *url.Error) string {
	if opErr, ok := err.Err.(*net.OpError); ok {
		switch e := opErr.Err.(type) {
		case *net.DNSError:
			return fmt.Sprintf("Could not resolve host: %s\n", e.Name)
		case *os.SyscallError:
			return fmt.Sprintf("Failed to connect: %s\n", err.URL)
		default:
			return Format(err.Err)
		}
	}

	return Format(err.Err)
}

func formatLuaError(err *lua.ApiError) string {
	builder := strings.Builder{}

	switch err.Type {
	case lua.ApiErrorSyntax:
		builder.WriteString(fmt.Sprintf("Syntax error while parsing custom script:\n"))
	case lua.ApiErrorRun, lua.ApiErrorError:
		builder.WriteString(fmt.Sprintf("Error while running Lua script:\n"))
	case lua.ApiErrorPanic:
		builder.WriteString(fmt.Sprintf("Panic while running Lua script:\n"))
	default:
		builder.WriteString(fmt.Sprintf("Unknown error %v: %v", err.Type, err.Error()))
	}

	switch err.Object.Type() {
	case lua.LTString:
		builder.WriteString(fmt.Sprintf("	%s\n", err.Object))
	case lua.LTTable:
		err.Object.(*lua.LTable).ForEach(func(_ lua.LValue, value lua.LValue) {
			builder.WriteString(fmt.Sprintf("	%v", value))
		})
	}

	switch e := err.Cause.(type) {
	case *luaParse.Error:
		builder.WriteString(fmt.Sprintf("%s:%d:%d: %s", e.Pos.Source, e.Pos.Line, e.Pos.Column, e.Token))
	case nil:
	default:
		builder.WriteString(fmt.Sprintf("Unknown cause %T: %v", e, e))
	}

	return builder.String()
}

// Format returns a human-readable version of the error.
func Format(err error) string {
	switch e := err.(type) {
	case *lua.ApiError:
		return formatLuaError(e)
	case *lua.CompileError:
		return fmt.Sprintf("Compile error: %v -- %s", e.Line, e.Message)
	case *os.SyscallError:
		return fmt.Sprintf("Syscall error: %s\n%v", e.Syscall, Format(e.Err))
	case *url.Error:
		return formatURLError(e)
	case syscall.Errno:
		return fmt.Sprintf("Errno: %s\n", e)
	default:
		return fmt.Sprintf("%v\n", err.Error())
	}
}

// BuildErrorTable returns an LTTable filled with the errors in `result`.
func BuildErrorTable(L *lua.LState, result *jsonschema.Result) *lua.LTable {
	table := L.NewTable()

	for _, err := range result.Errors() {
		table.Append(lua.LString(err.String()))
	}

	return table
}
