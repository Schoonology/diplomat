package http

type Client interface {
	Do(*Request) (*Response, error)
}

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    []byte
}

type Response struct {
	StatusCode int
	StatusText string
	Headers    map[string]string
	Body       []byte
}

func NewRequest(method string, path string) *Request {
	return &Request{
		Method:  method,
		Path:    path,
		Headers: make(map[string]string),
		Body:    make([]byte, 0),
	}
}

func NewResponse(statusCode int, statusText string) *Response {
	return &Response{
		StatusCode: statusCode,
		StatusText: statusText,
		Headers:    make(map[string]string),
		Body:       make([]byte, 0),
	}
}
