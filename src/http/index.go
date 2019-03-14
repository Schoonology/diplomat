package http

// A Client can execute a single Request, providing the Response.
type Client interface {
	Do(*Request) (*Response, error)
}

// Request contains all data required to define a single HTTP request.
type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    []byte
}

// Response contains all data required to define a single HTTP response.
type Response struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       []byte
}

// NewRequest returns a new Request wrapping the `method` and `path`.
func NewRequest(method string, path string) *Request {
	return &Request{
		Method:  method,
		Path:    path,
		Headers: make(map[string]string),
		Body:    make([]byte, 0),
	}
}

// NewResponse returns a new Response wrapping the `statusCode` and
// `statusText`.
func NewResponse(statusCode int, statusText string) *Response {
	return &Response{
		StatusCode: statusCode,
		StatusText: statusText,
		Headers:    make(map[string]string),
		Body:       make([]byte, 0),
	}
}
