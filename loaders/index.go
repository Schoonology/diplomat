package loaders

type Loader interface {
	Load(string) (*Body, error)
}

type Body struct {
	Lines []string
}
