package transforms

import (
	"github.com/testdouble/diplomat/parsers"
)

// Transform mutates a Spec in some meaningful way.
type Transform func(*parsers.Spec) error

// TransformStream mutates a Spec via stream.
type TransformStream func(chan *parsers.Spec, chan error) (chan *parsers.Spec, chan error)
