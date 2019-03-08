package transforms_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/diplomat/http"
	"github.com/testdouble/diplomat/parsers"
	"github.com/testdouble/diplomat/transforms"
)

// TODO(schoon) - There's probably a better way to test all this...do that.
func TestRenderHeaders(t *testing.T) {
	assert := assert.New(t)

	spec := parsers.Spec{
		Tests: []parsers.Test{parsers.Test{
			Request: &http.Request{
				Headers: map[string]string{
					"Key": "{{ __test }}",
				},
			},
			Response: &http.Response{
				Headers: map[string]string{
					"Key": "{{ __test }}",
				},
			},
		}},
	}

	err := transforms.RenderTemplates(&spec)
	assert.Nil(err)
	assert.Equal("this is a test", spec.Tests[0].Request.Headers["Key"])
	assert.Equal("this is a test", spec.Tests[0].Response.Headers["Key"])
}

func TestRenderBodies(t *testing.T) {
	assert := assert.New(t)

	spec := parsers.Spec{
		Tests: []parsers.Test{parsers.Test{
			Request: &http.Request{
				Body: []byte("{{ __test }}"),
			},
			Response: &http.Response{
				Body: []byte("{{ __test }}"),
			},
		}},
	}

	err := transforms.RenderTemplates(&spec)
	assert.Nil(err)
	assert.Equal([]byte("this is a test"), spec.Tests[0].Request.Body)
	assert.Equal([]byte("this is a test"), spec.Tests[0].Response.Body)
}
