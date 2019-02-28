package http

type Differ interface {
	Diff(*Response, *Response) (string, error)
}
