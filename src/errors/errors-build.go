package errors

import (
	"fmt"

	"github.com/testdouble/diplomat/parsers"
)

// NewBuildError wraps an error with a BuildError containing a LineNumber.
func NewBuildError(paragraph parsers.Paragraph, err error) *BuildError {
	return &BuildError{
		LineNumber: paragraph.LineNumber,
		Err:        err,
	}
}

// BuildError is the error type for any error found during the build step.
type BuildError struct {
	LineNumber int
	Err        error
}

func (err *BuildError) Error() string {
	return fmt.Sprintf("Error building spec: line %v\n	%s\n", err.LineNumber, err.Err.Error())
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
