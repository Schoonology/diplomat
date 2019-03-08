package loaders

// A Loader loads spec content from the provided location.
type Loader interface {
	Load(string) (*Body, error)
}

// Body contains all lines from a Loader.
type Body struct {
	Lines []string
}
