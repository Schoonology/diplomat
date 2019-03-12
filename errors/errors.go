package errors

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"syscall"
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

// Display emits a human-readable version of the error to STDOUT.
func Display(err error) {
	switch e := err.(type) {
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

// MissingResponse is the error type for a request without an expected response.
type MissingResponse struct{}

func (err *MissingResponse) Error() string {
	return "Found a request without a corresponding response."
}
