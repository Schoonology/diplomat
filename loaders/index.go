package loaders

type FileLoader interface {
	Load(string) (*File, error)
}

type File struct{}
