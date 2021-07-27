package transforms_test

import (
	"testing"

	"github.com/schoonology/diplomat/builders"
	"github.com/schoonology/diplomat/http"
	"github.com/schoonology/diplomat/transforms"
	"github.com/stretchr/testify/assert"
)

func TestRenderHeaders(t *testing.T) {
	input := builders.Test{
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
	}
	subject := new(transforms.TemplateRenderer)

	output, err := subject.Transform(input)

	assert.Nil(t, err)
	assert.Equal(t, builders.Test{
		Request: &http.Request{
			Headers: map[string]string{"Key": "this is a test"},
		},
		Response: &http.Response{
			Headers: map[string]string{"Key": "this is a test"},
		},
	}, output)
}

func TestRenderBodies(t *testing.T) {
	input := builders.Test{
		Request: &http.Request{
			Body: []byte("{{ __test }}"),
		},
		Response: &http.Response{
			Body: []byte("{{ __test }}"),
		},
	}
	subject := new(transforms.TemplateRenderer)

	output, err := subject.Transform(input)

	assert.Nil(t, err)
	assert.Equal(t, builders.Test{
		Request: &http.Request{
			Body: []byte("this is a test"),
		},
		Response: &http.Response{
			Body: []byte("this is a test"),
		},
	}, output)
}

func TestRenderPath(t *testing.T) {
	input := builders.Test{
		Request: &http.Request{
			Path: "/path?value={{ __test }}",
		},
		Response: &http.Response{},
	}
	subject := new(transforms.TemplateRenderer)

	output, err := subject.Transform(input)

	assert.Nil(t, err)
	assert.Equal(t, builders.Test{
		Request: &http.Request{
			Path: "/path?value=this is a test",
		},
		Response: &http.Response{},
	}, output)
}

func TestRenderWhitespace(t *testing.T) {
	input := builders.Test{
		Request: &http.Request{},
		Response: &http.Response{
			Body: []byte("{{ __test() .. __test() }}"),
		},
	}
	subject := new(transforms.TemplateRenderer)

	output, err := subject.Transform(input)

	assert.Nil(t, err)
	assert.Equal(t, builders.Test{
		Request: &http.Request{},
		Response: &http.Response{
			Body: []byte("this is a testthis is a test"),
		},
	}, output)
}
