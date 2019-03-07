package transforms

import (
	"github.com/testdouble/http-assertion-tool/parsers"
)

type Transform func(*parsers.Spec) error
