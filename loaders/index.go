package loaders

// A Loader loads spec content from the provided location, sending each line
// to the returned channel. If any error is received on the provided channel,
// the Loader will halt and close the output channel.
type Loader interface {
	Load(string, chan error) chan string
}
