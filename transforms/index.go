package transforms

import (
	"github.com/testdouble/diplomat/parsers"
)

// Transformer is capable of mutating Tests.
type Transformer interface {
	Transform(chan parsers.Test, chan error) chan parsers.Test
}
