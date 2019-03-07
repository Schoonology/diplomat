package transforms_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/http-assertion-tool/http"
	"github.com/testdouble/http-assertion-tool/parsers"
	"github.com/testdouble/http-assertion-tool/transforms"
)

// TODO(schoon) - There's probably a better way to test all this...do that.
func TestRenderHeaders(t *testing.T) {
	assert := assert.New(t)

	spec := parsers.Spec{
		Tests: []parsers.Test{parsers.Test{
			Request: &http.Request{
				Headers: map[string]string{
					"Key": "{{ test }}",
				},
			},
			Response: &http.Response{
				Headers: map[string]string{
					"Key": "{{ test }}",
				},
			},
		}},
	}

	err := transforms.RenderTemplates(&spec)
	assert.Nil(err)
	assert.Equal("test result", spec.Tests[0].Request.Headers["Key"])
	assert.Equal("test result", spec.Tests[0].Response.Headers["Key"])
}
