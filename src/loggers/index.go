package loggers

// A Logger is responsible for writing streams to an output source.
type Logger interface {
	Print(str string)
	PrintAll(chan string)
}
