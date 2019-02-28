package http

type Client interface {
	Do(*Request) (*Response, error)
}

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
}

type Response struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
}

func NewRequest(method string, path string) *Request {
	return &Request{
		Method:  method,
		Path:    path,
		Headers: make(map[string]string),
	}
}

func NewResponse(statusCode int, statusText string) *Response {
	return &Response{
		StatusCode: statusCode,
		StatusText: statusText,
		Headers:    make(map[string]string),
	}
}
