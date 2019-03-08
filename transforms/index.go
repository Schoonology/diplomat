package transforms

import (
	"github.com/testdouble/http-assertion-tool/parsers"
)

// Transform mutates a Spec in some meaningful way.
type Transform func(*parsers.Spec) error
