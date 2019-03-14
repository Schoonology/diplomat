package errors

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"syscall"

	jsonschema "github.com/xeipuuv/gojsonschema"
	lua "github.com/yuin/gopher-lua"
	luaParse "github.com/yuin/gopher-lua/parse"
)

func displayURLError(err *url.Error) {
	if opErr, ok := err.Err.(*net.OpError); ok {
		switch e := opErr.Err.(type) {
		case *net.DNSError:
			fmt.Printf("Could not resolve host: %s\n", e.Name)
		case *os.SyscallError:
			fmt.Printf("Failed to connect: %s\n", err.URL)
		default:
			Display(err.Err)
		}
		return
	}

	Display(err.Err)
}

func displayLuaError(err *lua.ApiError) {
	switch err.Type {
	case lua.ApiErrorSyntax:
		fmt.Printf("Syntax error while parsing custom script:\n")
	case lua.ApiErrorRun, lua.ApiErrorError:
		fmt.Printf("Error while running Lua script:\n")
	case lua.ApiErrorPanic:
		fmt.Printf("Panic while running Lua script:\n")
	default:
		fmt.Printf("Unknown error %v: %v", err.Type, err.Error())
	}

	switch err.Object.Type() {
	case lua.LTString:
		fmt.Printf("	%s\n", err.Object)
	case lua.LTTable:
		err.Object.(*lua.LTable).ForEach(func(_ lua.LValue, value lua.LValue) {
			fmt.Printf("	%v", value)
		})
	}

	switch e := err.Cause.(type) {
	case *luaParse.Error:
		fmt.Printf("%s:%d:%d: %s", e.Pos.Source, e.Pos.Line, e.Pos.Column, e.Token)
	case nil:
	default:
		fmt.Printf("Unknown cause %T: %v", e, e)
	}
}

// Display emits a human-readable version of the error to STDOUT.
func Display(err error) {
	switch e := err.(type) {
	case *lua.ApiError:
		displayLuaError(e)
	case *lua.CompileError:
		fmt.Printf("Compile error: %v -- %s", e.Line, e.Message)
	case *os.SyscallError:
		fmt.Printf("Syscall error: %s\n", e.Syscall)
		Display(e.Err)
	case *url.Error:
		displayURLError(e)
	case syscall.Errno:
		fmt.Printf("Errno: %s\n", e)
	default:
		fmt.Println(err.Error())
	}
}

// BadHeader is the error type for a badly formatted header.
type BadHeader struct {
	Header string
}

func (err *BadHeader) Error() string {
	return fmt.Sprintf("Failed to parse header: %s", err.Header)
}

// BadRequestLine is the error type for a badly formatted "request line":
// the first line of a request with `METHOD PATH PROTO`.
type BadRequestLine struct {
	Line string
}

func (err *BadRequestLine) Error() string {
	return fmt.Sprintf("Failed to parse request line: %s", err.Line)
}

// BadResponseStatus is the error type for a badly formatted "response status":
// the first line of a response with `PROTO CODE STATUS`.
type BadResponseStatus struct {
	Line string
}

func (err *BadResponseStatus) Error() string {
	return fmt.Sprintf("Failed to parse response line: %s", err.Line)
}

// MissingRequest is the error type for an expected response without a
// corresponding request.
type MissingRequest struct{}

func (err *MissingRequest) Error() string {
	return "Found a response without a corresponding request."
}

// MisplacedRequest is the error type for an unexpected request line in or after a response.
type MisplacedRequest struct{}

func (err *MisplacedRequest) Error() string {
	return "Found a request in or after the response. A spec should have one request that comes before the response."
}

// MissingResponse is the error type for a request without an expected response.
type MissingResponse struct{}

func (err *MissingResponse) Error() string {
	return "Found a request without a corresponding response."
}

// MissingTemplate is the error type for a template that could not be found.
type MissingTemplate struct {
	Template string
}

func (err *MissingTemplate) Error() string {
	return fmt.Sprintf("Template `%s` could not be found.", err.Template)
}

// BuildErrorTable returns an LTTable filled with the errors in `result`.
func BuildErrorTable(L *lua.LState, result *jsonschema.Result) *lua.LTable {
	table := L.NewTable()

	for _, err := range result.Errors() {
		table.Append(lua.LString(err.String()))
	}

	return table
}
