package builders_test

import (
	"testing"

	"github.com/schoonology/diplomat/builders"
	"github.com/schoonology/diplomat/http"
	"github.com/stretchr/testify/assert"
)

func assertTest(t *testing.T, expected builders.Test, actual builders.Test, err error) {
	if err != nil {
		assert.FailNow(t, "Error should not exist.")
		return
	}
	assert.Equal(t, expected, actual)
}

func fillRequest(method string, path string, headers map[string]string, body string) *http.Request {
	request := http.NewRequest(method, path)
	for key, value := range headers {
		request.Headers[key] = value
	}
	request.Body = []byte(body)
	return request
}

func fillResponse(code int, status string, headers map[string]string, body string) *http.Response {
	response := http.NewResponse(code, status)
	for key, value := range headers {
		response.Headers[key] = value
	}
	response.Body = []byte(body)
	return response
}
