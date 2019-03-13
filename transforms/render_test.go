package transforms_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testdouble/diplomat/http"
	"github.com/testdouble/diplomat/parsers"
	"github.com/testdouble/diplomat/transforms"
)

func streamTests(tests []parsers.Test) chan parsers.Test {
	testChannel := make(chan parsers.Test)

	go func() {
		for _, test := range tests {
			testChannel <- test
		}
	}()

	return testChannel
}

func assertTest(t *testing.T, expected parsers.Test, tests chan parsers.Test, errors chan error) {
	select {
	case err := <-errors:
		assert.FailNow(t, "Error should not exist.", err)
	case test := <-tests:
		assert.Equal(t, expected, test)
	}
}

// TODO(schoon) - There's probably a better way to test all this...do that.
func TestRenderHeaders(t *testing.T) {
	tests := streamTests([]parsers.Test{parsers.Test{
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
	}})
	errors := make(chan error)

	rendered := transforms.RenderTemplates(tests, errors)

	assertTest(t, parsers.Test{
		Request: &http.Request{
			Headers: map[string]string{"Key": "this is a test"},
		},
		Response: &http.Response{
			Headers: map[string]string{"Key": "this is a test"},
		},
	}, rendered, errors)
}

func TestRenderBodies(t *testing.T) {
	tests := streamTests([]parsers.Test{parsers.Test{
		Request: &http.Request{
			Body: []byte("{{ __test }}"),
		},
		Response: &http.Response{
			Body: []byte("{{ __test }}"),
		},
	}})
	errors := make(chan error)

	rendered := transforms.RenderTemplates(tests, errors)

	assertTest(t, parsers.Test{
		Request: &http.Request{
			Body: []byte("this is a test"),
		},
		Response: &http.Response{
			Body: []byte("this is a test"),
		},
	}, rendered, errors)
}
