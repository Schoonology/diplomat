package transforms

import (
	"github.com/testdouble/diplomat/parsers"
)

// Transform mutates a Spec in some meaningful way.
type Transform func(*parsers.Spec) error
