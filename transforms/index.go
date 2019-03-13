package transforms

import (
	"github.com/testdouble/diplomat/parsers"
)

// Transform mutates a Test, returning the mutated Test.
type Transform func(chan parsers.Test, chan error) chan parsers.Test
