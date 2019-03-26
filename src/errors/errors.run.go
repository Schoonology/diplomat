package errors

// NewAssertionError creates an error representing a test failure.
func NewAssertionError(diff string) *AssertionError {
	return &AssertionError{
		Diff: diff,
	}
}

// AssertionError is the error type for any error found during the build step.
type AssertionError struct {
	Diff string
}

func (err *AssertionError) Error() string {
	return err.Diff
}
