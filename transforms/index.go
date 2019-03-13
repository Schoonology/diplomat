package transforms

import (
	"github.com/testdouble/diplomat/parsers"
)

// Transformer is capable of mutating Tests.
type Transformer interface {
	Transform(parsers.Test) (parsers.Test, error)
	TransformAll(chan parsers.Test, chan error) chan parsers.Test
}
