package loaders

// A Loader loads spec content from the provided location.
type Loader interface {
	Streamer
	Load(string) (*Body, error)
}

// A Streamer returns a channel and streams file contents to that channel.
type Streamer interface {
	Stream(string) (chan *Body, chan error)
}

// Body contains all lines from a Loader.
type Body struct {
	Lines []string
}
