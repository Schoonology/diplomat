package transforms

import "github.com/testdouble/diplomat/builders"

// Transformer is capable of mutating Tests.
type Transformer interface {
	Transform(builders.Test) (builders.Test, error)
	TransformAll(chan builders.Test) chan builders.Test
}
